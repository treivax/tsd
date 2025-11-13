#!/bin/bash

# Script de renommage conservateur - Phase 1
# Renomme seulement les fichiers non-critiques pour éviter de casser les imports

echo "🔄 RENOMMAGE CONSERVATEUR - PHASE 1"
echo "================================="

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 FICHIERS À RENOMMER (Phase 1 - Non-critiques)${NC}"
echo "==============================================="

# Renommer seulement les fichiers test helper qui ont moins de dépendances
declare -A safe_renames=(
    # Helper tests - moins de dépendances
    ["./test/helper.go"]="./test/test_util.go"
)

for old_file in "${!safe_renames[@]}"; do
    new_file="${safe_renames[$old_file]}"
    if [[ -f "$old_file" ]]; then
        echo -e "${YELLOW}📝 $old_file → $new_file${NC}"
    fi
done

echo ""
read -p "Continuer ? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Annulé."
    exit 0
fi

echo -e "${GREEN}🚀 RENOMMAGE EN COURS${NC}"

for old_file in "${!safe_renames[@]}"; do
    new_file="${safe_renames[$old_file]}"
    
    if [[ -f "$old_file" ]]; then
        if mv "$old_file" "$new_file"; then
            echo -e "${GREEN}✅ Renommé: $old_file → $new_file${NC}"
        else
            echo -e "${RED}❌ Erreur: $old_file${NC}"
        fi
    fi
done

echo -e "${GREEN}✅ Phase 1 terminée${NC}"