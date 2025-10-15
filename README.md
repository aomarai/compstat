# Compression Benchmark Tool

A high-performance, extensible tool for benchmarking compression algorithms.

## Features

- Native parallelism with Go
- Plugin-based codec system
- Comprehensive metrics (ratio, speed, memory)
- Decompression verification
- Python analysis tools
- Interactive dashboard

## Quick Start

### Build
```bash
make build
```

### Run Benchmark
```bash
./compstat -files data.tar -codecs zstd,xz,gzip -parallelism 4
```

### Analyze Results
```bash
python python/analyze.py benchmark_results.csv --summary
python python/visualize.py benchmark_results.csv
streamlit run python/dashboard.py
```

## Supported Codecs

- zstd
- xz
- gzip (pigz)
- lz4
- bzip2 (pbzip2)
- brotli

## Adding New Codecs

Implement the `Codec` interface in Go:
```go
type MyCodec struct{}

func (m MyCodec) Name() string { return "mycodec" }
// ... implement other methods
```

Register in `codecRegistry`.

## CI/CD

GitHub Actions automatically:
- Builds for all platforms
- Runs tests and benchmarks
- Creates releases with binaries