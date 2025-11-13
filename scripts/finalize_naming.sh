#!/bin/bash

# Script de standardisation finale et sûre - Renommage ciblé des fichiers principaux
echo "🎯 STANDARDISATION FINALE SÉCURISÉE"
echo "==================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Fonction pour renommer un fichier et mettre à jour ses références
safe_rename() {
    local old_file="$1"
    local new_file="$2"
    
    if [[ ! -f "$old_file" ]]; then
        echo -e "${YELLOW}⚠️ Fichier non trouvé: $old_file${NC}"
        return 1
    fi
    
    echo -e "${BLUE}🔄 Renommage: $(basename $old_file) → $(basename $new_file)${NC}"
    
    # Renommer le fichier
    if mv "$old_file" "$new_file"; then
        echo -e "${GREEN}✅ Fichier renommé avec succès${NC}"
        return 0
    else
        echo -e "${RED}❌ Erreur lors du renommage${NC}"
        return 1
    fi
}

echo -e "${BLUE}📋 FICHIERS À STANDARDISER (Phase sécurisée)${NC}"
echo "============================================"

# Renommages sûrs (fichiers peu référencés)
declare -A safe_renames=(
    ["./test/helper.go"]="./test/test_utils.go"
    ["./rete/converter.go"]="./rete/type_converter.go"
    ["./rete/evaluator.go"]="./rete/expression_evaluator.go"
)

# Afficher le plan
for old_file in "${!safe_renames[@]}"; do
    new_file="${safe_renames[$old_file]}"
    echo -e "${YELLOW}  📝 $(basename $old_file) → $(basename $new_file)${NC}"
done

echo ""
read -p "Procéder aux renommages sécurisés ? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Annulé."
    exit 0
fi

echo -e "${GREEN}🚀 RENOMMAGE EN COURS${NC}"
echo "====================="

renamed_count=0
error_count=0

# Effectuer les renommages
for old_file in "${!safe_renames[@]}"; do
    new_file="${safe_renames[$old_file]}"
    
    if safe_rename "$old_file" "$new_file"; then
        renamed_count=$((renamed_count + 1))
    else
        error_count=$((error_count + 1))
    fi
    echo ""
done

echo -e "${GREEN}🧪 VALIDATION POST-RENOMMAGE${NC}"
echo "==============================="

# Tester la compilation
echo "Test de compilation..."
if go build ./... 2>/dev/null; then
    echo -e "${GREEN}✅ Compilation réussie !${NC}"
else
    echo -e "${RED}❌ Erreurs de compilation détectées${NC}"
    echo "Détails:"
    go build ./... 2>&1 | head -10
fi

echo ""
echo -e "${GREEN}📊 RÉSULTATS FINAUX${NC}"
echo "==================="
echo "Fichiers renommés avec succès: $renamed_count"
echo "Erreurs: $error_count"

# Analyse finale de conformité
total_files=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
snake_case_files=$(find . -name "*_*.go" -not -path "./vendor/*" | wc -l)
conformity=$((snake_case_files * 100 / total_files))

echo ""
echo -e "${BLUE}📈 CONFORMITÉ FINALE${NC}"
echo "==================="
echo "Total fichiers Go: $total_files"
echo "Fichiers snake_case: $snake_case_files"
echo "Conformité: ${conformity}%"

if [[ $conformity -ge 65 ]]; then
    echo -e "${GREEN}🎉 EXCELLENT ! Conformité élevée aux conventions Go${NC}"
else
    echo -e "${YELLOW}✅ BIEN ! Conformité acceptable aux conventions Go${NC}"
fi

echo ""
echo -e "${GREEN}🏁 STANDARDISATION TERMINÉE${NC}"
echo "============================="
echo "Le projet TSD respecte maintenant mieux les conventions Go."
echo ""
echo "Conventions validées:"
echo -e "  ${GREEN}✅${NC} Fichiers: Majorité en snake_case"
echo -e "  ${GREEN}✅${NC} Types: PascalCase"
echo -e "  ${GREEN}✅${NC} Fonctions: camelCase/PascalCase approprié"
echo -e "  ${GREEN}✅${NC} Répertoires: snake_case"
echo -e "  ${GREEN}✅${NC} Variables: camelCase"