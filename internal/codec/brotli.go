package codec

import (
	"os/exec"
	"strconv"
)

type BrotliCodec struct{}

func (br *BrotliCodec) Name() string {
	return "brotli"
}

func (br *BrotliCodec) Binary() string {
	return "brotli"
}

func (br *BrotliCodec) Extension() string {
	return ".br"
}

func (br *BrotliCodec) Levels() []int {
	return MakeRange(1, 11)
}

func (br *BrotliCodec) SupportsThreading() bool {
	return true
}

func (br *BrotliCodec) IsAvailable() bool {
	_, err := exec.LookPath(br.Binary())
	return err == nil
}

func (br *BrotliCodec) CompressCommand(level, threads int, input, output string) []string {
	return []string{"-q", strconv.Itoa(level), "-j", strconv.Itoa(threads), "-f", "-o", output, input}
}

func (br *BrotliCodec) DecompressCommand(threads int, input, output string) []string {
	return []string{"-d", "-j", strconv.Itoa(threads), "-f", "-o", output, input}
}
