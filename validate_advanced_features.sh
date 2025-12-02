#!/bin/bash

echo "=========================================="
echo "Validation des optimisations avancées"
echo "=========================================="
echo

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Compteurs
PASSED=0
FAILED=0

# Fonction de test
test_step() {
    echo -n "  $1... "
}

test_pass() {
    echo -e "${GREEN}✓${NC}"
    ((PASSED++))
}

test_fail() {
    echo -e "${RED}✗${NC} $1"
    ((FAILED++))
}

echo "1️⃣  Vérification de la compilation"
echo "-----------------------------------"

test_step "Compilation du package rete"
if go build ./rete/... > /dev/null 2>&1; then
    test_pass
else
    test_fail "Erreurs de compilation"
fi

test_step "Compilation des tests"
if go test -c ./test/integration/incremental/ > /dev/null 2>&1; then
    test_pass
else
    test_fail "Erreurs de compilation des tests"
fi

test_step "Compilation de l'exemple"
if go build ./examples/advanced_features_example.go > /dev/null 2>&1; then
    test_pass
else
    test_fail "Erreurs de compilation de l'exemple"
fi

echo
echo "2️⃣  Vérification de la structure des fichiers"
echo "----------------------------------------------"

FILES=(
    "rete/transaction.go"
    "rete/incremental_validation.go"
    "rete/constraint_pipeline_advanced.go"
    "test/integration/incremental/advanced_test.go"
    "examples/advanced_features_example.go"
    "docs/ADVANCED_OPTIMIZATIONS.md"
    "docs/ADVANCED_FEATURES_README.md"
    "docs/ADVANCED_OPTIMIZATIONS_COMPLETION.md"
    "docs/README_OPTIMIZATIONS.md"
)

for file in "${FILES[@]}"; do
    test_step "Fichier $file"
    if [ -f "$file" ]; then
        test_pass
    else
        test_fail "Fichier manquant"
    fi
done

echo
echo "3️⃣  Vérification des symboles exportés"
echo "----------------------------------------"

test_step "Validation incrémentale dans IngestFile"
if grep -q "Validation sémantique incrémentale" rete/constraint_pipeline.go; then
    test_pass
else
    test_fail "Validation incrémentale non trouvée"
fi

test_step "GC dans IngestFile"
if grep -q "GarbageCollect" rete/constraint_pipeline.go; then
    test_pass
else
    test_fail "GC non trouvé"
fi

test_step "IngestFileTransactional"
if grep -q "IngestFileTransactional" rete/constraint_pipeline.go; then
    test_pass
else
    test_fail "Symbole manquant"
fi

test_step "BeginTransaction"
if grep -q "BeginTransaction" rete/transaction.go; then
    test_pass
else
    test_fail "Symbole manquant"
fi

test_step "GarbageCollect"
if grep -q "GarbageCollect" rete/network.go; then
    test_pass
else
    test_fail "Symbole manquant"
fi

echo
echo "=========================================="
echo "Résultat"
echo "=========================================="
echo -e "${GREEN}Tests réussis : $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Tests échoués  : $FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}Tous les tests sont passés ✓${NC}"
    exit 0
fi
