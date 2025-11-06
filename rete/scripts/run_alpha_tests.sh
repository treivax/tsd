#!/bin/bash

# Script de tests automatisés pour les conditions Alpha du réseau RETE
# Ce script exécute une suite complète de tests pour valider tous les types d'expressions

set -e  # Arrêter en cas d'erreur

echo "🧪 TESTS AUTOMATISÉS DES CONDITIONS ALPHA"
echo "========================================"

# Couleurs pour l'affichage
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# Nettoyer les fichiers de couverture précédents
rm -f *.out coverage.html

echo -e "${BLUE}📋 Phase 1: Tests de couverture complète${NC}"
echo "Tests de toutes les expressions et opérateurs..."

go test -v -cover -coverprofile=comprehensive_coverage.out \
    -run="TestAlphaConditionEvaluator_ComprehensiveCoverage" .
print_result $? "Tests de couverture complète"

echo -e "\n${BLUE}📋 Phase 2: Tests des cas d'erreur${NC}"
echo "Validation de la gestion d'erreur robuste..."

go test -v -run="TestAlphaConditionEvaluator_ExtendedErrorCases" .
print_result $? "Tests des cas d'erreur"

echo -e "\n${BLUE}📋 Phase 3: Tests des cas limites${NC}"
echo "Validation des valeurs extrêmes et cas spéciaux..."

go test -v -run="TestAlphaConditionEvaluator_EdgeCases" .
print_result $? "Tests des cas limites"

echo -e "\n${BLUE}📋 Phase 4: Tests du constructeur Alpha${NC}"
echo "Validation de toutes les méthodes du builder..."

go test -v -run="TestAlphaConditionBuilder_AllMethods" .
print_result $? "Tests du constructeur Alpha"

echo -e "\n${BLUE}📋 Phase 5: Tests d'intégration RETE${NC}"
echo "Validation de l'intégration avec le réseau RETE..."

go test -v -run="TestAlphaConditionEvaluator_Integration" .
print_result $? "Tests d'intégration RETE"

echo -e "\n${BLUE}📋 Phase 6: Tests des liaisons de variables${NC}"
echo "Validation de la gestion des variables..."

go test -v -run="TestAlphaConditionEvaluator_VariableBindings" .
print_result $? "Tests des liaisons de variables"

echo -e "\n${BLUE}📊 Phase 7: Benchmark de performance${NC}"
echo "Mesure des performances des conditions Alpha..."

go test -bench="BenchmarkAlphaConditionEvaluator" -benchmem -run="^$" .
print_result $? "Benchmark de performance"

echo -e "\n${BLUE}📈 Phase 8: Analyse de couverture détaillée${NC}"
echo "Génération du rapport de couverture..."

# Combiner toutes les métriques de couverture
go test -cover -coverprofile=full_alpha_coverage.out \
    -run="TestAlphaConditionEvaluator_ComprehensiveCoverage|TestAlphaConditionEvaluator_ExtendedErrorCases|TestAlphaConditionEvaluator_EdgeCases|TestAlphaConditionBuilder_AllMethods|TestAlphaConditionEvaluator_Integration|TestAlphaConditionEvaluator_VariableBindings" .

# Générer le rapport HTML
go tool cover -html=full_alpha_coverage.out -o alpha_coverage.html

# Afficher les statistiques de couverture
echo -e "\n${YELLOW}📊 STATISTIQUES DE COUVERTURE:${NC}"
go tool cover -func=full_alpha_coverage.out | tail -1

echo -e "\n${GREEN}🎉 TOUS LES TESTS AUTOMATISÉS RÉUSSIS !${NC}"
echo -e "${GREEN}✨ Les conditions Alpha du réseau RETE sont entièrement validées${NC}"

echo -e "\n${BLUE}📂 Fichiers générés:${NC}"
echo "  - full_alpha_coverage.out (données de couverture)"
echo "  - alpha_coverage.html (rapport HTML détaillé)"

echo -e "\n${BLUE}🔍 Types d'expressions testés:${NC}"
echo "  ✅ Littéraux booléens (true/false)"
echo "  ✅ Opérateurs binaires (==, !=, <, <=, >, >=)"
echo "  ✅ Expressions logiques (AND, OR)"
echo "  ✅ Comparaisons numériques (int, float)"
echo "  ✅ Comparaisons de chaînes"
echo "  ✅ Comparaisons booléennes"
echo "  ✅ Valeurs négatives et zéro"
echo "  ✅ Valeurs limites (MaxInt64, MaxFloat64, Infinity)"
echo "  ✅ Expressions imbriquées complexes"
echo "  ✅ Conversion automatique de types"
echo "  ✅ Gestion d'erreurs robuste"
echo "  ✅ Intégration avec nœuds Alpha"
echo "  ✅ Liaisons de variables"

echo -e "\n${BLUE}⚡ Performances mesurées:${NC}"
echo "  - Temps d'évaluation par condition"
echo "  - Allocation mémoire"
echo "  - Nombre d'allocations"

echo -e "\n${GREEN}🚀 Le système RETE est prêt pour la production !${NC}"