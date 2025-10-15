#!/bin/bash
# Example benchmark script

set -e

# Configuration
FILES="data.tar"
CODECS="zstd,xz,gzip,lz4,bzip2,brotli"
ITERATIONS=3
PARALLELISM=4
OUTPUT_DIR="results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Check if compstat exists
if ! command -v compstat &> /dev/null && [ ! -f "./compstat" ]; then
    echo "Error: compstat not found"
    echo "Build it with: make build"
    exit 1
fi

COMPSTAT="./compstat"
if command -v compstat &> /dev/null; then
    COMPSTAT="compstat"
fi

echo "======================================"
echo "Compression Benchmark"
echo "======================================"
echo "Files:       $FILES"
echo "Codecs:      $CODECS"
echo "Iterations:  $ITERATIONS"
echo "Parallelism: $PARALLELISM"
echo "Output:      $OUTPUT_DIR"
echo "======================================"
echo ""

# Run benchmark
$COMPSTAT \
    -files "$FILES" \
    -codecs "$CODECS" \
    -iterations "$ITERATIONS" \
    -parallelism "$PARALLELISM" \
    -compress-threads "$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)" \
    -decompress-threads "$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)" \
    -output "$OUTPUT_DIR/benchmark_${TIMESTAMP}.csv" \
    -json "$OUTPUT_DIR/benchmark_${TIMESTAMP}.json"

echo ""
echo "======================================"
echo "Generating Analysis"
echo "======================================"

# Check if Python tools are available
if command -v python3 &> /dev/null; then
    # Generate analysis
    if [ -f "python/analyze.py" ]; then
        echo "Running statistical analysis..."
        python3 python/analyze.py \
            "$OUTPUT_DIR/benchmark_${TIMESTAMP}.csv" \
            --summary \
            --pareto \
            --output "$OUTPUT_DIR/analysis_${TIMESTAMP}.csv"
    fi

    # Generate visualizations
    if [ -f "python/visualize.py" ]; then
        echo "Generating visualizations..."
        python3 python/visualize.py \
            "$OUTPUT_DIR/benchmark_${TIMESTAMP}.csv" \
            --output-dir "$OUTPUT_DIR/plots_${TIMESTAMP}"
    fi

    echo ""
    echo "✓ Analysis complete!"
    echo ""
    echo "Results:"
    echo "  - CSV:   $OUTPUT_DIR/benchmark_${TIMESTAMP}.csv"
    echo "  - JSON:  $OUTPUT_DIR/benchmark_${TIMESTAMP}.json"
    echo "  - Analysis: $OUTPUT_DIR/analysis_${TIMESTAMP}.csv"
    echo "  - Plots: $OUTPUT_DIR/plots_${TIMESTAMP}/"
else
    echo "Python3 not found. Skipping analysis and visualization."
    echo "Install Python dependencies: pip install -r python/requirements.txt"
fi

echo ""
echo "======================================"
echo "Top Results"
echo "======================================"
echo ""
echo "Best Compression Ratio:"
if command -v python3 &> /dev/null; then
    python3 -c "
import pandas as pd
df = pd.read_csv('$OUTPUT_DIR/benchmark_${TIMESTAMP}.csv')
top = df.nsmallest(5, 'compression_ratio')[['algorithm', 'level', 'compression_ratio', 'compression_speed_mbs']]
print(top.to_string(index=False))
"
fi

echo ""
echo "Fastest Compression:"
if command -v python3 &> /dev/null; then
    python3 -c "
import pandas as pd
df = pd.read_csv('$OUTPUT_DIR/benchmark_${TIMESTAMP}.csv')
top = df.nlargest(5, 'compression_speed_mbs')[['algorithm', 'level', 'compression_speed_mbs', 'compression_ratio']]
print(top.to_string(index=False))
"
fi

echo ""
echo "======================================"
echo "Benchmark Complete!"
echo "======================================"