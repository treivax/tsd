#!/bin/bash
# Script de nettoyage en profondeur pour le projet TSD
# Usage: ./scripts/deep_clean.sh

set -e

echo "🧹 NETTOYAGE EN PROFONDEUR - TSD"
echo "================================"

# Répertoire racine du projet
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

echo "📁 Répertoire de travail: $PROJECT_DIR"
echo ""

# 1. Nettoyage des dépendances Go
echo "🔧 1. Nettoyage des dépendances Go..."
go mod tidy
echo "✅ Dépendances Go nettoyées"
echo ""

# 2. Formatage du code
echo "✨ 2. Formatage du code Go..."
go fmt ./...
echo "✅ Code formaté selon les standards Go"
echo ""

# 3. Analyse statique
echo "🔍 3. Analyse statique avec go vet..."
go vet ./...
echo "✅ Analyse statique passée"
echo ""

# 4. Vérification avec staticcheck (si disponible)
if command -v staticcheck &> /dev/null; then
    echo "🔬 4. Analyse avancée avec staticcheck..."
    staticcheck ./...
    echo "✅ Analyse staticcheck passée"
else
    echo "⚠️  4. staticcheck non disponible, installation recommandée:"
    echo "   go install honnef.co/go/tools/cmd/staticcheck@latest"
fi
echo ""

# 5. Compilation complète
echo "🔨 5. Vérification de la compilation..."
go build ./...
echo "✅ Compilation réussie"
echo ""

# 6. Tests rapides
echo "🧪 6. Tests rapides..."
go test -short ./...
echo "✅ Tests rapides passés"
echo ""

# 7. Nettoyage des fichiers temporaires
echo "🗑️  7. Nettoyage des fichiers temporaires..."
find . -name "*.tmp" -delete 2>/dev/null || true
find . -name "*~" -delete 2>/dev/null || true
find . -name "*.bak" -delete 2>/dev/null || true
find . -name ".#*" -delete 2>/dev/null || true
echo "✅ Fichiers temporaires supprimés"
echo ""

# 8. Rapport final
echo "📊 RAPPORT DE NETTOYAGE"
echo "======================"
echo "✅ Formatage Go: OK"
echo "✅ Analyse statique: OK"
echo "✅ Compilation: OK"
echo "✅ Tests rapides: OK"
echo "✅ Dépendances: Optimisées"
echo "✅ Fichiers temporaires: Nettoyés"
echo ""
echo "🎯 **PROJET NETTOYÉ ET OPTIMISÉ**"
echo ""
echo "💡 Prochaines commandes utiles:"
echo "   make test        # Tests complets"
echo "   make coverage    # Couverture de code"
echo "   make build       # Build production"
echo ""
