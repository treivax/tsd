#!/bin/bash

# Script pour lancer les tests du module constraint

echo "🧪 TESTS DU MODULE CONSTRAINT"
echo "============================="
echo ""

cd "$(dirname "$0")" || exit 1

echo "📁 Répertoire de travail: $(pwd)"
echo "📦 Module testé: constraint"
echo ""

echo "🔍 Vérification des fichiers de test disponibles:"
ls -la tests/*.txt | while read -r line; do
    filename=$(echo "$line" | awk '{print $NF}' | xargs basename)
    if [[ "$filename" == *"mismatch"* ]] || [[ "$filename" == *"error"* ]]; then
        echo "  ❌ $filename (erreur attendue)"
    else
        echo "  ✅ $filename (succès attendu)"
    fi
done
echo ""

echo "🏃 Exécution des tests unitaires..."
echo ""

# Exécuter les tests avec verbose output
go test -v

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 TOUS LES TESTS SONT PASSÉS !"
    echo ""
    
    echo "📊 Exécution des benchmarks..."
    go test -bench=. -benchmem
    
    echo ""
    echo "📈 Coverage des tests..."
    go test -cover
    
else
    echo ""
    echo "❌ ÉCHEC DE CERTAINS TESTS"
    echo "Vérifiez les erreurs ci-dessus"
    exit 1
fi

echo ""
echo "✨ Tests terminés avec succès !"