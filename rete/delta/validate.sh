#!/bin/bash
# Script de validation du package delta
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License

echo "🔍 Validation du package rete/delta"
echo "===================================="

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

PASS=0
FAIL=0

run_check() {
    local name="$1"
    local cmd="$2"
    echo -n "[$name] "
    if eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        ((PASS++))
        return 0
    else
        echo -e "${RED}✗${NC}"
        ((FAIL++))
        return 1
    fi
}

echo ""
echo "📋 Vérifications de code"
run_check "go fmt           " "test -z \"\$(go fmt ./rete/delta/... 2>&1)\""
run_check "goimports        " "test \$(goimports -l ./rete/delta/ 2>&1 | wc -l) -eq 0"
run_check "go vet           " "go vet ./rete/delta/..."
run_check "staticcheck      " "staticcheck ./rete/delta/..."

echo ""
echo "🧪 Tests"
run_check "Tests unitaires  " "go test ./rete/delta/... -short"
run_check "Race detector    " "go test ./rete/delta/... -race -short"

echo ""
echo "📊 Couverture"
COVERAGE=$(go test ./rete/delta/... -cover -short 2>&1 | grep coverage | awk '{print $5}' | sed 's/%//')
echo -n "Couverture: ${COVERAGE}% "
if [ -n "$COVERAGE" ] && (( $(echo "$COVERAGE >= 80.0" | bc -l 2>/dev/null || echo 0) )); then
    echo -e "${GREEN}✓${NC}"
    ((PASS++))
else
    echo -e "${RED}✗${NC}"
    ((FAIL++))
fi

echo ""
echo "===================================="
if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✓ Validation réussie ($PASS tests)${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAIL test(s) échoué(s)${NC}"
    exit 1
fi
