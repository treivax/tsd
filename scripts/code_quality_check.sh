#!/bin/bash
# Script d'analyse de la qualité du code pour le projet TSD
# Usage: ./scripts/code_quality_check.sh

set -e

echo "📊 ANALYSE QUALITÉ CODE - TSD"
echo "============================="

# Répertoire racine du projet
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

echo "📁 Répertoire de travail: $PROJECT_DIR"
echo ""

# Compteurs
TOTAL_ISSUES=0
WARNINGS=0
ERRORS=0

echo "🔍 1. ANALYSE STRUCTURELLE"
echo "========================="

# Analyser la taille des fichiers
echo "📋 Fichiers volumineux (>500 lignes):"
large_files=$(find . -name "*.go" -not -path "*/vendor/*" -exec wc -l {} + | awk '$1 > 500' | head -10)
if [ -n "$large_files" ]; then
    echo "$large_files"
    WARNINGS=$((WARNINGS + $(echo "$large_files" | wc -l)))
else
    echo "✅ Aucun fichier excessivement volumineux détecté"
fi
echo ""

# Analyser les fonctions complexes
echo "🔧 Fonctions avec beaucoup de paramètres (>5):"
complex_funcs=$(grep -rn "^func.*(" . --include="*.go" | grep -v "_test.go" | awk -F'[(),]' 'NF > 7' | head -5)
if [ -n "$complex_funcs" ]; then
    echo "$complex_funcs"
    WARNINGS=$((WARNINGS + $(echo "$complex_funcs" | wc -l)))
else
    echo "✅ Pas de fonctions avec trop de paramètres"
fi
echo ""

echo "🔬 2. ANALYSE STATIQUE"
echo "====================="

# Go vet
echo "🔍 Analyse go vet..."
if go vet ./... 2>&1; then
    echo "✅ go vet: Aucun problème détecté"
else
    echo "❌ go vet: Problèmes détectés"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# Staticcheck (si disponible)
if command -v staticcheck &> /dev/null; then
    echo "🔬 Analyse staticcheck..."
    staticcheck_output=$(staticcheck ./... 2>&1 || true)
    if [ -z "$staticcheck_output" ]; then
        echo "✅ staticcheck: Aucun problème détecté"
    else
        echo "⚠️  staticcheck: Problèmes détectés:"
        echo "$staticcheck_output" | head -10
        WARNINGS=$((WARNINGS + $(echo "$staticcheck_output" | wc -l)))
    fi
else
    echo "⚠️  staticcheck non installé"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

echo "📦 3. ANALYSE DES DÉPENDANCES"
echo "============================"

# Vérifier go.mod
echo "🔍 État de go.mod..."
go mod tidy
if git diff --quiet go.mod go.sum 2>/dev/null; then
    echo "✅ go.mod: Aucune modification nécessaire"
else
    echo "⚠️  go.mod: Modifications détectées après tidy"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# Analyser les dépendances indirectes
indirect_deps=$(go list -m all | grep "// indirect" | wc -l)
echo "📊 Dépendances indirectes: $indirect_deps"
if [ "$indirect_deps" -gt 20 ]; then
    echo "⚠️  Beaucoup de dépendances indirectes"
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ Nombre de dépendances indirectes acceptable"
fi
echo ""

echo "🧪 4. ANALYSE DES TESTS"
echo "====================="

# Couverture de tests
echo "🎯 Vérification de la couverture..."
coverage_output=$(go test -short -coverprofile=/tmp/coverage.out ./... 2>&1 || true)
if echo "$coverage_output" | grep -q "coverage:"; then
    coverage=$(echo "$coverage_output" | grep "coverage:" | tail -1 | awk '{print $5}' | sed 's/%//')
    echo "📊 Couverture globale: ${coverage}%"
    if (( $(echo "$coverage < 70" | bc -l) )); then
        echo "⚠️  Couverture inférieure à 70%"
        WARNINGS=$((WARNINGS + 1))
    else
        echo "✅ Couverture satisfaisante"
    fi
else
    echo "⚠️  Impossible de calculer la couverture"
    WARNINGS=$((WARNINGS + 1))
fi
rm -f /tmp/coverage.out
echo ""

echo "📈 5. MÉTRIQUES PROJET"
echo "====================="

# Statistiques générales
total_go_files=$(find . -name "*.go" -not -path "*/vendor/*" | wc -l)
total_lines=$(find . -name "*.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l)
total_functions=$(grep -r "^func " . --include="*.go" | wc -l)
total_types=$(grep -r "^type.*struct" . --include="*.go" | wc -l)

echo "📋 Statistiques du projet:"
echo "   📄 Fichiers Go: $total_go_files"
echo "   📏 Lignes de code: $total_lines"
echo "   🔧 Fonctions: $total_functions"
echo "   📦 Types (structs): $total_types"
echo ""

# Calcul du score qualité
TOTAL_ISSUES=$((ERRORS * 3 + WARNINGS))
if [ $TOTAL_ISSUES -eq 0 ]; then
    QUALITY_SCORE="A+"
    QUALITY_COLOR="🟢"
elif [ $TOTAL_ISSUES -le 5 ]; then
    QUALITY_SCORE="A"
    QUALITY_COLOR="🟢"
elif [ $TOTAL_ISSUES -le 10 ]; then
    QUALITY_SCORE="B"
    QUALITY_COLOR="🟡"
elif [ $TOTAL_ISSUES -le 20 ]; then
    QUALITY_SCORE="C"
    QUALITY_COLOR="🟠"
else
    QUALITY_SCORE="D"
    QUALITY_COLOR="🔴"
fi

echo "🎯 RAPPORT FINAL"
echo "==============="
echo "📊 Score Qualité: $QUALITY_COLOR $QUALITY_SCORE"
echo "⚠️  Avertissements: $WARNINGS"
echo "❌ Erreurs: $ERRORS"
echo "📈 Points à améliorer: $TOTAL_ISSUES"
echo ""

if [ $TOTAL_ISSUES -eq 0 ]; then
    echo "🎉 **EXCELLENT! Code de haute qualité.**"
elif [ $TOTAL_ISSUES -le 10 ]; then
    echo "✨ **TRÈS BIEN! Quelques améliorations mineures possibles.**"
else
    echo "🔧 **AMÉLIORATIONS RECOMMANDÉES pour optimiser la qualité.**"
fi

echo ""
echo "💡 Recommandations:"
echo "   1. Exécuter régulièrement ./scripts/deep_clean.sh"
echo "   2. Maintenir la couverture de tests > 80%"
echo "   3. Éviter les fichiers > 1000 lignes"
echo "   4. Limiter les fonctions à 5 paramètres max"
echo ""
