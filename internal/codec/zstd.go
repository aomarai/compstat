package codec

import (
	"fmt"
	"os/exec"
)

type ZstdCodec struct{}

func (z *ZstdCodec) Name() string {
	return "zstd"
}

func (z *ZstdCodec) Binary() string {
	return "zstd"
}

func (z *ZstdCodec) Extension() string {
	return ".zst"
}

func (z *ZstdCodec) Levels() []int {
	return MakeRange(1, 19)
}

func (z *ZstdCodec) SupportsThreading() bool {
	return true
}

func (z *ZstdCodec) IsAvailable() bool {
	_, err := exec.LookPath(z.Binary())
	return err == nil
}

func (z *ZstdCodec) CompressCommand(level, threads int, input, output string) []string {
	args := []string{fmt.Sprintf("-%d", level), fmt.Sprintf("-T%d", threads), "-q", "-f", "-o", output, input}
	if level == 19 {
		args = append([]string{"--ultra", "--long=31"}, args...)
	}
	return args
}

func (z *ZstdCodec) DecompressCommand(threads int, input, output string) []string {
	return []string{"-d", fmt.Sprintf("-T%d", threads), "-q", "-f", "-o", output, input}
}
