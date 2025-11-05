#!/bin/bash

# Script de build principal pour le projet TSD
set -e

echo "🚀 Build du projet TSD"
echo "======================"

# Build du module constraint
echo "📦 Construction du module constraint..."
cd constraint && ./build.sh && cd ..

echo ""
echo "🔧 Construction du client etcd..."
go mod tidy
go build -o etcd-client main.go operations.go put.go take.go

if [ $? -eq 0 ]; then
    echo "✅ Client etcd construit avec succès"
else
    echo "❌ Échec du build du client etcd"
    exit 1
fi

echo ""
echo "🎉 Build du projet terminé avec succès !"
echo ""
echo "💡 Utilisation:"
echo "   • Module constraint: cd constraint/cmd && ./constraint-parser ../tests/test_input.txt"
echo "   • Client etcd: ./etcd-client"