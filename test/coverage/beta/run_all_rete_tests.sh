#!/bin/bash
# Script automatisé pour exécuter tous les tests RETE beta
# Architecture refactorisée - Projet TSD

set -uo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
BETA_TESTS_DIR="$PROJECT_ROOT/beta_coverage_tests"
BINARY_PATH="$PROJECT_ROOT/bin/universal-rete-runner"

# Couleurs
GREEN='\033[32m'
RED='\033[31m'
YELLOW='\033[33m'
BLUE='\033[34m'
CYAN='\033[36m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔥 VALIDATION COMPLÈTE RETE - TOUS LES TESTS BETA${NC}"
echo -e "=================================================="

# Vérification du binaire
if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${YELLOW}⚠️  Binaire non trouvé, compilation...${NC}"
    cd "$PROJECT_ROOT"
    make build
fi

# Découverte des tests
TEST_FILES=($(find "$BETA_TESTS_DIR" -name "*.constraint" | sort))
TOTAL_TESTS=${#TEST_FILES[@]}

if [ $TOTAL_TESTS -eq 0 ]; then
    echo -e "${RED}❌ Aucun test trouvé dans $BETA_TESTS_DIR${NC}"
    exit 1
fi

echo -e "${CYAN}📊 Tests découverts: $TOTAL_TESTS${NC}"
echo

# Statistiques
PASSED=0
FAILED=0
TOTAL_TIME=0
FAILED_TESTS=()

# Fonction pour formater le temps
format_time() {
    local time_us=$1
    if [ $time_us -lt 1000 ]; then
        echo "${time_us}µs"
    elif [ $time_us -lt 1000000 ]; then
        echo "$((time_us / 1000))ms"
    else
        echo "$((time_us / 1000000))s"
    fi
}

# Exécution des tests
for constraint_file in "${TEST_FILES[@]}"; do
    test_name=$(basename "$constraint_file" .constraint)
    facts_file="${constraint_file%.constraint}.facts"

    if [ ! -f "$facts_file" ]; then
        echo -e "${RED}❌ $test_name - Fichier facts manquant${NC}"
        FAILED_TESTS+=("$test_name (fichier facts manquant)")
        ((FAILED++))
        continue
    fi

    echo -n -e "${CYAN}🎯 Test: ${test_name}${NC} ... "

    # Exécution du test
    start_time=$(date +%s%N)
    if output=$("$BINARY_PATH" "$constraint_file" "$facts_file" 2>&1); then
        end_time=$(date +%s%N)
        duration_ns=$((end_time - start_time))
        duration_us=$((duration_ns / 1000))
        TOTAL_TIME=$((TOTAL_TIME + duration_us))

        # Vérification du succès
        if echo "$output" | grep -q "✅ TEST VALIDÉ"; then
            # Extraction des métriques simplifiée
            tokens_observed=$(echo "$output" | grep "• Tokens observés.*:" | grep -o "[0-9]*" | head -1)
            tokens_expected=$(echo "$output" | grep "• Tokens attendus.*:" | grep -o "[0-9]*" | head -1)

            # Valeurs par défaut si extraction échoue
            tokens_observed=${tokens_observed:-"?"}
            tokens_expected=${tokens_expected:-"?"}

            echo -e "${GREEN}✅ RÉUSSI${NC} (${tokens_observed}/${tokens_expected} tokens, $(format_time $duration_us))"
            ((PASSED++))
        else
            echo -e "${RED}❌ ÉCHEC${NC} ($(format_time $duration_us))"
            FAILED_TESTS+=("$test_name (validation échoué)")
            ((FAILED++))
        fi
    else
        end_time=$(date +%s%N)
        duration_ns=$((end_time - start_time))
        duration_us=$((duration_ns / 1000))
        echo -e "${RED}❌ ERREUR${NC} ($(format_time $duration_us))"
        FAILED_TESTS+=("$test_name (erreur exécution)")
        ((FAILED++))
    fi
done

echo
echo -e "${BLUE}📊 RÉSULTATS FINAUX${NC}"
echo -e "=================="
echo -e "${GREEN}✅ Tests réussis: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}❌ Tests échoués: $FAILED${NC}"
    echo -e "${YELLOW}📋 Tests en échec:${NC}"
    for failed_test in "${FAILED_TESTS[@]}"; do
        echo -e "   ${RED}• $failed_test${NC}"
    done
else
    echo -e "${GREEN}❌ Tests échoués: $FAILED${NC}"
fi
echo -e "${CYAN}📊 Total: $TOTAL_TESTS${NC}"
echo -e "${CYAN}⏱️  Temps total: $(format_time $TOTAL_TIME)${NC}"

# Calcul du pourcentage
if [ $TOTAL_TESTS -gt 0 ]; then
    success_percentage=$((PASSED * 100 / TOTAL_TESTS))
    echo -e "${CYAN}📈 Taux de succès: ${success_percentage}%${NC}"
fi

echo

# Génération du rapport
REPORT_FILE="$PROJECT_ROOT/VALIDATION_RETE_$(date +%Y%m%d_%H%M%S).md"
cat > "$REPORT_FILE" << EOF
# RAPPORT VALIDATION RETE - $(date)

## Résumé Exécutif
- **Tests exécutés**: $TOTAL_TESTS
- **Tests réussis**: $PASSED
- **Tests échoués**: $FAILED
- **Taux de succès**: ${success_percentage:-0}%
- **Temps total**: $(format_time $TOTAL_TIME)

## Architecture
- **Méthode**: Tokens RÉELLEMENT extraits du réseau RETE
- **Binaire**: $BINARY_PATH
- **Tests**: $BETA_TESTS_DIR

## Détails des Tests
EOF

if [ $FAILED -gt 0 ]; then
    echo "### Tests en Échec" >> "$REPORT_FILE"
    for failed_test in "${FAILED_TESTS[@]}"; do
        echo "- $failed_test" >> "$REPORT_FILE"
    done
    echo "" >> "$REPORT_FILE"
fi

echo "### Conclusion" >> "$REPORT_FILE"
if [ $FAILED -eq 0 ]; then
    echo "✅ **VALIDATION COMPLÈTE RÉUSSIE** - Tous les tests RETE ont été validés avec succès." >> "$REPORT_FILE"
else
    echo "⚠️ **VALIDATION PARTIELLE** - $FAILED test(s) nécessitent une attention." >> "$REPORT_FILE"
fi

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 TOUS LES TESTS RETE ONT RÉUSSI !${NC}"
    echo -e "${GREEN}🚀 Validation complète du réseau RETE terminée avec succès${NC}"
    echo -e "${CYAN}📄 Rapport généré: $REPORT_FILE${NC}"
    exit 0
else
    echo -e "${RED}⚠️  $FAILED test(s) ont échoué${NC}"
    echo -e "${YELLOW}🔍 Vérifiez les contraintes et faits pour ces tests${NC}"
    echo -e "${CYAN}📄 Rapport généré: $REPORT_FILE${NC}"
    exit 1
fi
