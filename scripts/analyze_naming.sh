#!/bin/bash

# Script d'analyse des conventions de nommage dans le projet TSD
# Usage: ./analyze_naming.sh

echo "🔍 ANALYSE DES CONVENTIONS DE NOMMAGE TSD"
echo "=========================================="

# Couleurs pour l'output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Compteurs
total_files=0
camel_case_files=0
snake_case_files=0
mixed_case_files=0

echo -e "${BLUE}📁 ANALYSE DES NOMS DE FICHIERS GO${NC}"
echo "======================================"

# Analyser tous les fichiers .go
for file in $(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*"); do
    total_files=$((total_files + 1))
    basename_file=$(basename "$file" .go)
    
    # Détecter le style de nommage
    if [[ "$basename_file" =~ ^[a-z][a-zA-Z0-9]*$ ]]; then
        # camelCase ou un seul mot en minuscules
        camel_case_files=$((camel_case_files + 1))
        echo -e "${GREEN}✓ camelCase:${NC} $file"
    elif [[ "$basename_file" =~ ^[a-z][a-z0-9]*(_[a-z0-9]+)+$ ]]; then
        # snake_case
        snake_case_files=$((snake_case_files + 1))
        echo -e "${YELLOW}⚠ snake_case:${NC} $file"
    elif [[ "$basename_file" =~ ^[A-Z] ]]; then
        # PascalCase (généralement pour types)
        echo -e "${BLUE}ℹ PascalCase:${NC} $file"
    else
        # Cas mixtes ou autres
        mixed_case_files=$((mixed_case_files + 1))
        echo -e "${RED}❌ Mixed/Other:${NC} $file"
    fi
done

echo ""
echo -e "${BLUE}📊 STATISTIQUES NOMS DE FICHIERS${NC}"
echo "=================================="
echo "Total fichiers Go: $total_files"
echo -e "${GREEN}CamelCase: $camel_case_files${NC}"
echo -e "${YELLOW}Snake_case: $snake_case_files${NC}"
echo -e "${RED}Mixed/Other: $mixed_case_files${NC}"

echo ""
echo -e "${BLUE}🔤 ANALYSE DES FONCTIONS ET TYPES DANS LE CODE${NC}"
echo "============================================="

# Analyser les déclarations dans les fichiers Go
echo "Recherche de fonctions/types avec noms non-conformes..."

# Fonctions qui devraient être camelCase
echo -e "\n${YELLOW}🔧 Fonctions avec snake_case (devraient être camelCase):${NC}"
grep -n "func [a-z][a-z0-9]*_[a-z0-9_]*(" $(find . -name "*.go" -not -path "./vendor/*") | head -20

# Types qui devraient être PascalCase
echo -e "\n${YELLOW}📦 Types avec snake_case (devraient être PascalCase):${NC}"
grep -n "type [a-z][a-z0-9]*_[a-z0-9_]* " $(find . -name "*.go" -not -path "./vendor/*") | head -20

# Variables/constantes avec mixed case
echo -e "\n${YELLOW}🔀 Variables/constantes avec nommage mixte:${NC}"
grep -n "var [A-Z][a-zA-Z]*_" $(find . -name "*.go" -not -path "./vendor/*") | head -10
grep -n "const [A-Z][a-zA-Z]*_" $(find . -name "*.go" -not -path "./vendor/*") | head -10

echo ""
echo -e "${BLUE}📂 ANALYSE DES NOMS DE RÉPERTOIRES${NC}"
echo "================================="

# Analyser les noms de répertoires
for dir in $(find . -type d -not -path "./.git*" -not -path "./vendor/*" | sort); do
    basename_dir=$(basename "$dir")
    
    if [[ "$basename_dir" == "." ]]; then
        continue
    fi
    
    if [[ "$basename_dir" =~ ^[a-z][a-z0-9]*(_[a-z0-9]+)*$ ]]; then
        echo -e "${GREEN}✓ snake_case:${NC} $dir"
    elif [[ "$basename_dir" =~ ^[a-z][a-zA-Z0-9]*$ ]]; then
        echo -e "${YELLOW}⚠ camelCase:${NC} $dir"
    else
        echo -e "${RED}❌ Mixed/Other:${NC} $dir"
    fi
done

echo ""
echo -e "${BLUE}🎯 RECOMMANDATIONS GOLANG${NC}"
echo "==========================="
echo "✅ Fichiers Go: snake_case (ex: user_service.go)"
echo "✅ Types/Structs: PascalCase (ex: UserService)"
echo "✅ Fonctions: camelCase (ex: getUserInfo)"
echo "✅ Constantes exportées: UPPER_SNAKE_CASE"
echo "✅ Variables: camelCase (ex: userName)"
echo "✅ Répertoires: snake_case (ex: user_service/)"

echo ""
echo -e "${BLUE}📋 PROCHAINES ÉTAPES${NC}"
echo "==================="
echo "1. Renommer les fichiers non-conformes en snake_case"
echo "2. Standardiser les noms de fonctions en camelCase"
echo "3. Vérifier que les types sont en PascalCase"
echo "4. Mettre à jour les imports après renommage"

echo ""
echo "🏁 Analyse terminée."