package codec

import (
	"fmt"
	"os/exec"
	"strconv"
)

type GzipCodec struct{}

func (g *GzipCodec) Name() string {
	return "gzip"
}

func (g *GzipCodec) Binary() string {
	return "pigz"
}

func (g *GzipCodec) Extension() string {
	return ".gz"
}

func (g *GzipCodec) Levels() []int {
	return MakeRange(1, 9)
}

func (g *GzipCodec) SupportsThreading() bool {
	return true
}

func (g *GzipCodec) IsAvailable() bool {
	_, err := exec.LookPath(g.Binary())
	return err == nil
}

func (g *GzipCodec) CompressCommand(level, threads int, input, output string) []string {
	return []string{fmt.Sprintf("-%d", level), "-p", strconv.Itoa(threads), "-c", input}
}

func (g *GzipCodec) DecompressCommand(threads int, input, output string) []string {
	return []string{"-d", "-p", strconv.Itoa(threads), "-c", input}
}
