package codec

import (
	"fmt"
	"os/exec"
)

type XzCodec struct{}

func (x *XzCodec) Name() string {
	return "xz"
}

func (x *XzCodec) Binary() string {
	return "xz"
}

func (x *XzCodec) Extension() string {
	return ".xz"
}

func (x *XzCodec) Levels() []int {
	return MakeRange(0, 9)
}

func (x *XzCodec) SupportsThreading() bool {
	return true
}

func (x *XzCodec) IsAvailable() bool {
	_, err := exec.LookPath(x.Binary())
	return err == nil
}

func (x *XzCodec) CompressCommand(level, threads int, input, output string) []string {
	args := []string{fmt.Sprintf("-%d", level), fmt.Sprintf("-T%d", threads), "-c", input}
	if level == 9 {
		args = append([]string{"-e"}, args...)
	}
	return args
}

func (x *XzCodec) DecompressCommand(threads int, input, output string) []string {
	return []string{"-d", fmt.Sprintf("-T%d", threads), "-c", input}
}
