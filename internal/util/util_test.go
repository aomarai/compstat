package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSize(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	
	content := []byte("Hello, World!")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test file size
	size, err := FileSize(tmpFile)
	if err != nil {
		t.Fatalf("FileSize failed: %v", err)
	}

	expectedSize := int64(len(content))
	if size != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, size)
	}

	// Test non-existent file
	_, err = FileSize(filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestComputeFileHash(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	
	content := []byte("Hello, World!")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Compute hash
	hash1, err := ComputeFileHash(tmpFile)
	if err != nil {
		t.Fatalf("ComputeFileHash failed: %v", err)
	}

	if hash1 == "" {
		t.Error("Expected non-empty hash")
	}

	// Verify hash is consistent
	hash2, err := ComputeFileHash(tmpFile)
	if err != nil {
		t.Fatalf("ComputeFileHash failed on second call: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("Hash mismatch: %s != %s", hash1, hash2)
	}

	// Test with different content
	content2 := []byte("Different content")
	tmpFile2 := filepath.Join(tmpDir, "test2.txt")
	if err := os.WriteFile(tmpFile2, content2, 0644); err != nil {
		t.Fatalf("Failed to create second test file: %v", err)
	}

	hash3, err := ComputeFileHash(tmpFile2)
	if err != nil {
		t.Fatalf("ComputeFileHash failed for second file: %v", err)
	}

	if hash1 == hash3 {
		t.Error("Expected different hashes for different content")
	}

	// Test non-existent file
	_, err = ComputeFileHash(filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestNeedsStdoutRedirection(t *testing.T) {
	tests := []struct {
		binary   string
		expected bool
	}{
		{"xz", true},
		{"pigz", true},
		{"pbzip2", true},
		{"gzip", false},
		{"zstd", false},
		{"lz4", false},
		{"brotli", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.binary, func(t *testing.T) {
			result := NeedsStdoutRedirection(tt.binary)
			if result != tt.expected {
				t.Errorf("NeedsStdoutRedirection(%s) = %v, expected %v", tt.binary, result, tt.expected)
			}
		})
	}
}
