#!/bin/bash

# ═══════════════════════════════════════════════════════════════════════════
# Script de génération de rapport d'exécution E2E pour xuple-spaces
# ═══════════════════════════════════════════════════════════════════════════

set -e

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Vérifier les arguments
if [ $# -lt 1 ]; then
    echo "Usage: $0 <fichier.tsd>"
    exit 1
fi

TSD_FILE="$1"

if [ ! -f "$TSD_FILE" ]; then
    echo -e "${RED}❌ Fichier non trouvé: $TSD_FILE${NC}"
    exit 1
fi

# Déterminer le chemin du binaire TSD
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
TSD_BIN="$PROJECT_DIR/bin/tsd"

if [ ! -f "$TSD_BIN" ]; then
    echo -e "${RED}❌ Binaire TSD non trouvé: $TSD_BIN${NC}"
    echo -e "${YELLOW}💡 Exécutez 'make build' pour compiler TSD${NC}"
    exit 1
fi

# Fonction pour extraire les informations du fichier TSD
extract_info() {
    local file="$1"

    # Compter les types
    TYPE_COUNT=$(grep -c "^type " "$file" 2>/dev/null || echo "0")

    # Compter les xuple-spaces
    XUPLESPACE_COUNT=$(grep -c "^xuple-space " "$file" 2>/dev/null || echo "0")

    # Compter les actions
    ACTION_COUNT=$(grep -c "^action " "$file" 2>/dev/null || echo "0")

    # Compter les règles
    RULE_COUNT=$(grep -c "^rule " "$file" 2>/dev/null || echo "0")

    # Extraire les types
    TYPES=$(grep "^type " "$file" | sed 's/type \([^(]*\).*/\1/' 2>/dev/null || echo "")

    # Extraire les xuple-spaces
    XUPLESPACES=$(grep "^xuple-space " "$file" | sed 's/xuple-space \([^ ]*\).*/\1/' 2>/dev/null || echo "")

    # Extraire les actions
    ACTIONS=$(grep "^action " "$file" | sed 's/action \([^(]*\).*/\1/' 2>/dev/null || echo "")

    # Extraire les règles
    RULES=$(grep "^rule " "$file" | sed 's/rule \([^:]*\).*/\1/' 2>/dev/null || echo "")
}

# En-tête du rapport
print_header() {
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}  RAPPORT D'EXÉCUTION E2E - SYSTÈME XUPLE-SPACE${NC}"
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}📁 Fichier:${NC} $TSD_FILE"
    echo -e "${CYAN}⏰ Date:${NC} $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
}

# Section 1: Structure du programme
print_structure() {
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${BOLD}SECTION 1 : STRUCTURE DU PROGRAMME${NC}"
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo ""

    # Types
    echo -e "${MAGENTA}📋 TYPES DÉFINIS (${TYPE_COUNT}):${NC}"
    if [ -n "$TYPES" ]; then
        echo "$TYPES" | while IFS= read -r type; do
            if [ -n "$type" ]; then
                echo -e "  ${GREEN}•${NC} $type"
                # Extraire les champs du type
                grep "^type $type" "$TSD_FILE" | sed 's/.*(\(.*\))/\1/' | tr ',' '\n' | sed 's/^[ \t]*/    /' | sed 's/#/🔑 /'
            fi
        done
    fi
    echo ""

    # Xuple-spaces
    echo -e "${BLUE}🗄️  XUPLE-SPACES DÉCLARÉS (${XUPLESPACE_COUNT}):${NC}"
    if [ -n "$XUPLESPACES" ]; then
        echo "$XUPLESPACES" | while IFS= read -r xs; do
            if [ -n "$xs" ]; then
                echo -e "  ${GREEN}•${NC} ${BOLD}$xs${NC}"
                # Extraire les politiques
                awk "/^xuple-space $xs/,/^}/" "$TSD_FILE" | grep -E "selection:|consumption:|retention:" | sed 's/^/    /'
            fi
        done
    fi
    echo ""

    # Actions
    echo -e "${YELLOW}⚡ ACTIONS DÉFINIES (${ACTION_COUNT}):${NC}"
    if [ -n "$ACTIONS" ]; then
        echo "$ACTIONS" | while IFS= read -r action; do
            if [ -n "$action" ]; then
                grep "^action $action" "$TSD_FILE" | sed "s/action/${GREEN}•${NC}/"
            fi
        done
    fi
    echo ""

    # Règles
    echo -e "${CYAN}📜 RÈGLES DÉFINIES (${RULE_COUNT}):${NC}"
    if [ -n "$RULES" ]; then
        echo "$RULES" | while IFS= read -r rule; do
            if [ -n "$rule" ]; then
                echo -e "  ${GREEN}•${NC} $rule"
                # Extraire la condition et les actions
                grep "^rule $rule" "$TSD_FILE" | sed 's/.*\/ /    Condition: /' | sed 's/ ==>.*//'
                grep "^rule $rule" "$TSD_FILE" | sed 's/.*==>/    Actions: /'
            fi
        done
    fi
    echo ""
}

