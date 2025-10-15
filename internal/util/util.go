package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// ComputeFileHash computes the SHA256 hash of a file
func ComputeFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("warning: failed to close file: %v\n", err)
		}
	}(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// FileSize returns the size of a file in bytes
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// RunCommand executes a command and returns the elapsed time
func RunCommand(binary string, args []string, outputFile string) (time.Duration, error) {
	start := time.Now()
	cmd := exec.Command(binary, args...)

	// Handle stdout redirection for codecs that need it
	if outputFile != "" && NeedsStdoutRedirection(binary) {
		out, err := os.Create(outputFile)
		if err != nil {
			return 0, err
		}
		defer out.Close()
		cmd.Stdout = out
	}

	cmd.Stderr = nil // Discard stderr

	err := cmd.Run()
	elapsed := time.Since(start)

	return elapsed, err
}

// NeedsStdoutRedirection returns true if the binary writes to stdout with -c flag
func NeedsStdoutRedirection(binary string) bool {
	return binary == "xz" || binary == "pigz" || binary == "pbzip2"
}
