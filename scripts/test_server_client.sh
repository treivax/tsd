#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script de test pour le serveur et le client TSD
# Ce script démarre le serveur, exécute des tests avec le client, puis arrête le serveur

set -e

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVER_PORT=8080
SERVER_HOST="localhost"
SERVER_URL="http://${SERVER_HOST}:${SERVER_PORT}"
SERVER_PID=""
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${PROJECT_ROOT}/bin"
EXAMPLES_DIR="${PROJECT_ROOT}/examples/server"

# Fonction pour afficher des messages
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Fonction pour nettoyer à la fin
cleanup() {
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        log_info "Arrêt du serveur (PID: $SERVER_PID)..."
        kill "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
        log_success "Serveur arrêté"
    fi
}

trap cleanup EXIT

# Fonction pour vérifier si un port est utilisé
check_port() {
    local port=$1
    if lsof -Pi :${port} -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        return 0
    else
        return 1
    fi
}

# Fonction pour attendre que le serveur soit prêt
wait_for_server() {
    local max_attempts=30
    local attempt=0

    log_info "Attente du démarrage du serveur..."

    while [ $attempt -lt $max_attempts ]; do
        if curl -s "${SERVER_URL}/health" >/dev/null 2>&1; then
            log_success "Serveur prêt!"
            return 0
        fi

        attempt=$((attempt + 1))
        sleep 1
    done

    log_error "Timeout: le serveur n'a pas démarré dans les temps"
    return 1
}

