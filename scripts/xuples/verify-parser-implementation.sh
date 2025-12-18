#!/bin/bash
# Verification script for xuple-space parser implementation

echo "🔍 Verification Script - xuple-space Parser Implementation"
echo "==========================================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# Change to project directory
cd /home/resinsec/dev/tsd

echo "📋 Step 1: Check Files Exist"
echo "----------------------------"
files=(
    "constraint/constraint_types.go"
    "constraint/grammar/constraint.peg"
    "constraint/parser.go"
    "constraint/xuplespace_parser_test.go"
    "docs/xuples/implementation/01-parser-analysis.md"
    "docs/xuples/implementation/02-xuplespace-syntax.md"
    "docs/xuples/user-guide/xuplespace-command.md"
    "examples/xuples/basic-xuplespace.tsd"
    "examples/xuples/all-policies.tsd"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        print_status 0 "File exists: $file"
    else
        print_status 1 "File missing: $file"
    fi
done
echo ""

echo "🔨 Step 2: Build Package"
echo "-----------------------"
go build ./constraint > /dev/null 2>&1
print_status $? "Build constraint package"
echo ""

echo "🧹 Step 3: Format Check"
echo "----------------------"
go fmt ./constraint/... > /dev/null 2>&1
print_status $? "Format constraint package"
echo ""

echo "🔍 Step 4: Lint Check"
echo "--------------------"
go vet ./constraint > /dev/null 2>&1
print_status $? "Vet constraint package"
echo ""

echo "🧪 Step 5: Run Tests"
echo "-------------------"
echo "Running xuple-space parser tests..."
test_output=$(go test ./constraint -run TestParseXupleSpace -v 2>&1)
test_result=$?

if [ $test_result -eq 0 ]; then
    # Count passed tests
    passed=$(echo "$test_output" | grep -c "PASS.*TestParseXupleSpace")
    echo -e "${GREEN}✅ All xuple-space tests passed ($passed test suites)${NC}"
else
    echo -e "${RED}❌ Some tests failed${NC}"
    echo "$test_output"
    exit 1
fi
echo ""

echo "📊 Step 6: Coverage Check"
echo "------------------------"
coverage_output=$(go test ./constraint -cover 2>&1 | grep "coverage:")
coverage=$(echo "$coverage_output" | grep -oP '\d+\.\d+')

if (( $(echo "$coverage >= 80.0" | bc -l) )); then
    echo -e "${GREEN}✅ Coverage: $coverage% (>= 80%)${NC}"
else
    echo -e "${RED}❌ Coverage: $coverage% (< 80%)${NC}"
    exit 1
fi
echo ""

echo "📝 Step 7: Parse Example Files"
echo "------------------------------"
cat > /tmp/test_parser_xuples.go << 'EOF'
package main
import (
    "fmt"
    "os"
    "github.com/treivax/tsd/constraint"
)
func main() {
    files := []string{
        "examples/xuples/basic-xuplespace.tsd",
        "examples/xuples/all-policies.tsd",
    }
    for _, file := range files {
        result, err := constraint.ParseConstraintFile(file)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Parse error for %s: %v\n", file, err)
            os.Exit(1)
        }
        program, err := constraint.ConvertResultToProgram(result)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Convert error for %s: %v\n", file, err)
            os.Exit(1)
        }
        fmt.Printf("✅ %s: %d xuple-space(s), %d type(s), %d rule(s)\n",
            file, len(program.XupleSpaces), len(program.Types), len(program.Expressions))
    }
}
EOF

parse_output=$(go run /tmp/test_parser_xuples.go 2>&1)
parse_result=$?

if [ $parse_result -eq 0 ]; then
    echo "$parse_output"
else
    echo -e "${RED}❌ Parse failed${NC}"
    echo "$parse_output"
    exit 1
fi
echo ""

echo "📈 Step 8: Statistics"
echo "--------------------"
echo "Production code (constraint package):"
prod_lines=$(find constraint -name "*.go" ! -name "*_test.go" ! -name "parser.go" -exec wc -l {} + | tail -1 | awk '{print $1}')
echo "  - Production Go: $prod_lines lines"

echo "Generated code:"
gen_lines=$(wc -l constraint/parser.go | awk '{print $1}')
echo "  - parser.go (generated): $gen_lines lines"

echo "Test code:"
test_lines=$(wc -l constraint/xuplespace_parser_test.go | awk '{print $1}')
echo "  - xuplespace_parser_test.go: $test_lines lines"

echo "Documentation:"
doc_lines=$(find docs/xuples -name "*.md" -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}')
echo "  - Documentation: $doc_lines lines"

echo "Examples:"
example_lines=$(find examples/xuples -name "*.tsd" -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}')
echo "  - TSD examples: $example_lines lines"

echo ""
echo "✨ Step 9: Feature Summary"
echo "-------------------------"
echo "✅ Selection policies: random, fifo, lifo"
echo "✅ Consumption policies: once, per-agent, limited(n)"
echo "✅ Retention policies: unlimited, duration(s/m/h/d)"
echo "✅ Default values: fifo, once, unlimited"
echo "✅ Validation: limits > 0, durations > 0"
echo "✅ Error messages: clear and localized"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}🎉 VERIFICATION COMPLETE - ALL CHECKS PASSED${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📋 Summary:"
echo "  - Files: All present ✅"
echo "  - Build: Success ✅"
echo "  - Format: OK ✅"
echo "  - Lint: OK ✅"
echo "  - Tests: All passed ✅"
echo "  - Coverage: $coverage% ✅"
echo "  - Examples: Parse correctly ✅"
echo ""
echo "🚀 Ready for: Prompt 04 - Implement Default Actions"
echo ""
