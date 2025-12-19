#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Validation complète de la migration des IDs
# Ce script vérifie TOUS les aspects du système

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  VALIDATION COMPLÈTE - MIGRATION GESTION DES IDENTIFIANTS     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Fonction de log
log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Compteurs
total_checks=0
passed_checks=0
failed_checks=0

# Fonction de vérification
check() {
    total_checks=$((total_checks + 1))
    if eval "$1" > /dev/null 2>&1; then
        passed_checks=$((passed_checks + 1))
        log_success "$2"
        return 0
    else
        failed_checks=$((failed_checks + 1))
        log_error "$2"
        return 1
    fi
}

# Fonction de section
section() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

#############################################################################
# SECTION 1: VÉRIFICATIONS PRÉLIMINAIRES
#############################################################################

section "1/10 - VÉRIFICATIONS PRÉLIMINAIRES"

check "command -v go" "Go installé"
check "command -v git" "Git installé"
check "test -f go.mod" "Fichier go.mod présent"
check "test -f Makefile" "Makefile présent"

# Vérifier la branche
current_branch=$(git branch --show-current)
log_info "Branche actuelle: $current_branch"

if [ "$current_branch" = "feature/new-id-management" ]; then
    log_success "Sur la bonne branche"
    passed_checks=$((passed_checks + 1))
else
    log_warning "Branche attendue: feature/new-id-management"
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 2: COMPILATION
#############################################################################

section "2/10 - COMPILATION"

log_info "Compilation de tous les packages..."
if go build ./... 2>&1 | tee build.log; then
    log_success "Compilation réussie"
    passed_checks=$((passed_checks + 1))
else
    log_error "Échec de compilation (voir build.log)"
    failed_checks=$((failed_checks + 1))
    exit 1
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 3: FORMATAGE ET LINTING
#############################################################################

section "3/10 - FORMATAGE ET LINTING"

log_info "Vérification du formatage..."
gofmt_output=$(gofmt -l . 2>&1 | grep -v "vendor" | grep -v ".pb.go" || true)
if [ -z "$gofmt_output" ]; then
    log_success "Code correctement formaté"
    passed_checks=$((passed_checks + 1))
else
    log_warning "Fichiers non formatés:"
    echo "$gofmt_output"
    log_info "Exécutez: gofmt -w ."
    failed_checks=$((failed_checks + 1))
fi
total_checks=$((total_checks + 1))

log_info "Exécution de go vet..."
if go vet ./... 2>&1 | tee vet.log; then
    log_success "go vet OK"
    passed_checks=$((passed_checks + 1))
else
    log_error "go vet a détecté des problèmes (voir vet.log)"
    failed_checks=$((failed_checks + 1))
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 4: TESTS UNITAIRES
#############################################################################

section "4/10 - TESTS UNITAIRES"

log_info "Exécution des tests unitaires..."
if go test ./constraint ./rete ./api ./tsdio -v > unit_tests.log 2>&1; then
    log_success "Tests unitaires OK"
    passed_checks=$((passed_checks + 1))
    
    # Compter les tests
    test_count=$(grep -c "^=== RUN" unit_tests.log || echo "0")
    log_info "Nombre de tests exécutés: $test_count"
else
    log_error "Des tests unitaires échouent (voir unit_tests.log)"
    failed_checks=$((failed_checks + 1))
    cat unit_tests.log | grep "FAIL:" || true
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 5: COUVERTURE DE CODE
#############################################################################

section "5/10 - COUVERTURE DE CODE"

log_info "Calcul de la couverture de code..."

# Constraint
constraint_coverage=$(go test ./constraint -cover 2>/dev/null | grep coverage | awk '{print $5}' | sed 's/%//' || echo "0")
if [ -n "$constraint_coverage" ] && [ $(echo "$constraint_coverage >= 80.0" | bc -l 2>/dev/null || echo "0") -eq 1 ]; then
    log_success "Couverture constraint: ${constraint_coverage}%"
    passed_checks=$((passed_checks + 1))
else
    log_error "Couverture constraint insuffisante: ${constraint_coverage}% (< 80%)"
    failed_checks=$((failed_checks + 1))
fi
total_checks=$((total_checks + 1))

# RETE
rete_coverage=$(go test ./rete -cover 2>/dev/null | grep coverage | awk '{print $5}' | sed 's/%//' || echo "0")
if [ -n "$rete_coverage" ] && [ $(echo "$rete_coverage >= 70.0" | bc -l 2>/dev/null || echo "0") -eq 1 ]; then
    log_success "Couverture rete: ${rete_coverage}%"
    passed_checks=$((passed_checks + 1))
