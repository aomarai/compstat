package benchmark

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/aomarai/compstat/internal/codec"
	"github.com/aomarai/compstat/internal/util"
)

// Runner orchestrates benchmark execution
type Runner struct {
	config     Config
	results    []Result
	resultsMux sync.Mutex
	fileHashes map[string]string
	csvFile    *os.File
	csvWriter  *csv.Writer
	csvMux     sync.Mutex
}

// NewRunner creates a new benchmark runner
func NewRunner(config Config) (*Runner, error) {
	runner := &Runner{
		config:     config,
		results:    make([]Result, 0),
		fileHashes: make(map[string]string),
	}

	// Create tmpdir
	if err := os.MkdirAll(config.TmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create tmpdir: %w", err)
	}

	// Open CSV file for writing
	csvFile, err := os.OpenFile(config.OutputCSV, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	runner.csvFile = csvFile
	runner.csvWriter = csv.NewWriter(csvFile)

	// Write header if new file
	fileInfo, _ := csvFile.Stat()
	if fileInfo.Size() == 0 {
		runner.writeCSVHeader()
	}

	return runner, nil
}

// Close cleans up resources
func (r *Runner) Close() {
	if r.csvWriter != nil {
		r.csvWriter.Flush()
	}
	if r.csvFile != nil {
		r.csvFile.Close()
	}
}

// ResultCount returns the number of completed benchmarks
func (r *Runner) ResultCount() int {
	r.resultsMux.Lock()
	defer r.resultsMux.Unlock()
	return len(r.results)
}

func (r *Runner) writeCSVHeader() {
	header := []string{
		"run_id", "algorithm", "level", "compress_threads", "decompress_threads",
		"file_path", "uncompressed_bytes", "compressed_bytes", "compression_ratio",
		"compression_time_s", "decompression_time_s", "compression_speed_mbs",
		"decompression_speed_mbs", "compression_max_rss_mb", "decompression_max_rss_mb",
		"verified", "iteration",
	}
	r.csvWriter.Write(header)
	r.csvWriter.Flush()
}

func (r *Runner) writeResult(result Result) {
	r.csvMux.Lock()
	defer r.csvMux.Unlock()

	row := []string{
		result.RunID,
		result.Algorithm,
		strconv.Itoa(result.Level),
		strconv.Itoa(result.CompressThreads),
		strconv.Itoa(result.DecompressThreads),
		result.FilePath,
		strconv.FormatInt(result.UncompressedBytes, 10),
		strconv.FormatInt(result.CompressedBytes, 10),
		fmt.Sprintf("%.4f", result.CompressionRatio),
		fmt.Sprintf("%.3f", result.CompressionTimeS),
		fmt.Sprintf("%.3f", result.DecompressionTimeS),
		fmt.Sprintf("%.2f", result.CompressionSpeedMBs),
		fmt.Sprintf("%.2f", result.DecompressionSpeedMBs),
		fmt.Sprintf("%.2f", result.CompressionMaxRSSMB),
		fmt.Sprintf("%.2f", result.DecompressionMaxRSSMB),
		strconv.FormatBool(result.Verified),
		strconv.Itoa(result.Iteration),
	}
	err := r.csvWriter.Write(row)
	if err != nil {
		fmt.Printf("Failed to write CSV row: %v\n", err)
		return
	}
	r.csvWriter.Flush()
	if err := r.csvWriter.Error(); err != nil {
		fmt.Printf("Failed to flush CSV Writer: %v\n", err)
	}
}

// PrecomputeHashes computes file hashes upfront for verification
func (r *Runner) PrecomputeHashes() error {
	if !r.config.VerifyDecompression {
		return nil
	}

	fmt.Println("\n=== Pre-computing file hashes ===")
	for _, filePath := range r.config.Files {
		fmt.Printf("Hashing %s... ", filepath.Base(filePath))
		hash, err := util.ComputeFileHash(filePath)
		if err != nil {
			fmt.Printf("failed: %v\n", err)
			continue
		}
		r.fileHashes[filePath] = hash
		fmt.Println("✓")
	}
	return nil
}

// GetAvailableCodecs filters codecs to only available ones
func (r *Runner) GetAvailableCodecs() []codec.Codec {
	available := make([]codec.Codec, 0)
	for _, name := range r.config.Codecs {
		c, ok := codec.Registry[name]
		if !ok {
			fmt.Printf("Unknown codec: %s\n", name)
			continue
		}
		if c.IsAvailable() {
			available = append(available, c)
		} else {
			fmt.Printf("Skipping %s: %s not found\n", name, c.Binary())
		}
	}
	return available
}

// Run executes the full benchmark suite
func (r *Runner) Run() error {
	if err := r.PrecomputeHashes(); err != nil {
		return err
	}

	codecs := r.GetAvailableCodecs()
	if len(codecs) == 0 {
		return fmt.Errorf("no codecs available")
	}

	fmt.Printf("\n=== Benchmarking %d file(s) with %d codec(s) ===\n", len(r.config.Files), len(codecs))

	// Build job queue
	jobs := make([]job, 0)
	for _, filePath := range r.config.Files {
		for _, c := range codecs {
			for _, level := range c.Levels() {
				for iter := 1; iter <= r.config.Iterations; iter++ {
					jobs = append(jobs, job{
						filePath:  filePath,
						codec:     c,
						level:     level,
						iteration: iter,
					})
				}
			}
		}
	}

	fmt.Printf("Total benchmark runs: %d\n", len(jobs))

	// Process jobs in parallel
	jobChan := make(chan job, len(jobs))
	for _, j := range jobs {
		jobChan <- j
	}
	close(jobChan)

	var wg sync.WaitGroup
	for i := 0; i < r.config.Parallelism; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range jobChan {
				result := r.runSingleBenchmark(j)
				if result != nil {
					r.resultsMux.Lock()
					r.results = append(r.results, *result)
					r.resultsMux.Unlock()
					r.writeResult(*result)
				}
			}
		}(i)
	}

	wg.Wait()
	return nil
}

