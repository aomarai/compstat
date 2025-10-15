package benchmark

// Result Represents a single benchmark result
type Result struct {
	RunID                 string  `json:"run_id"`
	Algorithm             string  `json:"algorithm"`
	Level                 int     `json:"level"`
	CompressThreads       int     `json:"compress_threads"`
	DecompressThreads     int     `json:"decompress_threads"`
	FilePath              string  `json:"file_path"`
	UncompressedBytes     int64   `json:"uncompressed_bytes"`
	CompressedBytes       int64   `json:"compressed_bytes"`
	CompressionRatio      float64 `json:"compression_ratio"`
	CompressionTimeS      float64 `json:"compression_time_s"`
	DecompressionTimeS    float64 `json:"decompression_time_s"`
	CompressionSpeedMBs   float64 `json:"compression_speed_mbs"`
	DecompressionSpeedMBs float64 `json:"decompression_speed_mbs"`
	CompressionMaxRSSMB   float64 `json:"compression_max_rss_mb"`
	DecompressionMaxRSSMB float64 `json:"decompression_max_rss_mb"`
	Verified              bool    `json:"verified"`
	Iteration             int     `json:"iteration"`
}

// Config holds benchmark configuration
type Config struct {
	Files               []string
	Codecs              []string
	CompressThreads     int
	DecompressThreads   int
	Iterations          int
	TmpDir              string
	OutputCSV           string
	OutputJSON          string
	VerifyDecompression bool
	SkipDecompression   bool
	Parallelism         int
}

// Job represents a single benchmark job
type job struct {
	filePath  string
	codec     interface{} // Will be codec.Codec, using interface{} to avoid import cycle
	level     int
	iteration int
}
