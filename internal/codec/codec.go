package codec

// Codec defines the interface for compression codecs
type Codec interface {
	Name() string
	Binary() string
	Extension() string
	Levels() []int
	IsAvailable() bool
	CompressCommand(level, threads int, input, output string) []string
	DecompressCommand(threads int, input, output string) []string
	SupportsThreading() bool
}

// Registry holds all available codecs
var Registry = map[string]Codec{
	"zstd":   &ZstdCodec{},
	"xz":     &XzCodec{},
	"gzip":   &GzipCodec{},
	"lz4":    &Lz4Codec{},
	"bzip2":  &Bzip2Codec{},
	"brotli": &BrotliCodec{},
}

// MakeRange creates a slice of integers from min to max (inclusive)
func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