// WriteJSONSummary writes results to JSON file
func (r *Runner) WriteJSONSummary() error {
	if r.config.OutputJSON == "" {
		return nil
	}

	data, err := json.MarshalIndent(r.results, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.config.OutputJSON, data, 0644)
}

func (r *Runner) runSingleBenchmark(j job) *Result {
	timestamp := time.Now().Unix()
	c := j.codec.(codec.Codec)
	runID := fmt.Sprintf("%d_%s_%s_%d_%d", timestamp, filepath.Base(j.filePath), c.Name(), j.level, j.iteration)

	fmt.Printf("[%d/%d] %s - %s level %d\n", j.iteration, r.config.Iterations, filepath.Base(j.filePath), c.Name(), j.level)

	// Setup paths
	compOut := filepath.Join(r.config.TmpDir, fmt.Sprintf("%s.%s.%d.%d%s", filepath.Base(j.filePath), c.Name(), j.level, j.iteration, c.Extension()))
	decompOut := filepath.Join(r.config.TmpDir, fmt.Sprintf("%s.decompressed.%d", filepath.Base(j.filePath), j.iteration))

	// Determine thread counts
	compThreads := r.config.CompressThreads
	decompThreads := r.config.DecompressThreads
	if !c.SupportsThreading() {
		compThreads = 1
		decompThreads = 1
	}

	// Get uncompressed size
	uncompSize, err := util.FileSize(j.filePath)
	if err != nil {
		fmt.Printf("  ! Failed to get file size: %v\n", err)
		return nil
	}

	// Compression
	compCmd := c.CompressCommand(j.level, compThreads, j.filePath, compOut)
	compTime, err := util.RunCommand(c.Binary(), compCmd, compOut)
	if err != nil {
		fmt.Printf("  ! Compression failed: %v\n", err)
		os.Remove(compOut)
		return nil
	}

	compSize, err := util.FileSize(compOut)
	if err != nil {
		fmt.Printf("  ! Failed to get compressed size: %v\n", err)
		os.Remove(compOut)
		return nil
	}

	// Calculate compression metrics
	compTimeSec := compTime.Seconds()
	compSpeedMBs := float64(uncompSize) / (1024 * 1024) / compTimeSec
	ratio := float64(compSize) / float64(uncompSize)

	result := &Result{
		RunID:               runID,
		Algorithm:           c.Name(),
		Level:               j.level,
		CompressThreads:     compThreads,
		DecompressThreads:   decompThreads,
		FilePath:            j.filePath,
		UncompressedBytes:   uncompSize,
		CompressedBytes:     compSize,
		CompressionRatio:    ratio,
		CompressionTimeS:    compTimeSec,
		CompressionSpeedMBs: compSpeedMBs,
		Iteration:           j.iteration,
	}

	// Decompression
	if !r.config.SkipDecompression {
		decompCmd := c.DecompressCommand(decompThreads, compOut, decompOut)
		decompTime, err := util.RunCommand(c.Binary(), decompCmd, decompOut)
		if err != nil {
			fmt.Printf("  ! Decompression failed: %v\n", err)
		} else {
			decompTimeSec := decompTime.Seconds()
			result.DecompressionTimeS = decompTimeSec
			result.DecompressionSpeedMBs = float64(uncompSize) / (1024 * 1024) / decompTimeSec

			// Verify if requested
			if r.config.VerifyDecompression {
				if origHash, ok := r.fileHashes[j.filePath]; ok {
					decompHash, err := util.ComputeFileHash(decompOut)
					if err == nil && decompHash == origHash {
						result.Verified = true
						fmt.Println("  ✓ Verified")
					} else {
						fmt.Println("  ! Verification failed")
					}
				}
			}
		}
		os.Remove(decompOut)
	}

	// Cleanup
	os.Remove(compOut)

	return result
}
