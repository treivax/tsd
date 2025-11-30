#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
#
# Profile multi-source aggregation performance
# This script runs benchmarks with CPU and memory profiling enabled

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RETE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROFILE_DIR="$RETE_DIR/profiles"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Multi-Source Aggregation Profiler${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Create profiles directory
mkdir -p "$PROFILE_DIR"
echo -e "${GREEN}✓${NC} Profile directory: $PROFILE_DIR"
echo ""

# Change to rete directory
cd "$RETE_DIR"

# Parse command line arguments
BENCHMARK_FILTER="${1:-BenchmarkMultiSourceAggregation}"
RUN_COUNT="${2:-3}"
BENCHMARK_TIME="${3:-10s}"

echo -e "${YELLOW}Configuration:${NC}"
echo -e "  Benchmark Filter: $BENCHMARK_FILTER"
echo -e "  Run Count: $RUN_COUNT"
echo -e "  Benchmark Time: $BENCHMARK_TIME"
echo ""

# Function to run a specific benchmark with profiling
run_profiled_benchmark() {
    local bench_name=$1
    local profile_suffix=$2

    echo -e "${BLUE}Running: $bench_name${NC}"

    local cpu_profile="$PROFILE_DIR/cpu_${profile_suffix}.prof"
    local mem_profile="$PROFILE_DIR/mem_${profile_suffix}.prof"
    local bench_output="$PROFILE_DIR/bench_${profile_suffix}.txt"

    go test -bench="^${bench_name}$" \
        -benchtime="$BENCHMARK_TIME" \
        -count="$RUN_COUNT" \
        -cpuprofile="$cpu_profile" \
        -memprofile="$mem_profile" \
        -benchmem \
        2>&1 | tee "$bench_output"

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} CPU profile: $cpu_profile"
        echo -e "${GREEN}✓${NC} Memory profile: $mem_profile"
        echo -e "${GREEN}✓${NC} Benchmark output: $bench_output"
    else
        echo -e "${RED}✗${NC} Benchmark failed"
        return 1
    fi
    echo ""
}

# Run comprehensive profiling suite
echo -e "${YELLOW}======================================${NC}"
echo -e "${YELLOW}Running Comprehensive Profile Suite${NC}"
echo -e "${YELLOW}======================================${NC}"
echo ""

# Small scale - 2 sources
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_TwoSources_SmallScale" \
    "two_sources_small"

# Medium scale - 2 sources
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_TwoSources_MediumScale" \
    "two_sources_medium"

# Large scale - 2 sources
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_TwoSources_LargeScale" \
    "two_sources_large"

# Small scale - 3 sources
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_ThreeSources_SmallScale" \
    "three_sources_small"

# Medium scale - 3 sources
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_ThreeSources_MediumScale" \
    "three_sources_medium"

# High fanout
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_HighFanout" \
    "high_fanout"

# Low fanout
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_LowFanout" \
    "low_fanout"

# Many aggregates
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_ManyAggregates" \
    "many_aggregates"

# With thresholds
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_WithThresholds" \
    "with_thresholds"

# Retraction
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_Retraction" \
    "retraction"

# Incremental
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_IncrementalUpdate" \
    "incremental"

# Memory benchmarks
run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_Memory_SmallScale" \
    "memory_small"

run_profiled_benchmark \
    "BenchmarkMultiSourceAggregation_Memory_LargeScale" \
    "memory_large"

echo -e "${YELLOW}======================================${NC}"
echo -e "${YELLOW}Profile Analysis${NC}"
echo -e "${YELLOW}======================================${NC}"
echo ""

# Generate profile analysis reports
echo -e "${BLUE}Analyzing profiles...${NC}"
echo ""

# CPU profile analysis for the large scale benchmark
echo -e "${GREEN}Top 10 CPU consumers (2-source large scale):${NC}"
go tool pprof -top -nodecount=10 "$PROFILE_DIR/cpu_two_sources_large.prof" 2>/dev/null | head -20 || true
echo ""

