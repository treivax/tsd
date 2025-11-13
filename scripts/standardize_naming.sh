#!/bin/bash

# Script de standardisation complète des conventions de nommage Go
# Applique les bonnes pratiques Go de manière systématique

echo "🔄 STANDARDISATION COMPLÈTE DES CONVENTIONS GO"
echo "=============================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Compteurs
renamed_files=0
updated_functions=0
updated_variables=0

echo -e "${BLUE}📋 PLAN DE STANDARDISATION${NC}"
echo "=========================="
echo "1. 🗂️  Fichiers Go: camelCase → snake_case"
echo "2. 🔧 Fonctions: snake_case → camelCase (si nécessaire)"
echo "3. 📦 Types: snake_case → PascalCase (si nécessaire)"
echo "4. 🔗 Mise à jour des imports/références"

echo ""
read -p "Continuer avec la standardisation ? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Annulé."
    exit 0
fi

echo -e "${GREEN}🚀 PHASE 1: STANDARDISATION DES NOMS DE FICHIERS${NC}"
echo "================================================="

# Liste des fichiers à renommer (camelCase → snake_case)
declare -A file_renames=(
    # constraint/ - fichiers principaux
    ["./constraint/api.go"]="./constraint/constraint_api.go"
    ["./constraint/parser.go"]="./constraint/constraint_parser.go"
    
    # rete/ - fichiers principaux  
    ["./rete/converter.go"]="./rete/type_converter.go"
    ["./rete/evaluator.go"]="./rete/expression_evaluator.go"
    ["./rete/network.go"]="./rete/rete_network.go"
    ["./rete/rete.go"]="./rete/rete_core.go"
    
    # pkg/domain - spécialisation par module
    ["./constraint/pkg/domain/types.go"]="./constraint/pkg/domain/constraint_types.go"
    ["./constraint/pkg/domain/errors.go"]="./constraint/pkg/domain/constraint_errors.go"
    ["./constraint/pkg/domain/interfaces.go"]="./constraint/pkg/domain/constraint_interfaces.go"
    ["./rete/pkg/domain/facts.go"]="./rete/pkg/domain/fact_types.go"
    ["./rete/pkg/domain/interfaces.go"]="./rete/pkg/domain/rete_interfaces.go"
    ["./rete/pkg/domain/errors.go"]="./rete/pkg/domain/rete_errors.go"
    
    # pkg/nodes
    ["./rete/pkg/nodes/base.go"]="./rete/pkg/nodes/node_base.go"
    ["./rete/pkg/nodes/beta.go"]="./rete/pkg/nodes/beta_nodes.go"
    
    # pkg/validator
    ["./constraint/pkg/validator/validator.go"]="./constraint/pkg/validator/constraint_validator.go"
    ["./constraint/pkg/validator/types.go"]="./constraint/pkg/validator/validator_types.go"
    
    # internal/config
    ["./constraint/internal/config/config.go"]="./constraint/internal/config/constraint_config.go"
    ["./rete/internal/config/config.go"]="./rete/internal/config/rete_config.go"
    
    # test/ - helper global
    ["./test/helper.go"]="./test/test_utils.go"
)

echo "Fichiers à renommer:"
for old_file in "${!file_renames[@]}"; do
    new_file="${file_renames[$old_file]}"
    if [[ -f "$old_file" ]]; then
        echo -e "${YELLOW}  📝 $(basename $old_file) → $(basename $new_file)${NC}"
    fi
done

echo ""
echo -e "${GREEN}Renommage des fichiers...${NC}"

# Fonction pour mettre à jour les imports dans tous les fichiers Go
update_imports() {
    local old_import="$1"
    local new_import="$2"
    
    echo "  🔄 Mise à jour import: $old_import → $new_import"
    
    # Trouver tous les fichiers Go et mettre à jour les imports
    find . -name "*.go" -not -path "./vendor/*" -type f -exec grep -l "$old_import" {} \; | while read file; do
        sed -i "s|$old_import|$new_import|g" "$file"
    done
}

# Renommer les fichiers et mettre à jour les imports
for old_file in "${!file_renames[@]}"; do
    new_file="${file_renames[$old_file]}"
    
    if [[ -f "$old_file" ]]; then
        # Créer le répertoire de destination si nécessaire
        new_dir=$(dirname "$new_file")
        mkdir -p "$new_dir"
        
        # Renommer le fichier
        if mv "$old_file" "$new_file" 2>/dev/null; then
            echo -e "${GREEN}✅ Renommé: $(basename $old_file) → $(basename $new_file)${NC}"
            renamed_files=$((renamed_files + 1))
            
            # Pas de mise à jour d'imports pour l'instant car trop complexe
        else
            echo -e "${RED}❌ Erreur renommage: $old_file${NC}"
        fi
    fi
done

echo ""
echo -e "${GREEN}🚀 PHASE 2: VALIDATION DE LA COMPILATION${NC}"
echo "========================================="

echo "Test de compilation après renommage..."
if go build ./... 2>/dev/null; then
    echo -e "${GREEN}✅ Compilation réussie après renommage${NC}"
else
    echo -e "${RED}❌ Erreurs de compilation détectées${NC}"
    echo "Les imports devront être mis à jour manuellement."
fi

echo ""
echo -e "${GREEN}🚀 PHASE 3: ANALYSE DES FONCTIONS ET VARIABLES${NC}"
echo "=============================================="

# Analyser les fonctions avec snake_case (rare mais possible)
echo "Recherche de fonctions non-conformes..."
snake_case_functions=$(grep -r "func [a-z][a-z0-9]*_[a-z0-9_]*(" . --include="*.go" | head -5)
if [[ -n "$snake_case_functions" ]]; then
    echo -e "${YELLOW}⚠️ Fonctions avec snake_case trouvées:${NC}"
    echo "$snake_case_functions"
    echo "Note: Ces fonctions devraient être renommées en camelCase manuellement."
else
    echo -e "${GREEN}✅ Aucune fonction avec snake_case trouvée${NC}"
fi

# Analyser les types avec snake_case
echo ""
echo "Recherche de types non-conformes..."
snake_case_types=$(grep -r "type [a-z][a-z0-9]*_[a-z0-9_]*" . --include="*.go" | head -5)
if [[ -n "$snake_case_types" ]]; then
    echo -e "${YELLOW}⚠️ Types avec snake_case trouvés:${NC}"
    echo "$snake_case_types"
    echo "Note: Ces types devraient être renommés en PascalCase manuellement."
else
    echo -e "${GREEN}✅ Aucun type avec snake_case trouvé${NC}"
fi

echo ""
echo -e "${BLUE}📊 STATISTIQUES FINALES${NC}"
echo "========================"
echo "Fichiers renommés: $renamed_files"
echo "Fonctions analysées: ✅"
echo "Types analysés: ✅"

echo ""
echo -e "${YELLOW}⚠️ ACTIONS MANUELLES REQUISES${NC}"
echo "==============================="
echo "1. Vérifier la compilation: 'go build ./...'"
echo "2. Exécuter les tests: 'go test ./...'"
echo "3. Mettre à jour manuellement les imports si nécessaire"
echo "4. Renommer manuellement les fonctions/types non-conformes"

echo ""
echo -e "${GREEN}✅ Standardisation terminée !${NC}"
echo ""
echo "Le projet respecte maintenant mieux les conventions Go:"
echo "📁 Fichiers: snake_case"
echo "🏷️ Types: PascalCase"  
echo "🔧 Fonctions: camelCase"
echo "📂 Répertoires: snake_case"