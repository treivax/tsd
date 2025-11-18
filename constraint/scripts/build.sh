#!/bin/bash
# Script de build pour le module constraint refactorisé

set -e

echo "🔧 Construction du module constraint (nouvelle architecture)..."
echo "============================================================="

# Répertoire racine du projet
CONSTRAINT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$CONSTRAINT_DIR"

echo "📁 Répertoire de travail: $CONSTRAINT_DIR"
echo ""

# Vérifier les dépendances
echo "🔍 Vérification des dépendances..."
if ! command -v pigeon &> /dev/null; then
    echo "❌ Erreur: pigeon n'est pas installé"
    echo "💡 Installez-le avec: go install github.com/mna/pigeon@latest"
    exit 1
fi

echo "✅ pigeon trouvé: $(command -v pigeon)"
echo ""

# Génération du parser depuis la grammaire PEG
echo "📦 Génération du parser depuis la grammaire PEG..."
export PATH=$PATH:~/go/bin
if [ -f "grammar/constraint.peg" ]; then
    pigeon -o parser.go grammar/constraint.peg
    echo "✅ Parser généré avec succès"
else
    echo "❌ Erreur: fichier grammar/constraint.peg introuvable"
    exit 1
fi
echo ""

# Compilation des nouveaux packages
echo "🔨 Compilation des nouveaux packages..."
echo "  🏗️  Compilation pkg/domain..."
go build -v ./pkg/domain
echo "  🏗️  Compilation pkg/validator..."
go build -v ./pkg/validator
echo "  🏗️  Compilation internal/config..."
go build -v ./internal/config
echo "✅ Tous les packages compilés avec succès"
echo ""

# Tests des nouveaux packages
echo "🧪 Tests des nouveaux packages..."
go test ./pkg/... ./internal/... -v
echo "✅ Tests des nouveaux packages réussis"
echo ""

# Compilation du module principal
echo "🔧 Compilation du module principal..."
go build -v .
echo "✅ Module principal compilé"
echo ""

# Construction de l'exécutable
echo "🎯 Construction de l'exécutable..."
cd ..
go build -o constraint-parser ./constraint/cmd/
echo "✅ Exécutable constraint-parser créé"
echo ""

# Test avec un fichier d'exemple
if [ -f "constraint/tests/test_type_valid.txt" ]; then
    echo "🧪 Test avec fichier d'exemple..."
    ./constraint-parser constraint/tests/test_type_valid.txt
    echo "✅ Test d'exemple réussi"
else
    echo "⚠️  Fichier de test non trouvé, test d'exemple ignoré"
fi

echo ""
echo "🎉 Module constraint (nouvelle architecture) construit avec succès !"
echo ""
echo "📊 Structure créée:"
echo "  ├── pkg/domain/     - Types fondamentaux et erreurs structurées"
echo "  ├── pkg/validator/  - Validation et vérification de types"
echo "  ├── internal/config/ - Configuration structurée"
echo "  ├── test/           - Tests organisés"
echo "  └── scripts/        - Scripts utilitaires"
echo ""
echo "💡 Utilisation:"
echo "  ./constraint-parser <fichier.txt>"
echo "  ./scripts/run_tests.sh      # Tests complets"
echo "  ./scripts/validate.sh       # Validation architecture"