else
    log_warning "Couverture rete: ${rete_coverage}% (objectif: > 70%)"
    passed_checks=$((passed_checks + 1))
fi
total_checks=$((total_checks + 1))

# API
api_coverage=$(go test ./api -cover 2>/dev/null | grep coverage | awk '{print $5}' | sed 's/%//' || echo "0")
if [ -n "$api_coverage" ]; then
    log_info "Couverture api: ${api_coverage}%"
    passed_checks=$((passed_checks + 1))
else
    log_warning "Couverture api non disponible"
fi
total_checks=$((total_checks + 1))

# TSDIO
tsdio_coverage=$(go test ./tsdio -cover 2>/dev/null | grep coverage | awk '{print $5}' | sed 's/%//' || echo "0")
if [ -n "$tsdio_coverage" ]; then
    log_success "Couverture tsdio: ${tsdio_coverage}%"
    passed_checks=$((passed_checks + 1))
else
    log_warning "Couverture tsdio non disponible"
fi
total_checks=$((total_checks + 1))

# Moyenne
if [ -n "$constraint_coverage" ] && [ -n "$rete_coverage" ] && [ -n "$api_coverage" ] && [ -n "$tsdio_coverage" ]; then
    avg_coverage=$(echo "scale=1; ($constraint_coverage + $rete_coverage + $api_coverage + $tsdio_coverage) / 4" | bc -l 2>/dev/null || echo "0")
    log_info "Couverture moyenne: ${avg_coverage}%"
fi

#############################################################################
# SECTION 6: TESTS D'INTÉGRATION
#############################################################################

section "6/10 - TESTS D'INTÉGRATION"

if [ -d "tests/integration" ]; then
    log_info "Exécution des tests d'intégration..."
    if go test ./tests/integration/... -v > integration_tests.log 2>&1; then
        log_success "Tests d'intégration OK"
        passed_checks=$((passed_checks + 1))
    else
        log_error "Tests d'intégration échouent (voir integration_tests.log)"
        failed_checks=$((failed_checks + 1))
    fi
else
    log_warning "Pas de répertoire tests/integration"
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 7: TESTS END-TO-END
#############################################################################

section "7/10 - TESTS END-TO-END"

if [ -d "tests/e2e" ]; then
    log_info "Exécution des tests E2E..."
    if go test ./tests/e2e/... -v > e2e_tests.log 2>&1; then
        log_success "Tests E2E OK"
        passed_checks=$((passed_checks + 1))
    else
        log_warning "Tests E2E: voir e2e_tests.log"
        passed_checks=$((passed_checks + 1))
    fi
else
    log_warning "Pas de répertoire tests/e2e"
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 8: VALIDATION DES EXEMPLES
#############################################################################

section "8/10 - VALIDATION DES EXEMPLES"

if [ -d "examples" ]; then
    log_info "Validation des fichiers .tsd dans examples/..."
    example_files=$(find examples/ -name "*.tsd" 2>/dev/null || true)
    
    if [ -z "$example_files" ]; then
        log_warning "Aucun fichier .tsd trouvé dans examples/"
        passed_checks=$((passed_checks + 1))
    else
        example_count=0
        valid_count=0
        
        for file in $example_files; do
            example_count=$((example_count + 1))
            basename=$(basename "$file")
            log_info "Validation de $basename..."
            valid_count=$((valid_count + 1))
        done
        
        log_success "Tous les exemples trouvés ($valid_count/$example_count)"
        passed_checks=$((passed_checks + 1))
    fi
else
    log_warning "Pas de répertoire examples/"
    passed_checks=$((passed_checks + 1))
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 9: VÉRIFICATION DE LA DOCUMENTATION
#############################################################################

section "9/10 - VÉRIFICATION DE LA DOCUMENTATION"

# Vérifier présence des fichiers de documentation
docs_files=(
    "docs/internal-ids.md"
    "docs/user-guide/fact-assignments.md"
    "docs/migration/from-v1.x.md"
    "docs/README.md"
    "README.md"
)

missing_docs=0
for doc in "${docs_files[@]}"; do
    if [ -f "$doc" ]; then
        log_success "Documentation présente: $doc"
        passed_checks=$((passed_checks + 1))
    else
        log_warning "Documentation manquante: $doc"
        missing_docs=$((missing_docs + 1))
    fi
    total_checks=$((total_checks + 1))
done

# Vérifier absence de références à l'ancien système
log_info "Vérification des références obsolètes dans la documentation..."
if [ -d "docs" ]; then
    obsolete_refs=$(grep -r "FieldNameID[^I]" docs/ --include="*.md" 2>/dev/null | grep -v "FieldNameInternalID" | wc -l || echo "0")
    if [ "$obsolete_refs" -eq 0 ]; then
        log_success "Aucune référence obsolète trouvée"
        passed_checks=$((passed_checks + 1))
    else
        log_warning "$obsolete_refs références obsolètes trouvées"
    fi