# Fonction principale
main() {
    log_info "=========================================="
    log_info "Test du Serveur et Client TSD"
    log_info "=========================================="
    echo ""

    # Vérifier que nous sommes dans le bon répertoire
    if [ ! -f "${PROJECT_ROOT}/go.mod" ]; then
        log_error "Ce script doit être exécuté depuis le répertoire racine du projet"
        exit 1
    fi

    # Créer le répertoire bin s'il n'existe pas
    mkdir -p "${BIN_DIR}"

    # Étape 1: Compiler le serveur et le client
    log_info "Compilation du serveur et du client..."

    cd "${PROJECT_ROOT}"

    if ! go build -o "${BIN_DIR}/tsd-server" ./cmd/tsd-server; then
        log_error "Échec de la compilation du serveur"
        exit 1
    fi
    log_success "Serveur compilé"

    if ! go build -o "${BIN_DIR}/tsd-client" ./cmd/tsd-client; then
        log_error "Échec de la compilation du client"
        exit 1
    fi
    log_success "Client compilé"

    echo ""

    # Étape 2: Vérifier que le port est libre
    if check_port ${SERVER_PORT}; then
        log_warning "Le port ${SERVER_PORT} est déjà utilisé"
        log_info "Tentative d'utilisation du port 9080 à la place..."
        SERVER_PORT=9080
        SERVER_URL="http://${SERVER_HOST}:${SERVER_PORT}"

        if check_port ${SERVER_PORT}; then
            log_error "Le port ${SERVER_PORT} est également utilisé. Arrêtez les serveurs existants."
            exit 1
        fi
    fi

    # Étape 3: Démarrer le serveur
    log_info "Démarrage du serveur sur ${SERVER_URL}..."

    "${BIN_DIR}/tsd-server" -port ${SERVER_PORT} > /tmp/tsd-server.log 2>&1 &
    SERVER_PID=$!

    log_info "Serveur démarré (PID: ${SERVER_PID})"

    # Attendre que le serveur soit prêt
    if ! wait_for_server; then
        log_error "Le serveur n'a pas pu démarrer"
        cat /tmp/tsd-server.log
        exit 1
    fi

    echo ""

    # Étape 4: Tests du client
    log_info "=========================================="
    log_info "Tests du Client"
    log_info "=========================================="
    echo ""

    # Test 1: Health check
    log_info "Test 1: Health check..."
    if "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" -health; then
        log_success "Health check OK"
    else
        log_error "Health check échoué"
        exit 1
    fi
    echo ""

    # Test 2: Exécution d'un fichier simple
    log_info "Test 2: Exécution d'un fichier simple..."
    if [ -f "${EXAMPLES_DIR}/simple.tsd" ]; then
        if "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" "${EXAMPLES_DIR}/simple.tsd"; then
            log_success "Exécution du fichier simple réussie"
        else
            log_error "Exécution du fichier simple échouée"
            exit 1
        fi
    else
        log_warning "Fichier ${EXAMPLES_DIR}/simple.tsd non trouvé, test ignoré"
    fi
    echo ""

    # Test 3: Exécution avec format JSON
    log_info "Test 3: Exécution avec format JSON..."
    if [ -f "${EXAMPLES_DIR}/simple.tsd" ]; then
        if "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" -format json "${EXAMPLES_DIR}/simple.tsd" > /tmp/tsd-client-output.json; then
            log_success "Format JSON OK"

            # Vérifier que la sortie est du JSON valide
            if command -v jq >/dev/null 2>&1; then
                if jq empty /tmp/tsd-client-output.json 2>/dev/null; then
                    log_success "JSON valide"
                else
                    log_error "JSON invalide"
                    cat /tmp/tsd-client-output.json
                    exit 1
                fi
            else
                log_warning "jq non disponible, validation JSON ignorée"
            fi
        else
            log_error "Exécution avec format JSON échouée"
            exit 1
        fi
    fi
    echo ""

    # Test 4: Exécution avec stdin
    log_info "Test 4: Exécution avec stdin..."
    echo 'type Person : <id: string, name: string>
action notify : <message: string>
rule person_rule : {p: Person} / p.name == "Alice" ==> notify(p.id)
Person("p1", "Alice")' | "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" -stdin

    if [ $? -eq 0 ]; then
        log_success "Exécution avec stdin réussie"
    else
        log_error "Exécution avec stdin échouée"
        exit 1
    fi
    echo ""

    # Test 5: Exécution avec code direct
    log_info "Test 5: Exécution avec code direct..."
    if "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" -text 'type Person : <id: string, name: string>
Person("p1", "Alice")'; then
        log_success "Exécution avec code direct réussie"
    else
        log_error "Exécution avec code direct échouée"
        exit 1
    fi
    echo ""

    # Test 6: Multiples activations
    log_info "Test 6: Exécution avec multiples activations..."
    if [ -f "${EXAMPLES_DIR}/multiple_activations.tsd" ]; then
        if "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" "${EXAMPLES_DIR}/multiple_activations.tsd"; then
            log_success "Exécution avec multiples activations réussie"
        else
            log_error "Exécution avec multiples activations échouée"
            exit 1
        fi
    else
        log_warning "Fichier ${EXAMPLES_DIR}/multiple_activations.tsd non trouvé, test ignoré"
    fi
    echo ""

    # Test 7: Mode verbeux
    log_info "Test 7: Mode verbeux..."
    if [ -f "${EXAMPLES_DIR}/simple.tsd" ]; then
        if "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" -v "${EXAMPLES_DIR}/simple.tsd" > /tmp/tsd-verbose-output.txt 2>&1; then
            log_success "Mode verbeux OK"

            # Vérifier que la sortie contient des informations supplémentaires
            if grep -q "Envoi requête" /tmp/tsd-verbose-output.txt; then
                log_success "Informations verboses présentes"
            else
                log_warning "Informations verboses manquantes"
            fi
        else
            log_error "Mode verbeux échoué"
            exit 1
        fi
    fi
    echo ""

    # Test 8: Gestion des erreurs - code invalide
    log_info "Test 8: Gestion des erreurs (code invalide)..."
    if "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" -text "invalid code !!!" 2>/dev/null; then
        log_error "Le code invalide aurait dû échouer"
        exit 1
    else
        log_success "Erreur correctement gérée"
    fi
    echo ""

    # Test 9: Test de performance simple
    log_info "Test 9: Test de performance (10 requêtes)..."
    start_time=$(date +%s)

    for i in {1..10}; do
        echo 'type Person : <id: string, name: string>
Person("p1", "Alice")' | "${BIN_DIR}/tsd-client" -server "${SERVER_URL}" -stdin -format json > /dev/null 2>&1

        if [ $? -ne 0 ]; then
            log_error "Requête $i échouée"
            exit 1
        fi
    done

    end_time=$(date +%s)
    duration=$((end_time - start_time))
    log_success "10 requêtes exécutées en ${duration}s"
    echo ""

    # Résumé
    log_info "=========================================="
    log_success "TOUS LES TESTS SONT PASSÉS!"
    log_info "=========================================="
    echo ""
    log_info "Résumé:"
    log_info "  - Serveur: ${SERVER_URL}"
    log_info "  - Logs serveur: /tmp/tsd-server.log"
    log_info "  - Binaires: ${BIN_DIR}/"
    echo ""

    # Afficher les dernières lignes des logs du serveur
    log_info "Dernières lignes des logs du serveur:"
    echo "---"
    tail -n 10 /tmp/tsd-server.log
    echo "---"
    echo ""

    log_success "Script de test terminé avec succès!"
}

# Exécuter le main
main "$@"