# Section 2: Faits injectés
print_facts() {
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${BOLD}SECTION 2 : FAITS INJECTÉS${NC}"
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo ""

    # Compter et afficher les faits par type
    if [ -n "$TYPES" ]; then
        echo "$TYPES" | while IFS= read -r type; do
            if [ -n "$type" ]; then
                FACT_COUNT=$(grep "^$type(" "$TSD_FILE" | wc -l)
                if [ "$FACT_COUNT" -gt 0 ]; then
                    echo -e "${MAGENTA}📊 $type (${FACT_COUNT} fait(s)):${NC}"
                    grep "^$type(" "$TSD_FILE" | nl -w2 -s'. ' | sed 's/^/  /'
                    echo ""
                fi
            fi
        done
    fi
}

# Section 3: Exécution
print_execution() {
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${BOLD}SECTION 3 : EXÉCUTION DU PROGRAMME${NC}"
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo ""

    echo -e "${CYAN}🔄 Exécution de TSD...${NC}"
    echo ""

    # Exécuter TSD et capturer la sortie
    TSD_OUTPUT=$("$TSD_BIN" "$TSD_FILE" -v 2>&1)
    TSD_EXIT_CODE=$?

    if [ $TSD_EXIT_CODE -eq 0 ]; then
        echo -e "${GREEN}✅ Exécution réussie${NC}"
        echo ""

        # Afficher la sortie de TSD (filtrée)
        echo -e "${BOLD}Sortie de l'exécution:${NC}"
        echo "────────────────────────────────────────────────────────────────────────────"
        echo "$TSD_OUTPUT" | grep -v "^\[" | sed 's/^/  /'
        echo "────────────────────────────────────────────────────────────────────────────"
        echo ""

        # Extraire les statistiques si disponibles
        if echo "$TSD_OUTPUT" | grep -q "ACTION"; then
            echo -e "${YELLOW}📋 Actions déclenchées:${NC}"
            echo "$TSD_OUTPUT" | grep "ACTION" | sed 's/^/  /'
            echo ""
        fi
    else
        echo -e "${RED}❌ Erreur lors de l'exécution${NC}"
        echo ""
        echo "$TSD_OUTPUT"
        echo ""
        return 1
    fi
}

# Section 4: Xuples générés (simulation)
print_xuples() {
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${BOLD}SECTION 4 : XUPLES GÉNÉRÉS (ATTENDUS)${NC}"
    echo -e "${BOLD}───────────────────────────────────────────────────────────────────────────${NC}"
    echo ""

    echo -e "${BLUE}🎯 Analyse des xuples potentiels basée sur les règles:${NC}"
    echo ""

    # Pour chaque xuple-space, on peut inférer quels xuples seraient créés
    # basé sur les règles qui utilisent l'action Xuple
    if grep -q "Xuple(" "$TSD_FILE" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Des actions Xuple ont été détectées dans les règles"
        echo ""

        # Lister les xuple-spaces qui reçoivent des xuples
        if [ -n "$XUPLESPACES" ]; then
            echo "$XUPLESPACES" | while IFS= read -r xs; do
                if [ -n "$xs" ]; then
                    XUPLE_RULES=$(grep -c "Xuple(\"$xs\"" "$TSD_FILE" 2>/dev/null || echo "0")
                    if [ "$XUPLE_RULES" -gt 0 ]; then
                        echo -e "${MAGENTA}📦 Xuple-space: ${BOLD}$xs${NC}"
                        echo -e "   Règles qui y écrivent: $XUPLE_RULES"
                        grep "Xuple(\"$xs\"" "$TSD_FILE" | sed 's/^/   /'
                        echo ""
                    fi
                fi
            done
        fi
    else
        echo -e "${YELLOW}⚠️  Aucune action Xuple détectée dans ce programme${NC}"
        echo -e "   Les xuple-spaces sont déclarés mais non utilisés"
        echo ""
    fi
}

# Résumé final
print_summary() {
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}RÉSUMÉ FINAL${NC}"
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════════${NC}"
    echo ""

    echo -e "${GREEN}✓${NC} Types définis: ${BOLD}$TYPE_COUNT${NC}"
    echo -e "${GREEN}✓${NC} Xuple-spaces déclarés: ${BOLD}$XUPLESPACE_COUNT${NC}"
    echo -e "${GREEN}✓${NC} Actions définies: ${BOLD}$ACTION_COUNT${NC}"
    echo -e "${GREEN}✓${NC} Règles définies: ${BOLD}$RULE_COUNT${NC}"

    # Compter le total de faits
    TOTAL_FACTS=0
    if [ -n "$TYPES" ]; then
        while IFS= read -r type; do
            if [ -n "$type" ]; then
                FACT_COUNT=$(grep "^$type(" "$TSD_FILE" 2>/dev/null | wc -l)
                TOTAL_FACTS=$((TOTAL_FACTS + FACT_COUNT))
            fi
        done <<< "$TYPES"
    fi
    echo -e "${GREEN}✓${NC} Faits injectés: ${BOLD}$TOTAL_FACTS${NC}"

    echo ""
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════════${NC}"
}

# ═══════════════════════════════════════════════════════════════════════════
# MAIN
# ═══════════════════════════════════════════════════════════════════════════

# Extraire les informations du fichier
extract_info "$TSD_FILE"

# Afficher les sections
print_header
print_structure
print_facts
print_execution
print_xuples
print_summary

echo ""
echo -e "${CYAN}📄 Rapport généré avec succès${NC}"
echo ""
