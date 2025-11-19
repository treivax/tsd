#!/bin/bash

echo "=== TEST EXISTS SIMPLE UNIQUEMENT ==="
echo "Fichier constraint:"
head -10 beta_coverage_tests/exists_simple.constraint
echo ""
echo "Fichier facts:"
head -10 beta_coverage_tests/exists_simple.facts
echo ""
echo "=== RÉSULTAT DU TEST ==="

# Créer un dossier temporaire avec seulement le test EXISTS
mkdir -p test_temp
cp beta_coverage_tests/exists_simple.constraint test_temp/
cp beta_coverage_tests/exists_simple.facts test_temp/

# Exécuter le test sur ce dossier isolé
go run test/coverage/beta/runner.go test_temp/exists_simple.constraint test_temp/exists_simple.facts

# Nettoyer
rm -rf test_temp
