#!/bin/bash

# Script de nettoyage pour TSD
# Supprime les fichiers temporaires et artefacts de build

echo "🧹 TSD CLEANUP"
echo "=============="

# Supprimer les binaires
echo "🗑️ Nettoyage des binaires..."
rm -rf bin/
rm -f cmd/main
rm -f constraint/cmd/main
rm -f constraint/cmd/constraint-parser

# Supprimer les fichiers de couverture Go
echo "🗑️ Nettoyage des fichiers de couverture..."
find . -name "*.out" -type f -delete
find . -name "coverage.out" -type f -delete
find . -name "profile.out" -type f -delete

# Supprimer les logs temporaires
echo "🗑️ Nettoyage des logs temporaires..."
find . -name "*.log" -type f -delete
find . -name "*.tmp" -type f -delete

# Supprimer les fichiers de cache Go
echo "🗑️ Nettoyage du cache Go..."
go clean -cache -modcache -testcache 2>/dev/null || true

# Supprimer les fichiers système
echo "🗑️ Nettoyage des fichiers système..."
find . -name ".DS_Store" -type f -delete 2>/dev/null || true
find . -name "Thumbs.db" -type f -delete 2>/dev/null || true

# Supprimer les répertoires de build temporaires
echo "🗑️ Nettoyage des répertoires temporaires..."
rm -rf tmp/
rm -rf temp/
rm -rf build/

echo "✅ Nettoyage terminé !"
echo ""
echo "📁 Structure propre maintenue :"
echo "  ✅ Code source préservé"
echo "  ✅ Tests préservés" 
echo "  ✅ Documentation préservée"
echo "  ✅ Configuration préservée"