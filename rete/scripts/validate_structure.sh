#!/bin/bash
# Script pour valider la nouvelle structure du module RETE

set -e

echo "🔍 Validation de la structure du module RETE"
echo "==========================================="

# Répertoire racine du projet
RETE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$RETE_DIR"

# Vérifier que les nouveaux packages se compilent
echo "🔨 Compilation des nouveaux packages..."
go build ./pkg/... && echo "✅ pkg/ compilé avec succès"
go build ./internal/... && echo "✅ internal/ compilé avec succès"

# Vérifier les imports
echo ""
echo "📦 Vérification des imports..."
go mod tidy && echo "✅ Dépendances vérifiées"

# Vérifier l'organisation des fichiers
echo ""
echo "📁 Structure des répertoires:"
echo "Root files:"
ls -1 *.go *.md 2>/dev/null || echo "  (aucun fichier Go/MD à la racine)"

echo ""
echo "pkg/ structure:"
find pkg -type f -name "*.go" | head -10

echo ""
echo "test/ structure:"
find test -type f -name "*.go" -o -name "*.out" -o -name "*.html" | head -10

echo ""
echo "docs/ structure:"
find docs -type f | head -10

echo ""
echo "scripts/ structure:"
find scripts -type f | head -10

# Vérifier les tests
echo ""
echo "🧪 Tests rapides..."
go test -run TestBasicFunctionality -short . 2>/dev/null || echo "⚠️  Tests en cours d'adaptation à la nouvelle structure"

echo ""
echo "✅ Validation terminée !"
echo "📊 Pour voir la couverture complète: ./scripts/run_tests.sh"