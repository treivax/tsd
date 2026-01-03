#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License

set -e

echo "🔍 Profiling Delta Propagation System"
echo "======================================"

PROFILE_DIR="profile_results"
mkdir -p "$PROFILE_DIR"

# 1. CPU Profiling
echo ""
echo "📊 CPU Profiling..."
go test ./rete/delta/... -bench=. -cpuprofile="$PROFILE_DIR/cpu.prof" -benchtime=10s > "$PROFILE_DIR/bench.txt"

echo "   CPU profile saved to: $PROFILE_DIR/cpu.prof"
echo "   View with: go tool pprof -http=:8080 $PROFILE_DIR/cpu.prof"

# 2. Memory Profiling
echo ""
echo "💾 Memory Profiling..."
go test ./rete/delta/... -bench=. -memprofile="$PROFILE_DIR/mem.prof" -benchtime=10s >> "$PROFILE_DIR/bench.txt"

echo "   Memory profile saved to: $PROFILE_DIR/mem.prof"
echo "   View with: go tool pprof -http=:8081 $PROFILE_DIR/mem.prof"

# 3. Allocation Profiling
echo ""
echo "🔢 Allocation Analysis..."
go test ./rete/delta/... -bench=. -benchmem -benchtime=5s > "$PROFILE_DIR/alloc_report.txt"

echo "   Allocation report: $PROFILE_DIR/alloc_report.txt"

# 4. Trace Analysis
echo ""
echo "🔬 Trace Analysis..."
go test ./rete/delta/... -run=^$ -bench=BenchmarkDeltaDetector_DetectDelta_SingleChange -trace="$PROFILE_DIR/trace.out" -benchtime=3s

echo "   Trace saved to: $PROFILE_DIR/trace.out"
echo "   View with: go tool trace $PROFILE_DIR/trace.out"

# 5. Escape Analysis
echo ""
echo "🏃 Escape Analysis..."
cd rete/delta
go build -gcflags='-m -m' . 2> "../../$PROFILE_DIR/escape_analysis.txt"
cd ../..

echo "   Escape analysis: $PROFILE_DIR/escape_analysis.txt"

echo ""
echo "✅ Profiling complete!"
echo ""
echo "📁 Results saved in: $PROFILE_DIR/"
echo ""
echo "🌐 To view profiles:"
echo "   CPU:    go tool pprof -http=:8080 $PROFILE_DIR/cpu.prof"
echo "   Memory: go tool pprof -http=:8081 $PROFILE_DIR/mem.prof"
echo "   Trace:  go tool trace $PROFILE_DIR/trace.out"
