#!/bin/bash

# Script de build pour le module constraint

echo "🔧 Construction du module constraint..."

# On est déjà dans constraint/ donc pas besoin de cd

echo "📦 Génération du parser depuis la grammaire PEG..."
# Régénérer le parser à partir de la grammaire PEG
export PATH=$PATH:~/go/bin
if command -v pigeon &> /dev/null; then
    pigeon -o parser.go grammar/constraint.peg
    echo "✅ Parser généré avec succès"
else
    echo "❌ Erreur: pigeon n'est pas installé. Installez-le avec: go install github.com/mna/pigeon@latest"
    exit 1
fi

echo "🧪 Tests du module constraint..."
# Construire l'exécutable depuis la racine du projet
cd ..
go build -o constraint-parser ./constraint/cmd/

if [ $? -eq 0 ]; then
    echo "✅ Build réussi"
    echo "🎯 Test avec un fichier d'exemple..."
    ./constraint-parser constraint/tests/test_type_valid.txt
else
    echo "❌ Échec du build"
    exit 1
fi

echo "🎉 Module constraint construit avec succès !"
echo "💡 Utilisation: ./constraint-parser constraint/tests/test_type_valid.txt"