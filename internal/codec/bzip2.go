package codec

import (
	"fmt"
	"os/exec"
)

type Bzip2Codec struct{}

func (b *Bzip2Codec) Name() string {
	return "bzip2"
}

func (b *Bzip2Codec) Binary() string {
	return "pbzip2"
}

func (b *Bzip2Codec) Extension() string {
	return ".bz2"
}

func (b *Bzip2Codec) Levels() []int {
	return MakeRange(1, 9)
}

func (b *Bzip2Codec) SupportsThreading() bool {
	return true
}

func (b *Bzip2Codec) IsAvailable() bool {
	_, err := exec.LookPath(b.Binary())
	return err == nil
}

func (b *Bzip2Codec) CompressCommand(level, threads int, input, output string) []string {
	return []string{fmt.Sprintf("-%d", level), fmt.Sprintf("-p%d", threads), "-c", input}
}

func (b *Bzip2Codec) DecompressCommand(threads int, input, output string) []string {
	return []string{"-d", fmt.Sprintf("-p%d", threads), "-c", input}
}
