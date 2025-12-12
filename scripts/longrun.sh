#!/bin/bash
# Long Run Script - Exécution séquentielle avec GitHub Copilot CLI
# Applique review.md sur chaque fichier *.md d'un sous-répertoire
set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Vérifier que copilot CLI est installé
check_copilot() {
    if ! command -v copilot &> /dev/null; then
        echo -e "${RED}❌ GitHub Copilot CLI non installé${NC}"
        echo ""
        echo "Installation :"
        echo "  npm install -g @githubnext/github-copilot-cli"
        echo ""
        echo "Configuration :"
        echo "  copilot auth login"
        exit 1
    fi
}

# Lister les fichiers *.md par ordre lexicographique dans le sous-répertoire
get_session_files() {
    local subdir="$1"
    find "$SCRIPT_DIR/$subdir" -maxdepth 1 -name "*.md" -type f | sort
}

# Exécuter review pour une session
run_session_review() {
    local session_file="$1"
    local session_name=$(basename "$session_file")

    echo ""
    echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
    echo -e "${YELLOW}  REVIEW SESSION : $session_name${NC}"
    echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
    echo ""

    # Construire le prompt
    local prompt="Execute, as the linux user resinsec, the prompt .github/prompts/review.md (de l'analyse jusqu'au refactoring du code que tu dois mener en appliquant l'ensemble des préconisations et solutions identifiées) en l'appliquant sur le périmètre et les contraintes définis dans ${session_file} ainsi que les règles et bonnes pratiques définies dans .github/prompts/common.md. Effectue les modifications sans conservation de l'existant même si elles impliquent une modification du code qui utilise cet existant. Dans le cas où le nouveau code ne serait pas compatible avec l'existant, si tu ne peux corriger le code appelant, décris clairement en TODO les actions qui seront nécessaires pour rendre fonctionnel le code qui utilisera les modifications faites."

    echo -e "${BLUE}📝 Prompt : ${prompt:0:100}...${NC}"
    echo -e "${BLUE}🚀 Lancement Copilot CLI...${NC}"
    echo ""

    # Exécuter copilot
    cd "$PROJECT_ROOT"

    if copilot -p "$prompt" --allow-all-tools; then
        echo ""
        echo -e "${GREEN}✅ Session $session_name terminée avec succès${NC}"
        return 0
    else
        echo ""
        echo -e "${RED}❌ Session $session_name échouée${NC}"
        return 1
    fi
}

