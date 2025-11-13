#!/bin/bash

# Script de renommage automatique pour standardiser les noms de fichiers Go
# Usage: ./rename_files_to_snake_case.sh

echo "🔄 RENOMMAGE AUTOMATIQUE VERS SNAKE_CASE"
echo "========================================"

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Compteurs
renamed_count=0
skipped_count=0

# Liste des fichiers à renommer (camelCase vers snake_case)
declare -A file_renames=(
    # constraint/
    ["./constraint/api.go"]="./constraint/constraint_api.go"
    
    # rete/ - fichiers principaux
    ["./rete/converter.go"]="./rete/type_converter.go"
    ["./rete/evaluator.go"]="./rete/expression_evaluator.go"
    ["./rete/network.go"]="./rete/rete_network.go"
    ["./rete/rete.go"]="./rete/rete_core.go"
    
    # pkg/domain - fichiers core
    ["./constraint/pkg/domain/types.go"]="./constraint/pkg/domain/constraint_types.go"
    ["./constraint/pkg/domain/errors.go"]="./constraint/pkg/domain/constraint_errors.go"
    ["./constraint/pkg/domain/interfaces.go"]="./constraint/pkg/domain/constraint_interfaces.go"
    ["./rete/pkg/domain/facts.go"]="./rete/pkg/domain/fact_types.go"
    ["./rete/pkg/domain/interfaces.go"]="./rete/pkg/domain/rete_interfaces.go"
    ["./rete/pkg/domain/errors.go"]="./rete/pkg/domain/rete_errors.go"
    
    # pkg/nodes
    ["./rete/pkg/nodes/base.go"]="./rete/pkg/nodes/node_base.go"
    ["./rete/pkg/nodes/beta.go"]="./rete/pkg/nodes/beta_node.go"
    
    # pkg/validator
    ["./constraint/pkg/validator/validator.go"]="./constraint/pkg/validator/constraint_validator.go"
    ["./constraint/pkg/validator/types.go"]="./constraint/pkg/validator/validator_types.go"
    
    # pkg/storage
    ["./rete/pkg/storage/storage.go"]="./rete/pkg/storage/memory_storage.go"
    
    # internal/config
    ["./constraint/internal/config/config.go"]="./constraint/internal/config/constraint_config.go"
    ["./rete/internal/config/config.go"]="./rete/internal/config/rete_config.go"
    
    # constraint/ - fichiers root
    ["./constraint/parser.go"]="./constraint/constraint_parser.go"
    
    # test/ - helper global
    ["./test/helper.go"]="./test/test_helper.go"
)

echo -e "${BLUE}📋 PLAN DE RENOMMAGE${NC}"
echo "==================="

for old_file in "${!file_renames[@]}"; do
    new_file="${file_renames[$old_file]}"
    if [[ -f "$old_file" ]]; then
        echo -e "${YELLOW}📝 $old_file → $new_file${NC}"
    else
        echo -e "${RED}❌ Fichier non trouvé: $old_file${NC}"
    fi
done

echo ""
read -p "Continuer avec le renommage ? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Annulé."
    exit 0
fi

echo -e "${GREEN}🚀 DÉBUT DU RENOMMAGE${NC}"
echo "======================="

for old_file in "${!file_renames[@]}"; do
    new_file="${file_renames[$old_file]}"
    
    if [[ -f "$old_file" ]]; then
        # Créer le répertoire de destination si nécessaire
        new_dir=$(dirname "$new_file")
        mkdir -p "$new_dir"
        
        # Renommer le fichier
        if mv "$old_file" "$new_file" 2>/dev/null; then
            echo -e "${GREEN}✅ Renommé: $old_file → $new_file${NC}"
            renamed_count=$((renamed_count + 1))
        else
            echo -e "${RED}❌ Erreur renommage: $old_file${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️ Ignoré (fichier absent): $old_file${NC}"
        skipped_count=$((skipped_count + 1))
    fi
done

echo ""
echo -e "${BLUE}📊 STATISTIQUES RENOMMAGE${NC}"
echo "============================"
echo "Fichiers renommés: $renamed_count"
echo "Fichiers ignorés: $skipped_count"

echo ""
echo -e "${YELLOW}⚠️ ATTENTION: MISE À JOUR DES IMPORTS NÉCESSAIRE${NC}"
echo "=================================================="
echo "Les imports dans les fichiers Go devront être mis à jour manuellement."
echo "Utilisez 'go mod tidy' pour nettoyer les dépendances après mise à jour."

echo ""
echo -e "${GREEN}✅ Renommage terminé !${NC}"
echo "Prochaines étapes:"
echo "1. Mettre à jour les imports dans les fichiers Go"
echo "2. Exécuter 'go mod tidy'"
echo "3. Tester la compilation"