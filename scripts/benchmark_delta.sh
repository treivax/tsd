#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License

set -e

echo "⚡ Benchmarking Delta Propagation System"
echo "========================================"

ITERATIONS=5
OUTPUT_DIR="benchmark_results"

mkdir -p "$OUTPUT_DIR"

echo ""
echo "🏃 Running benchmarks ($ITERATIONS iterations)..."

for i in $(seq 1 $ITERATIONS); do
    echo "  Iteration $i/$ITERATIONS..."
    go test ./rete/delta/... \
        -bench=. \
        -benchmem \
        -benchtime=5s \
        -count=1 \
        > "$OUTPUT_DIR/bench_$i.txt" 2>&1
done

echo ""
echo "📊 Aggregating results..."

# Utiliser benchstat si disponible
if command -v benchstat &> /dev/null; then
    benchstat "$OUTPUT_DIR"/bench_*.txt > "$OUTPUT_DIR/aggregate.txt"
    echo ""
    echo "📈 Benchmark Results (Aggregated):"
    echo "=================================="
    cat "$OUTPUT_DIR/aggregate.txt"
else
    echo "⚠️  benchstat not installed - showing last run only"
    echo "   Install with: go install golang.org/x/perf/cmd/benchstat@latest"
    echo ""
    echo "📈 Benchmark Results (Last Run):"
    echo "================================"
    cat "$OUTPUT_DIR/bench_$ITERATIONS.txt"
fi

echo ""
echo "✅ Benchmarks complete!"
echo ""
echo "📁 Results saved in: $OUTPUT_DIR/"
echo ""
echo "🔍 Analysis:"

# Extraire quelques statistiques clés
if [ -f "$OUTPUT_DIR/aggregate.txt" ]; then
    echo ""
    echo "Top 5 fastest benchmarks:"
    grep "BenchmarkDelta" "$OUTPUT_DIR/aggregate.txt" | sort -k3 -n | head -5 || true
    
    echo ""
    echo "Top 5 most allocations:"
    grep "BenchmarkDelta" "$OUTPUT_DIR/aggregate.txt" | sort -k5 -rn | head -5 || true
fi

echo ""
echo "💡 Tips:"
echo "   - Compare with previous runs to track improvements"
echo "   - Use 'benchstat old.txt new.txt' to compare runs"
echo "   - Focus on ns/op and allocs/op for optimization"