echo -e "${GREEN}Top 10 CPU consumers (3-source medium scale):${NC}"
go tool pprof -top -nodecount=10 "$PROFILE_DIR/cpu_three_sources_medium.prof" 2>/dev/null | head -20 || true
echo ""

# Memory profile analysis
echo -e "${GREEN}Top 10 memory allocators (large scale):${NC}"
go tool pprof -top -nodecount=10 "$PROFILE_DIR/mem_two_sources_large.prof" 2>/dev/null | head -20 || true
echo ""

# Generate summary report
SUMMARY_FILE="$PROFILE_DIR/summary_report.txt"
echo -e "${BLUE}Generating summary report: $SUMMARY_FILE${NC}"

cat > "$SUMMARY_FILE" << EOF
Multi-Source Aggregation Performance Profile Summary
Generated: $(date)
==========================================

BENCHMARK CONFIGURATIONS:
- Run Count: $RUN_COUNT
- Benchmark Time: $BENCHMARK_TIME
- Filter: $BENCHMARK_FILTER

PROFILED BENCHMARKS:
EOF

for bench_file in "$PROFILE_DIR"/bench_*.txt; do
    if [ -f "$bench_file" ]; then
        echo "" >> "$SUMMARY_FILE"
        echo "--- $(basename "$bench_file" .txt) ---" >> "$SUMMARY_FILE"
        grep -E "^Benchmark|ns/op|allocs/op|B/op|facts/sec|activations" "$bench_file" >> "$SUMMARY_FILE" 2>/dev/null || true
    fi
done

echo "" >> "$SUMMARY_FILE"
echo "==========================================" >> "$SUMMARY_FILE"
echo "CPU PROFILE ANALYSIS (2-source large scale):" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
go tool pprof -top -nodecount=20 "$PROFILE_DIR/cpu_two_sources_large.prof" 2>/dev/null >> "$SUMMARY_FILE" || true

echo "" >> "$SUMMARY_FILE"
echo "==========================================" >> "$SUMMARY_FILE"
echo "MEMORY PROFILE ANALYSIS (large scale):" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
go tool pprof -top -nodecount=20 "$PROFILE_DIR/mem_two_sources_large.prof" 2>/dev/null >> "$SUMMARY_FILE" || true

echo -e "${GREEN}✓${NC} Summary report saved"
echo ""

# Instructions for further analysis
echo -e "${YELLOW}======================================${NC}"
echo -e "${YELLOW}Next Steps${NC}"
echo -e "${YELLOW}======================================${NC}"
echo ""
echo -e "Profile files are available in: ${GREEN}$PROFILE_DIR${NC}"
echo ""
echo -e "To analyze profiles interactively:"
echo -e "  ${BLUE}go tool pprof -http=:8080 $PROFILE_DIR/cpu_two_sources_large.prof${NC}"
echo -e "  ${BLUE}go tool pprof -http=:8080 $PROFILE_DIR/mem_two_sources_large.prof${NC}"
echo ""
echo -e "To generate flamegraphs (if go-torch is installed):"
echo -e "  ${BLUE}go-torch --file=$PROFILE_DIR/cpu_flame.svg $PROFILE_DIR/cpu_two_sources_large.prof${NC}"
echo ""
echo -e "To view specific profiles:"
echo -e "  ${BLUE}go tool pprof $PROFILE_DIR/cpu_<name>.prof${NC}"
echo -e "  ${BLUE}go tool pprof $PROFILE_DIR/mem_<name>.prof${NC}"
echo ""
echo -e "Common pprof commands:"
echo -e "  ${BLUE}top${NC}       - Show top functions by CPU/memory"
echo -e "  ${BLUE}list <func>${NC} - Show source code for function"
echo -e "  ${BLUE}web${NC}       - Generate and view graph (requires graphviz)"
echo -e "  ${BLUE}peek <func>${NC} - Show callers and callees"
echo ""

echo -e "${GREEN}✓${NC} Profiling complete!"
echo ""
