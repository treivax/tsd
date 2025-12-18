#!/bin/bash
# Script de validation post-déploiement du fix bug 'once'
# Version: 1.1.0
# Date: 2025-12-18

set -e  # Exit on error

# Couleurs pour output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonctions d'affichage
print_header() {
    echo -e "\n${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}\n"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Variables
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TESTS_PASSED=0
TESTS_FAILED=0

print_header "🚀 VALIDATION POST-DÉPLOIEMENT - FIX BUG 'ONCE'"

print_info "Projet: $PROJECT_ROOT"
print_info "Go version: $(go version)"
echo ""

# ═══════════════════════════════════════════════════════════════
# ÉTAPE 1: Tests unitaires xuples
# ═══════════════════════════════════════════════════════════════
print_header "ÉTAPE 1: Tests Unitaires Xuples"

print_info "Exécution de la suite complète..."
if go test -v ./xuples -count=1 > /tmp/xuples_test.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Suite xuples: PASS"

    # Compter les tests
    TOTAL_TESTS=$(grep -c "^--- PASS:" /tmp/xuples_test.log || echo "0")
    print_info "  Tests passés: $TOTAL_TESTS"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Suite xuples: FAIL"
    echo "Voir détails dans: /tmp/xuples_test.log"
fi

# ═══════════════════════════════════════════════════════════════
# ÉTAPE 2: Tests spécifiques au bug fix
# ═══════════════════════════════════════════════════════════════
print_header "ÉTAPE 2: Tests Spécifiques Bug Fix"

# Test 1: RetrieveAutomaticallyMarksConsumed
print_info "Test: RetrieveAutomaticallyMarksConsumed..."
if go test ./xuples -run TestRetrieveAutomaticallyMarksConsumed -v > /tmp/test_once.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Bug 'once' corrigé validé"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Bug 'once' NON corrigé!"
    cat /tmp/test_once.log
fi

# Test 2: RetrievePerAgentPolicy
print_info "Test: RetrievePerAgentPolicy..."
if go test ./xuples -run TestRetrievePerAgentPolicy -v > /tmp/test_per_agent.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Politique per-agent validée"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Politique per-agent échouée"
fi

# Test 3: RetrieveLimitedPolicy
print_info "Test: RetrieveLimitedPolicy..."
if go test ./xuples -run TestRetrieveLimitedPolicy -v > /tmp/test_limited.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Politique limited(n) validée"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Politique limited(n) échouée"
fi

# Test 4: MultipleXuplesWithOncePolicy
print_info "Test: MultipleXuplesWithOncePolicy..."
if go test ./xuples -run TestMultipleXuplesWithOncePolicy -v > /tmp/test_multiple.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Consommation multiple avec 'once' validée"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Consommation multiple échouée"
fi

# ═══════════════════════════════════════════════════════════════
# ÉTAPE 3: Tests E2E
# ═══════════════════════════════════════════════════════════════
print_header "ÉTAPE 3: Tests End-to-End"

print_info "Exécution tests E2E xuples..."
if go test ./tests/e2e -run XuplesE2E -v > /tmp/e2e_test.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Tests E2E: PASS"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Tests E2E: FAIL"
    echo "Voir détails dans: /tmp/e2e_test.log"
fi

# ═══════════════════════════════════════════════════════════════
# ÉTAPE 4: Tests de race conditions
# ═══════════════════════════════════════════════════════════════
print_header "ÉTAPE 4: Tests Race Conditions"

print_info "Exécution avec -race detector..."
if go test ./xuples -race -run "TestConcurrent|TestRace" > /tmp/race_test.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Race detector: PASS (aucune race détectée)"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Race conditions détectées!"
    cat /tmp/race_test.log
fi

# ═══════════════════════════════════════════════════════════════
# ÉTAPE 5: Vérification compilation
# ═══════════════════════════════════════════════════════════════
print_header "ÉTAPE 5: Vérification Compilation"

print_info "Build du package xuples..."
if go build ./xuples > /tmp/build.log 2>&1; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    print_success "Compilation: PASS"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    print_error "Compilation: FAIL"
    cat /tmp/build.log
fi

# ═══════════════════════════════════════════════════════════════
# ÉTAPE 6: Vérification documentation
# ═══════════════════════════════════════════════════════════════
print_header "ÉTAPE 6: Vérification Documentation"

# Vérifier que les fichiers de doc existent
DOCS=(
    "RAPPORT_DEPLOIEMENT_BUG_FIX.md"
    "CHANGELOG_v1.1.0.md"
)

DOC_OK=true
for doc in "${DOCS[@]}"; do
    if [ -f "$doc" ]; then
        print_success "Documentation trouvée: $doc"
    else
        print_warning "Documentation manquante: $doc"
        DOC_OK=false
    fi
done

if [ "$DOC_OK" = true ]; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# ═══════════════════════════════════════════════════════════════
# ÉTAPE 7: Vérification code coverage (optionnel)
# ═══════════════════════════════════════════════════════════════
print_header "ÉTAPE 7: Code Coverage (Optionnel)"

print_info "Calcul du code coverage..."
if go test ./xuples -coverprofile=/tmp/coverage.out > /dev/null 2>&1; then
    COVERAGE=$(go tool cover -func=/tmp/coverage.out | grep total | awk '{print $3}')
    print_info "Coverage total: $COVERAGE"

    # Extraire le pourcentage
    COVERAGE_PCT=$(echo "$COVERAGE" | sed 's/%//')
    if (( $(echo "$COVERAGE_PCT >= 80.0" | bc -l) )); then
        print_success "Coverage satisfaisant (>= 80%)"
    else
        print_warning "Coverage sous 80%: $COVERAGE"
    fi
else
    print_warning "Impossible de calculer le coverage"
fi

# ═══════════════════════════════════════════════════════════════
# RAPPORT FINAL
# ═══════════════════════════════════════════════════════════════
print_header "📊 RAPPORT FINAL"

echo -e "${BLUE}Tests passés:${NC} ${GREEN}$TESTS_PASSED${NC}"
echo -e "${BLUE}Tests échoués:${NC} ${RED}$TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    print_success "TOUS LES TESTS PASSENT! ✅"
    echo ""
    print_success "Le fix du bug 'once' est validé et prêt pour production."
    echo ""
    print_info "Prochaines étapes recommandées:"
    echo "  1. Revue du code par un pair"
    echo "  2. Merge dans la branche principale"
    echo "  3. Tag version v1.1.0"
    echo "  4. Communication aux utilisateurs"
    echo ""
    exit 0
else
    print_error "ÉCHECS DÉTECTÉS! ❌"
    echo ""
    print_error "Le déploiement n'est PAS validé."
    echo ""
    print_info "Logs disponibles dans /tmp/*.log"
    echo "  - /tmp/xuples_test.log"
    echo "  - /tmp/e2e_test.log"
    echo "  - /tmp/race_test.log"
    echo "  - /tmp/build.log"
    echo ""
    exit 1
fi
