#!/bin/bash

# Script de mise à jour complète des statistiques du projet TSD
# Usage: ./update_stats.sh

set -e

echo "════════════════════════════════════════════════════════════════"
echo "  📊 TSD Code Statistics Update"
echo "════════════════════════════════════════════════════════════════"
echo

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

TIMESTAMP=$(date -Iseconds)
DATE=$(date +%Y-%m-%d)
COMMIT=$(git rev-parse HEAD)
COMMIT_SHORT=$(git rev-parse --short HEAD)
BRANCH=$(git branch --show-current)

echo "📅 Date: $DATE"
echo "🔖 Commit: $COMMIT_SHORT"
echo "🌿 Branch: $BRANCH"
echo

# 1. Exécuter les tests avec couverture
echo -e "${YELLOW}Step 1/5:${NC} Running tests with coverage..."
go test -coverprofile=coverage.out ./... > test_output.txt 2>&1
echo -e "${GREEN}✓${NC} Tests completed"
echo

# 2. Générer le rapport HTML de couverture
echo -e "${YELLOW}Step 2/5:${NC} Generating HTML coverage report..."
go tool cover -html=coverage.out -o docs/reports/coverage_report.html
echo -e "${GREEN}✓${NC} HTML report generated"
echo

# 3. Calculer les métriques
echo -e "${YELLOW}Step 3/5:${NC} Calculating metrics..."

TOTAL_LINES=$(find . -name "*.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l)
TEST_LINES=$(find . -name "*_test.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l)
GENERATED_LINES=$(find . -name "parser.go" -path "*/constraint/*" -exec cat {} \; | wc -l 2>/dev/null || echo 0)
MANUAL_LINES=$((TOTAL_LINES - TEST_LINES - GENERATED_LINES))

TOTAL_FILES=$(find . -name "*.go" -not -path "*/vendor/*" | wc -l)
TEST_FILES=$(find . -name "*_test.go" -not -path "*/vendor/*" | wc -l)
PROD_FILES=$((TOTAL_FILES - TEST_FILES))

GLOBAL_COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $NF}')

echo "  Lines: $TOTAL_LINES total, $MANUAL_LINES manual, $TEST_LINES tests"
echo "  Files: $TOTAL_FILES total, $PROD_FILES prod, $TEST_FILES tests"
echo "  Coverage: $GLOBAL_COVERAGE"
echo -e "${GREEN}✓${NC} Metrics calculated"
echo

# 4. Mettre à jour le JSON
echo -e "${YELLOW}Step 4/5:${NC} Updating JSON metrics..."
./generate_metrics.sh > /dev/null 2>&1
echo -e "${GREEN}✓${NC} JSON updated"
echo

# 5. Afficher un résumé
echo -e "${YELLOW}Step 5/5:${NC} Summary"
echo
echo "┌─────────────────────────────────────────────────────────────┐"
echo "│                    STATISTICS SUMMARY                       │"
echo "├─────────────────────────────────────────────────────────────┤"
printf "│ Total Lines:        %8s                                │\n" "$TOTAL_LINES"
printf "│ Manual Code:        %8s                                │\n" "$MANUAL_LINES"
printf "│ Tests:              %8s                                │\n" "$TEST_LINES"
printf "│ Files:              %8s                                │\n" "$TOTAL_FILES"
printf "│ Coverage:           %8s                                │\n" "$GLOBAL_COVERAGE"
echo "└─────────────────────────────────────────────────────────────┘"
echo

# Afficher les packages à 0%
echo "🔴 Packages at 0% coverage:"
go tool cover -func=coverage.out | grep '0.0%' | grep -v test | awk '{print $1}' | \
    sed 's|github.com/treivax/tsd/||' | sort -u | head -10 | while read pkg; do
    echo "   - $pkg"
done
echo

echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  ✓ Statistics updated successfully!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo
echo "📁 Reports available:"
echo "   - docs/reports/CODE_STATS_2025-11-26.md (detailed)"
echo "   - docs/reports/DASHBOARD.md (visual)"
echo "   - docs/reports/coverage_report.html (interactive)"
echo "   - docs/reports/code_metrics.json (machine-readable)"
echo
echo "💡 Next steps:"
echo "   - Review priority packages in DASHBOARD.md"
echo "   - Add tests for packages at 0% coverage"
echo "   - Update docs/reports/CODE_STATS_YYYY-MM-DD.md if needed"
echo