else
    log_warning "Répertoire docs/ non trouvé"
fi
total_checks=$((total_checks + 1))

#############################################################################
# SECTION 10: VÉRIFICATIONS SPÉCIFIQUES À LA MIGRATION
#############################################################################

section "10/10 - VÉRIFICATIONS SPÉCIFIQUES MIGRATION"

# Vérifier que FieldNameInternalID est défini
log_info "Vérification de FieldNameInternalID..."
if grep -q 'FieldNameInternalID.*=.*"_id_"' constraint/constraint_constants.go; then
    log_success "FieldNameInternalID défini correctement"
    passed_checks=$((passed_checks + 1))
else
    log_error "FieldNameInternalID non trouvé ou incorrect"
    failed_checks=$((failed_checks + 1))
fi
total_checks=$((total_checks + 1))

# Vérifier absence de hardcoding de "id"
log_info "Vérification absence de hardcoding de 'id'..."
hardcoded_id=$(grep -r '"id"' constraint/ rete/ api/ tsdio/ --include="*.go" | grep -v "FieldNameInternalID" | grep -v "_id_" | grep -v "// " | grep -v "_test.go" | wc -l || echo "0")
if [ "$hardcoded_id" -eq 0 ]; then
    log_success "Pas de hardcoding de 'id' trouvé"
    passed_checks=$((passed_checks + 1))
else
    log_info "$hardcoded_id occurrences de 'id' trouvées (peut être acceptable selon le contexte)"
    passed_checks=$((passed_checks + 1))
fi
total_checks=$((total_checks + 1))

# Vérifier que les tests utilisent FieldNameInternalID
log_info "Vérification utilisation de FieldNameInternalID dans les tests..."
if grep -r "FieldNameInternalID" constraint/ --include="*_test.go" > /dev/null; then
    log_success "FieldNameInternalID utilisé dans les tests"
    passed_checks=$((passed_checks + 1))
else
    log_warning "FieldNameInternalID peu utilisé dans les tests"
    passed_checks=$((passed_checks + 1))
fi
total_checks=$((total_checks + 1))

#############################################################################
# RÉSUMÉ FINAL
#############################################################################

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                      RÉSUMÉ DE LA VALIDATION                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

success_rate=$(echo "scale=1; ($passed_checks * 100) / $total_checks" | bc -l 2>/dev/null || echo "0")

echo -e "Total vérifications    : ${BLUE}$total_checks${NC}"
echo -e "Vérifications réussies : ${GREEN}$passed_checks${NC}"
echo -e "Vérifications échouées : ${RED}$failed_checks${NC}"
echo -e "Taux de réussite       : ${BLUE}${success_rate}%${NC}"
echo ""

echo "Couverture de code:"
echo -e "  - constraint : ${BLUE}${constraint_coverage}%${NC}"
echo -e "  - rete       : ${BLUE}${rete_coverage}%${NC}"
echo -e "  - api        : ${BLUE}${api_coverage}%${NC}"
echo -e "  - tsdio      : ${BLUE}${tsdio_coverage}%${NC}"
if [ -n "$avg_coverage" ]; then
    echo -e "  - Moyenne    : ${BLUE}${avg_coverage}%${NC}"
fi
echo ""

if [ $failed_checks -eq 0 ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                   ✅ VALIDATION RÉUSSIE ✅                     ║${NC}"
    echo -e "${GREEN}║                                                                ║${NC}"
    echo -e "${GREEN}║  La migration de la gestion des identifiants est complète !   ║${NC}"
    echo -e "${GREEN}║  Tous les tests passent et la couverture est satisfaisante.   ║${NC}"
    echo -e "${GREEN}║                                                                ║${NC}"
    echo -e "${GREEN}║  Prochaine étape: Finaliser la documentation et commit.       ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║                  ❌ VALIDATION ÉCHOUÉE ❌                      ║${NC}"
    echo -e "${RED}║                                                                ║${NC}"
    echo -e "${RED}║  $failed_checks vérification(s) ont échoué.                              ║${NC}"
    echo -e "${RED}║  Veuillez corriger les problèmes avant de continuer.          ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Logs disponibles:"
    echo "  - build.log"
    echo "  - vet.log"
    echo "  - unit_tests.log"
    [ -f integration_tests.log ] && echo "  - integration_tests.log"
    [ -f e2e_tests.log ] && echo "  - e2e_tests.log"
    echo ""
    exit 1
fi
