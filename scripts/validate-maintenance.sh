#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

set -e

# Couleurs pour output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🔧 TSD - Validation Maintenance${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Compteurs
WARNINGS=0
ERRORS=0

# Fonction pour afficher succès
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Fonction pour afficher warning
warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS++))
}

# Fonction pour afficher erreur
error() {
    echo -e "${RED}❌ $1${NC}"
    ((ERRORS++))
}

# Fonction pour afficher info
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo "=== 1. Vérification Fichiers Temporaires ==="
TEMP_FILES=$(find . -type f \( -name "*.prof" -o -name "*.out" -o -name "*.test" \) ! -path "./.git/*" ! -path "*/vendor/*" 2>/dev/null || true)
if [ -z "$TEMP_FILES" ]; then
    success "Aucun fichier temporaire trouvé"
else
    warning "Fichiers temporaires détectés:"
    echo "$TEMP_FILES"
    info "Exécutez: make clean ou supprimez manuellement"
fi
echo ""

echo "=== 2. Vérification Dépendances ==="
if go mod verify &>/dev/null; then
    success "go mod verify - OK"
else
    error "go mod verify - ÉCHEC"
fi

if go mod tidy -diff &>/dev/null; then
    success "go mod tidy - À jour"
else
    warning "go mod tidy - Modifications nécessaires"
    info "Exécutez: go mod tidy"
fi
echo ""

echo "=== 3. Formatage Code ==="
UNFORMATTED=$(gofmt -l . 2>/dev/null | grep -v vendor || true)
if [ -z "$UNFORMATTED" ]; then
    success "Tous les fichiers sont formatés"
else
    warning "Fichiers non formatés détectés"
    echo "$UNFORMATTED"
    info "Exécutez: go fmt ./..."
fi
echo ""

echo "=== 4. Imports ==="
if command -v goimports &> /dev/null; then
    UNORGANIZED=$(goimports -l . 2>/dev/null | grep -v vendor | head -10 || true)
    if [ -z "$UNORGANIZED" ]; then
        success "Imports organisés"
    else
        warning "Imports à réorganiser (premiers 10):"
        echo "$UNORGANIZED"
        info "Exécutez: goimports -w ."
    fi
else
    info "goimports non installé - skip"
fi
echo ""

echo "=== 5. Analyse Statique ==="
if command -v staticcheck &> /dev/null; then
    STATIC_ISSUES=$(staticcheck ./... 2>&1 | wc -l)
    if [ "$STATIC_ISSUES" -eq 0 ]; then
        success "Aucun problème staticcheck"
    else
        warning "staticcheck a détecté $STATIC_ISSUES problème(s)"
        info "Exécutez: staticcheck ./..."
    fi
else
    info "staticcheck non installé - skip"
fi
echo ""

echo "=== 6. Tests ==="
info "Exécution des tests..."
if go test ./... -short &>/dev/null; then
    success "Tests passent (mode short)"
else
    error "Certains tests échouent"
    info "Exécutez: go test ./... pour détails"
fi
echo ""

echo "=== 7. Couverture Globale ==="
COVERAGE=$(go test -cover ./... 2>&1 | grep "coverage:" | awk '{sum+=$5; count++} END {if(count>0) print sum/count; else print 0}' | cut -d'%' -f1)
if [ ! -z "$COVERAGE" ]; then
    COVERAGE_INT=$(echo "$COVERAGE" | cut -d'.' -f1)
    if [ "$COVERAGE_INT" -ge 80 ]; then
        success "Couverture moyenne: ${COVERAGE}% (objectif: 80%+)"
    elif [ "$COVERAGE_INT" -ge 70 ]; then
        warning "Couverture moyenne: ${COVERAGE}% (objectif: 80%+)"
    else
        error "Couverture moyenne: ${COVERAGE}% (objectif: 80%+)"
    fi
else
    info "Impossible de calculer la couverture"
fi
echo ""

echo "=== 8. Complexité Cyclomatique ==="
if command -v gocyclo &> /dev/null; then
    HIGH_COMPLEXITY=$(gocyclo -over 20 . 2>/dev/null | grep -v "_test.go" | wc -l)
    if [ "$HIGH_COMPLEXITY" -eq 0 ]; then
        success "Complexité acceptable (seuil: 20)"
    else
        warning "$HIGH_COMPLEXITY fonction(s) de production avec complexité > 20"
        info "Exécutez: gocyclo -over 20 . pour détails"
    fi
else
    info "gocyclo non installé - skip"
fi
echo ""

echo "=== 9. TODOs et FIXMEs ==="
TODO_COUNT=$(grep -rn "TODO\|FIXME\|XXX\|HACK" --include="*.go" . 2>/dev/null | grep -v vendor | wc -l)
if [ "$TODO_COUNT" -eq 0 ]; then
    success "Aucun TODO/FIXME trouvé"
elif [ "$TODO_COUNT" -lt 20 ]; then
    info "$TODO_COUNT TODO/FIXME trouvé(s) - acceptable"
else
    warning "$TODO_COUNT TODO/FIXME trouvé(s)"
    info "Consultez: REPORTS/MAINTENANCE_TODO.md"
fi
echo ""

echo "=== 10. Code Non Utilisé (deadcode) ==="
if command -v deadcode &> /dev/null; then
    DEAD_COUNT=$(deadcode ./... 2>/dev/null | wc -l)
    if [ "$DEAD_COUNT" -eq 0 ]; then
        success "Aucun code mort détecté"
    else
        warning "$DEAD_COUNT élément(s) de code potentiellement non utilisé(s)"
        info "Exécutez: deadcode ./..."
    fi
else
    info "deadcode non installé - skip"
    info "Installez avec: go install golang.org/x/tools/cmd/deadcode@latest"
fi
echo ""

echo "=== 11. Vulnérabilités ==="
if command -v govulncheck &> /dev/null; then
    if govulncheck ./... &>/dev/null; then
        success "Aucune vulnérabilité connue"
    else
        warning "Vulnérabilités détectées"
        info "Exécutez: govulncheck ./..."
    fi
else
    info "govulncheck non installé - skip"
    info "Installez avec: go install golang.org/x/vuln/cmd/govulncheck@latest"
fi
echo ""

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📊 Résumé${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ Parfait! Aucun problème détecté${NC}"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  $WARNINGS warning(s) - voir ci-dessus${NC}"
    exit 0
else
    echo -e "${RED}❌ $ERRORS erreur(s), $WARNINGS warning(s)${NC}"
    echo ""
    echo "Recommandations:"
    echo "1. Corrigez les erreurs en priorité"
    echo "2. Traitez les warnings si possible"
    echo "3. Exécutez: make validate pour validation complète"
    exit 1
fi
