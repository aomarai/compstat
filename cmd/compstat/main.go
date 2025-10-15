package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aomarai/compstat/internal/benchmark"
	"github.com/aomarai/compstat/internal/codec"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	files := flag.String("files", "", "Comma-separated list of input files (required)")
	codecs := flag.String("codecs", "", "Comma-separated codecs (default: all available)")
	compThreads := flag.Int("compress-threads", 0, "Compression threads (default: CPU count)")
	decompThreads := flag.Int("decompress-threads", 0, "Decompression threads (default: CPU count)")
	iterations := flag.Int("iterations", 1, "Number of iterations per configuration")
	tmpDir := flag.String("tmpdir", "", "Temporary directory (default: system temp)")
	output := flag.String("output", "compstat_results.csv", "CSV output file")
	jsonOutput := flag.String("json", "", "Optional JSON output file")
	noVerify := flag.Bool("no-verify", false, "Skip decompression verification")
	skipDecomp := flag.Bool("skip-decompression", false, "Skip decompression entirely")
	parallelism := flag.Int("parallelism", 1, "Number of parallel benchmark jobs")
	version := flag.Bool("version", false, "Show version information")

	flag.Parse()

	if *version {
		fmt.Printf("compstat version %s (built %s)\n", Version, BuildTime)
		os.Exit(0)
	}

	if *files == "" {
		fmt.Println("Error: -files is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse inputs
	fileList := strings.Split(*files, ",")
	for i := range fileList {
		fileList[i] = strings.TrimSpace(fileList[i])
	}

	codecList := make([]string, 0)
	if *codecs == "" {
		for name := range codec.Registry {
			codecList = append(codecList, name)
		}
	} else {
		for _, c := range strings.Split(*codecs, ",") {
			codecList = append(codecList, strings.TrimSpace(c))
		}
	}

	cpuCount := runtime.NumCPU()
	if *compThreads == 0 {
		*compThreads = cpuCount
	}
	if *decompThreads == 0 {
		*decompThreads = cpuCount
	}

	tmpDirPath := *tmpDir
	if tmpDirPath == "" {
		tmpDirPath = filepath.Join(os.TempDir(), "compstat_tmp")
	}

	config := benchmark.Config{
		Files:               fileList,
		Codecs:              codecList,
		CompressThreads:     *compThreads,
		DecompressThreads:   *decompThreads,
		Iterations:          *iterations,
		TmpDir:              tmpDirPath,
		OutputCSV:           *output,
		OutputJSON:          *jsonOutput,
		VerifyDecompression: !*noVerify,
		SkipDecompression:   *skipDecomp,
		Parallelism:         *parallelism,
	}

	runner, err := benchmark.NewRunner(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer runner.Close()

	if err := runner.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if config.OutputJSON != "" {
		if err := runner.WriteJSONSummary(); err != nil {
			fmt.Printf("Warning: Failed to write JSON: %v\n", err)
		}
	}

	fmt.Printf("\n✓ Benchmark complete! Results: %s\n", config.OutputCSV)
	fmt.Printf("  Total runs: %d\n", runner.ResultCount())
}