# Fonction principale
main() {
    local subdir="$1"

    if [ -z "$subdir" ]; then
        echo -e "${RED}❌ Erreur : vous devez spécifier un sous-répertoire${NC}"
        echo ""
        show_help
        exit 1
    fi

    if [ ! -d "$SCRIPT_DIR/$subdir" ]; then
        echo -e "${RED}❌ Erreur : le sous-répertoire '$subdir' n'existe pas dans $SCRIPT_DIR${NC}"
        exit 1
    fi

    echo -e "${BLUE}╔═══════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║   REVIEW AUTOMATISÉE - COPILOT CLI       ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════╝${NC}"
    echo ""

    # Vérifications préliminaires
    check_copilot

    echo -e "${GREEN}✅ Copilot CLI détecté${NC}"
    echo -e "${BLUE}📂 Projet : $PROJECT_ROOT${NC}"
    echo -e "${BLUE}📁 Scripts : $SCRIPT_DIR${NC}"
    echo -e "${BLUE}📁 Sous-répertoire : $subdir${NC}"
    echo ""

    # Récupérer les fichiers session
    local session_files=($(get_session_files "$subdir"))
    local total_sessions=${#session_files[@]}

    if [ $total_sessions -eq 0 ]; then
        echo -e "${RED}❌ Aucun fichier *.md trouvé dans $SCRIPT_DIR/$subdir${NC}"
        exit 1
    fi

    echo -e "${CYAN}📋 Sessions trouvées : $total_sessions${NC}"
    for file in "${session_files[@]}"; do
        echo -e "${CYAN}   - $(basename "$file")${NC}"
    done
    echo ""

    # Demander confirmation
    if [ "${AUTO_CONFIRM:-0}" != "1" ]; then
        read -p "Lancer les $total_sessions sessions ? (y/N) " -n 1 -r REPLY
        echo ""
        echo $REPLY
        echo "lu"
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Annulé."
            exit 0
        fi
    fi

    echo "execution : start"

    # Exécuter chaque session
    local success_count=0
    local fail_count=0
    local current=0

    echo "files : " "${session_files[@]}"

    for session_file in "${session_files[@]}"; do
        current=$((current + 1))
        echo "session file :" $session_file "."

        echo -e "${CYAN}[${current}/${total_sessions}]${NC}"

        if run_session_review "$session_file"; then
             echo "true"
             success_count=$((success_count + 1))
        else
             echo "false"
             fail_count=$((fail_count + 1))

            # Demander si continuer après échec
            if [ "${AUTO_CONTINUE:-0}" != "1" ]; then
                echo ""
                read -p "Continuer malgré l'erreur ? (y/N) " -n 1 -r REPLY
                echo ""
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    echo "Arrêt."
                    break
                fi
            else
                echo -e "${YELLOW}⚠️  Continuer malgré l'erreur (AUTO_CONTINUE=1)...${NC}"
            fi
        fi
        echo "execution : end"

        if [ $current -lt $total_sessions ]; then
            local pause_time=${PAUSE_SECONDS:-10}
            echo ""
            echo -e "${BLUE}⏸️  Pause ${pause_time}s avant session suivante...${NC}"
            sleep $pause_time
        fi
   done

    # Résumé final
    echo ""
    echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
    echo -e "${YELLOW}  RÉSUMÉ${NC}"
    echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}Total sessions : $total_sessions${NC}"
    echo -e "${GREEN}✅ Réussies : $success_count${NC}"
    if [ $fail_count -gt 0 ]; then
        echo -e "${RED}❌ Échouées : $fail_count${NC}"
    fi
    echo ""

    if [ $fail_count -eq 0 ]; then
        echo -e "${GREEN}🎉 Review complète terminée avec succès !${NC}"
        exit 0
    else
        echo -e "${YELLOW}⚠️  Review terminée avec des erreurs${NC}"
        exit 1
    fi
}

# Afficher aide
show_help() {
    echo "Usage: $0 <sous-répertoire> [OPTIONS]"
    echo ""
    echo "Arguments:"
    echo "  sous-répertoire         Sous-répertoire contenant les fichiers *.md à traiter"
    echo ""
    echo "Options:"
    echo "  -h, --help              Afficher cette aide"
    echo "  -y, --yes               Lancer sans confirmation (AUTO_CONFIRM=1)"
    echo "  -c, --continue          Continuer après erreur (AUTO_CONTINUE=1)"
    echo "  -p, --pause SECONDS     Pause entre sessions (défaut: 10s)"
    echo ""
    echo "Variables d'environnement:"
    echo "  AUTO_CONFIRM=1          Lancer sans confirmation"
    echo "  AUTO_CONTINUE=1         Continuer après erreur"
    echo "  PAUSE_SECONDS=N         Pause entre sessions (secondes)"
    echo ""
    echo "Exemples:"
    echo "  $0 mon_dossier          Mode interactif"
    echo "  $0 mon_dossier -y       Automatique sans confirmation"
    echo "  $0 mon_dossier -y -c    Automatique + continue sur erreur"
    echo "  $0 mon_dossier -p 30    Pause 30s entre sessions"
    echo ""
    echo "  AUTO_CONFIRM=1 $0 mon_dossier       Variable d'environnement"
}

# Parser les arguments
parse_args() {
    local subdir=""

    # Premier argument doit être le sous-répertoire
    if [[ $# -eq 0 ]] || [[ "$1" == -* ]]; then
        echo -e "${RED}❌ Erreur : vous devez spécifier un sous-répertoire${NC}"
        echo ""
        show_help
        exit 1
    fi

    subdir="$1"
    shift

    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -y|--yes)
                export AUTO_CONFIRM=1
                shift
                ;;
            -c|--continue)
                export AUTO_CONTINUE=1
                shift
                ;;
            -p|--pause)
                export PAUSE_SECONDS="$2"
                shift 2
                ;;
            *)
                echo -e "${RED}❌ Option inconnue : $1${NC}"
                echo ""
                show_help
                exit 1
                ;;
        esac
    done

    echo "$subdir"
}

# Point d'entrée
if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

SUBDIR=$(parse_args "$@")
main "$SUBDIR"
