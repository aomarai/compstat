package codec

import (
	"fmt"
	"os/exec"
)

type Lz4Codec struct{}

func (l *Lz4Codec) Name() string {
	return "lz4"
}

func (l *Lz4Codec) Binary() string {
	return "lz4"
}

func (l *Lz4Codec) Extension() string {
	return ".lz4"
}

func (l *Lz4Codec) Levels() []int {
	return MakeRange(1, 9)
}

func (l *Lz4Codec) SupportsThreading() bool {
	return false
}

func (l *Lz4Codec) IsAvailable() bool {
	_, err := exec.LookPath(l.Binary())
	return err == nil
}

func (l *Lz4Codec) CompressCommand(level, threads int, input, output string) []string {
	return []string{fmt.Sprintf("-%d", level), input, output}
}

func (l *Lz4Codec) DecompressCommand(threads int, input, output string) []string {
	return []string{"-d", input, output}
}
