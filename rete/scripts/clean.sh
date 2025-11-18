#!/bin/bash
# Script pour nettoyer les artefacts de compilation et de test

echo "🧹 Nettoyage du module RETE"
echo "==========================="

# Répertoire racine du projet
RETE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$RETE_DIR"

# Supprimer les fichiers de couverture temporaires
echo "📊 Suppression des fichiers de couverture temporaires..."
find . -name "*.out" -type f -delete
find . -name "*.html" -type f -not -path "./test/coverage/reports/*" -delete

# Nettoyer le cache Go
echo "🗂️  Nettoyage du cache Go..."
go clean -cache -testcache -modcache 2>/dev/null || echo "⚠️  Certains caches nécessitent des privilèges root"

# Supprimer les fichiers de build temporaires
echo "🔨 Suppression des artefacts de build..."
find . -name "*.exe" -delete 2>/dev/null || true
find . -name "debug" -delete 2>/dev/null || true

# Nettoyer les répertoires vides
echo "📁 Nettoyage des répertoires vides..."
find . -type d -empty -delete 2>/dev/null || true

echo "✅ Nettoyage terminé !"
