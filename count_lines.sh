#!/bin/bash

echo "=== Analyse des lignes de code ==="
echo

# Total lignes Go (hors vendor)
total_lines=$(find . -name "*.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l)
echo "Total lignes .go (hors vendor): $total_lines"

# Lignes de tests
test_lines=$(find . -name "*_test.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l)
echo "Lignes de tests (*_test.go): $test_lines"

# Lignes de code généré (parser.go)
generated_lines=$(find . -name "parser.go" -path "*/constraint/*" -exec cat {} \; | wc -l)
echo "Lignes de code généré (constraint/parser.go): $generated_lines"

# Lignes de code manuel (hors tests et généré)
manual_lines=$((total_lines - test_lines - generated_lines))
echo "Lignes de code manuel (hors tests et généré): $manual_lines"

echo
echo "=== Fichiers par catégorie ==="
# Nombre de fichiers
total_files=$(find . -name "*.go" -not -path "*/vendor/*" | wc -l)
test_files=$(find . -name "*_test.go" -not -path "*/vendor/*" | wc -l)
prod_files=$((total_files - test_files))
echo "Fichiers Go total: $total_files"
echo "Fichiers de tests: $test_files"
echo "Fichiers de production: $prod_files"

echo
echo "=== Top 10 fichiers les plus volumineux (hors tests) ==="
find . -name "*.go" -not -name "*_test.go" -not -path "*/vendor/*" -exec wc -l {} \; | sort -rn | head -10

echo
echo "=== Top 10 fichiers de tests les plus volumineux ==="
find . -name "*_test.go" -not -path "*/vendor/*" -exec wc -l {} \; | sort -rn | head -10
