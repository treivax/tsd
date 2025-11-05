#!/usr/bin/env bash

echo "=== Structure finale des fichiers ==="
echo ""

echo "📁 FICHIERS PARSER PEG GO (nécessaires):"
echo "✅ constraint.peg - Grammaire PEG principale"
echo "✅ constraint_types.go - Structures Go pour l'AST"
echo "✅ constraint_utils.go - Fonctions utilitaires et validation"
echo "✅ constraint_main.go - Programme principal"
echo "✅ build.sh - Script de construction"
echo "✅ test_input.txt - Exemple d'entrée"
echo "✅ PARSER_README.md - Documentation du parser"
echo "✅ go.mod - Dépendances Go"
echo ""

echo "📁 AUTRES FICHIERS (projets séparés):"
echo "ℹ️  SetConstraint.g4 - Grammaire ANTLR (projet séparé)"
echo "ℹ️  main.go, operations.go, take.go, put.go - Projet etcd (séparé)"
echo "ℹ️  README.md - Documentation du projet etcd"
echo ""

echo "🗑️ FICHIERS SUPPRIMÉS (étaient inutiles):"
echo "❌ constraint_parser.py - Version Python (erreur)"
echo "❌ SetConstraint.pegjs - Version JavaScript (remplacée)"
echo ""

echo "=== Structure cohérente et prête à l'emploi ! ==="