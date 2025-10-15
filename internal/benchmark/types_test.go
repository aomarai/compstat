package benchmark

import (
	"testing"
)

func TestResultStructure(t *testing.T) {
	result := Result{
		RunID:                 "test_run_1",
		Algorithm:             "zstd",
		Level:                 3,
		CompressThreads:       4,
		DecompressThreads:     4,
		FilePath:              "/path/to/file.bin",
		UncompressedBytes:     1000000,
		CompressedBytes:       500000,
		CompressionRatio:      0.5,
		CompressionTimeS:      1.5,
		DecompressionTimeS:    0.5,
		CompressionSpeedMBs:   666.67,
		DecompressionSpeedMBs: 2000.0,
		CompressionMaxRSSMB:   100.0,
		DecompressionMaxRSSMB: 50.0,
		Verified:              true,
		Iteration:             1,
	}

	// Verify fields are set correctly
	if result.Algorithm != "zstd" {
		t.Errorf("Expected Algorithm 'zstd', got '%s'", result.Algorithm)
	}
	if result.Level != 3 {
		t.Errorf("Expected Level 3, got %d", result.Level)
	}
	if result.CompressionRatio != 0.5 {
		t.Errorf("Expected CompressionRatio 0.5, got %f", result.CompressionRatio)
	}
	if !result.Verified {
		t.Error("Expected Verified to be true")
	}
}

func TestConfigStructure(t *testing.T) {
	config := Config{
		Files:               []string{"file1.bin", "file2.bin"},
		Codecs:              []string{"zstd", "gzip"},
		CompressThreads:     4,
		DecompressThreads:   4,
		Iterations:          3,
		TmpDir:              "/tmp",
		OutputCSV:           "results.csv",
		OutputJSON:          "results.json",
		VerifyDecompression: true,
		SkipDecompression:   false,
		Parallelism:         2,
	}

	// Verify fields are set correctly
	if len(config.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(config.Files))
	}
	if len(config.Codecs) != 2 {
		t.Errorf("Expected 2 codecs, got %d", len(config.Codecs))
	}
	if config.Iterations != 3 {
		t.Errorf("Expected 3 iterations, got %d", config.Iterations)
	}
	if config.Parallelism != 2 {
		t.Errorf("Expected parallelism 2, got %d", config.Parallelism)
	}
	if !config.VerifyDecompression {
		t.Error("Expected VerifyDecompression to be true")
	}
	if config.SkipDecompression {
		t.Error("Expected SkipDecompression to be false")
	}
}
