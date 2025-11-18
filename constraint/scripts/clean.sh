#!/bin/bash
# Script de nettoyage pour le module constraint

echo "🧹 Nettoyage du module constraint"
echo "=================================="

# Répertoire racine du projet
CONSTRAINT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$CONSTRAINT_DIR"

echo "📁 Répertoire: $CONSTRAINT_DIR"
echo ""

# Supprimer les fichiers de couverture temporaires
echo "📊 Suppression des fichiers de couverture temporaires..."
find . -name "*.out" -type f -not -path "./test/coverage/reports/*" -delete
find . -name "*.html" -type f -not -path "./test/coverage/reports/*" -delete
echo "✅ Fichiers de couverture temporaires supprimés"

# Nettoyer le cache Go
echo "🗂️  Nettoyage du cache Go..."
go clean -cache -testcache 2>/dev/null || echo "⚠️  Certains caches nécessitent des privilèges root"

# Supprimer les fichiers de build temporaires
echo "🔨 Suppression des artefacts de build..."
find . -name "debug" -delete 2>/dev/null || true
find . -name "*.exe" -delete 2>/dev/null || true
find . -name "*.tmp" -delete 2>/dev/null || true

# Nettoyer les logs temporaires (s'ils existent)
echo "📝 Nettoyage des logs temporaires..."
find . -name "*.log" -path "./tmp/*" -delete 2>/dev/null || true
find . -name "*.temp" -delete 2>/dev/null || true

# Nettoyer les répertoires vides (sauf structure)
echo "📁 Nettoyage des répertoires vides..."
find . -type d -empty -not -path "./test/*" -not -path "./pkg/*" -not -path "./internal/*" -not -path "./scripts/*" -delete 2>/dev/null || true

# Nettoyer les fichiers de backup des éditeurs
echo "✏️  Suppression des fichiers de backup..."
find . -name "*~" -delete 2>/dev/null || true
find . -name "*.bak" -delete 2>/dev/null || true
find . -name ".#*" -delete 2>/dev/null || true

echo ""
echo "✅ Nettoyage terminé !"
echo ""
echo "📊 Structure préservée:"
echo "  📁 pkg/      - Packages publics"
echo "  📁 internal/ - Packages internes"
echo "  📁 test/     - Tests et rapports"
echo "  📁 scripts/  - Scripts utilitaires"
echo "  📁 docs/     - Documentation"
echo ""
echo "💡 Les rapports finaux dans test/coverage/reports/ sont préservés"
