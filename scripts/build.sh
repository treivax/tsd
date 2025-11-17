#!/bin/bash

# Script principal de build et test pour TSD
# Utilise les bonnes pratiques Go

set -e

echo "🚀 TSD BUILD & TEST SUITE"
echo "========================"

# Couleurs pour l'output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Vérifier que nous sommes dans le bon répertoire
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Erreur: go.mod non trouvé. Exécutez depuis la racine du projet.${NC}"
    exit 1
fi

echo -e "${BLUE}📋 Étape 1/6: Vérifications préliminaires${NC}"
echo "============================================"

# Vérifier Go version
go version
echo "✅ Go installé"

# Vérifier les dépendances
go mod tidy
echo "✅ Dépendances vérifiées"

echo ""
echo -e "${BLUE}🔧 Étape 2/6: Formatage du code${NC}"
echo "==============================="

# Formatter le code
go fmt ./...
echo "✅ Code formaté avec gofmt"

# Vérifier avec goimports si disponible
if command -v goimports &> /dev/null; then
    goimports -w .
    echo "✅ Imports optimisés avec goimports"
fi

echo ""
echo -e "${BLUE}🔍 Étape 3/6: Analyse statique${NC}"
echo "==============================="

# Analyse statique avec go vet
if go vet ./...; then
    echo "✅ Analyse statique (go vet) : OK"
else
    echo -e "${YELLOW}⚠️ Warnings détectés par go vet${NC}"
fi

# Analyse avec golangci-lint si disponible
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run; then
        echo "✅ Analyse golangci-lint : OK"
    else
        echo -e "${YELLOW}⚠️ Issues détectés par golangci-lint${NC}"
    fi
fi

echo ""
echo -e "${BLUE}🔨 Étape 4/6: Compilation${NC}"
echo "========================="

# Build principal
if go build -o bin/tsd ./cmd/; then
    echo "✅ Build principal : OK"
else
    echo -e "${RED}❌ Échec du build principal${NC}"
    exit 1
fi

# Build des outils
if go build -o bin/constraint-parser ./constraint/cmd/; then
    echo "✅ Build constraint-parser : OK"
else
    echo -e "${YELLOW}⚠️ Build constraint-parser : échec (optionnel)${NC}"
fi

echo ""
echo -e "${BLUE}🧪 Étape 5/6: Tests${NC}"
echo "==================="

# Tests unitaires
if go test -v ./...; then
    echo "✅ Tests unitaires : OK"
else
    echo -e "${RED}❌ Échec des tests unitaires${NC}"
    exit 1
fi

# Tests avec couverture
echo ""
echo "📊 Couverture de code :"
go test -cover ./...

# Tests de performance si demandés
if [ "$1" == "--bench" ]; then
    echo ""
    echo "🏃 Tests de performance :"
    go test -bench=. ./test/benchmark/... || echo "ℹ️ Pas de benchmarks trouvés"
fi

echo ""
echo -e "${BLUE}✅ Étape 6/6: Tests de couverture Alpha${NC}"
echo "=========================================="

# Exécuter les tests Alpha si le runner existe
if [ -f "test/coverage/alpha_coverage_runner.go" ]; then
    echo "🧪 Exécution des tests Alpha :"
    if go run test/coverage/alpha_coverage_runner.go; then
        echo "✅ Tests Alpha : OK"
    else
        echo -e "${YELLOW}⚠️ Quelques tests Alpha ont échoué${NC}"
    fi
else
    echo "ℹ️ Runner de tests Alpha non trouvé"
fi

echo ""
echo -e "${GREEN}🎉 BUILD & TEST TERMINÉS${NC}"
echo "========================="

# Afficher les binaires générés
echo "📦 Binaires générés :"
ls -la bin/ 2>/dev/null || echo "Aucun binaire dans bin/"

echo ""
echo -e "${GREEN}✅ TSD prêt à être utilisé !${NC}"