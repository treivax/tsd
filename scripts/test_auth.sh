#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script de test complet pour l'authentification TSD
# Ce script teste l'authentification par clé API et JWT

set -e

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Compteurs
TESTS_PASSED=0
TESTS_FAILED=0

# Fonction d'affichage
print_header() {
    echo ""
    echo -e "${BLUE}=================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}=================================${NC}"
}

print_test() {
    echo -e "${YELLOW}▶ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    ((TESTS_FAILED++))
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Fonction de nettoyage
cleanup() {
    if [ -n "$SERVER_PID" ] && kill -0 $SERVER_PID 2>/dev/null; then
        print_info "Arrêt du serveur (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
    fi

    # Nettoyer les fichiers temporaires
    rm -f /tmp/tsd_test_*.tsd
    rm -f /tmp/tsd_server_*.log
}

trap cleanup EXIT

# Configuration
BIN_DIR="./bin"
SERVER="$BIN_DIR/tsd-server"
CLIENT="$BIN_DIR/tsd-client"
AUTH_TOOL="$BIN_DIR/tsd-auth"
SERVER_PORT=18080
SERVER_URL="http://localhost:$SERVER_PORT"
SERVER_PID=""

# Vérifier que les binaires existent
print_header "Vérification des binaires"

if [ ! -f "$SERVER" ]; then
    print_error "Binaire serveur non trouvé: $SERVER"
    print_info "Compilez avec: go build -o $SERVER ./cmd/tsd-server"
    exit 1
fi
print_success "Serveur trouvé"

if [ ! -f "$CLIENT" ]; then
    print_error "Binaire client non trouvé: $CLIENT"
    print_info "Compilez avec: go build -o $CLIENT ./cmd/tsd-client"
    exit 1
fi
print_success "Client trouvé"

if [ ! -f "$AUTH_TOOL" ]; then
    print_error "Outil d'authentification non trouvé: $AUTH_TOOL"
    print_info "Compilez avec: go build -o $AUTH_TOOL ./cmd/tsd-auth"
    exit 1
fi
print_success "Outil d'authentification trouvé"

# Créer un fichier TSD de test
print_header "Préparation des fichiers de test"

TEST_FILE="/tmp/tsd_test_auth.tsd"
cat > "$TEST_FILE" << 'EOF'
type Person : <
  id: string,
  name: string,
  age: int
>

Person("p1", "Alice", 30)
Person("p2", "Bob", 25)
EOF
print_success "Fichier de test créé: $TEST_FILE"

# Test 1: Serveur sans authentification
print_header "Test 1: Serveur sans authentification"

print_test "Démarrage du serveur sans auth..."
$SERVER -port $SERVER_PORT > /tmp/tsd_server_noauth.log 2>&1 &
SERVER_PID=$!
sleep 2

if ! kill -0 $SERVER_PID 2>/dev/null; then
    print_error "Le serveur n'a pas démarré"
    cat /tmp/tsd_server_noauth.log
    exit 1
fi
print_success "Serveur démarré (PID: $SERVER_PID)"

print_test "Test health check sans auth..."
if $CLIENT -server "$SERVER_URL" -health > /dev/null 2>&1; then
    print_success "Health check OK sans authentification"
else
    print_error "Health check échoué"
fi

print_test "Test exécution sans auth..."
if $CLIENT -server "$SERVER_URL" "$TEST_FILE" > /dev/null 2>&1; then
    print_success "Exécution OK sans authentification"
else
    print_error "Exécution échouée"
fi

print_test "Arrêt du serveur..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
SERVER_PID=""
sleep 1
print_success "Serveur arrêté"

# Test 2: Auth Key
print_header "Test 2: Authentification par clé API"

print_test "Génération d'une clé API..."
API_KEY=$($AUTH_TOOL generate-key -format json 2>/dev/null | grep -o '"keys":\["[^"]*"' | cut -d'"' -f4)
if [ -z "$API_KEY" ]; then
    print_error "Génération de clé API échouée"
    exit 1
fi
print_success "Clé API générée: ${API_KEY:0:20}..."

print_test "Validation de la clé API..."
if $AUTH_TOOL validate -type key -token "$API_KEY" -keys "$API_KEY" > /dev/null 2>&1; then
    print_success "Validation de clé OK"
else
    print_error "Validation de clé échouée"
fi

print_test "Démarrage du serveur avec Auth Key..."
export TSD_AUTH_KEYS="$API_KEY"
$SERVER -port $SERVER_PORT -auth key > /tmp/tsd_server_authkey.log 2>&1 &
SERVER_PID=$!
sleep 2

if ! kill -0 $SERVER_PID 2>/dev/null; then
    print_error "Le serveur n'a pas démarré"
    cat /tmp/tsd_server_authkey.log
    exit 1
fi
print_success "Serveur démarré avec Auth Key"

print_test "Test health check SANS token (doit échouer)..."
if $CLIENT -server "$SERVER_URL" -health > /dev/null 2>&1; then
    print_error "Health check devrait échouer sans token"
else
    print_success "Health check échoué comme attendu"
fi

print_test "Test health check AVEC token..."
if $CLIENT -server "$SERVER_URL" -token "$API_KEY" -health > /dev/null 2>&1; then
    print_success "Health check OK avec token"
else
    print_error "Health check échoué avec token"
fi

print_test "Test exécution avec token..."
if $CLIENT -server "$SERVER_URL" -token "$API_KEY" "$TEST_FILE" > /dev/null 2>&1; then
    print_success "Exécution OK avec token"
else
    print_error "Exécution échouée avec token"
fi

print_test "Test avec mauvais token (doit échouer)..."
if $CLIENT -server "$SERVER_URL" -token "bad-token-12345" -health > /dev/null 2>&1; then
    print_error "Devrait échouer avec un mauvais token"
else
    print_success "Rejet du mauvais token comme attendu"
fi

print_test "Test via variable d'environnement..."
export TSD_AUTH_TOKEN="$API_KEY"
if $CLIENT -server "$SERVER_URL" -health > /dev/null 2>&1; then
    print_success "Authentification via variable d'environnement OK"
else
    print_error "Authentification via variable d'environnement échouée"
fi
unset TSD_AUTH_TOKEN

print_test "Arrêt du serveur..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
SERVER_PID=""
sleep 1
print_success "Serveur arrêté"

# Test 3: JWT
print_header "Test 3: Authentification JWT"

print_test "Génération d'un secret JWT..."
JWT_SECRET=$($AUTH_TOOL generate-key -format json 2>/dev/null | grep -o '"keys":\["[^"]*"' | cut -d'"' -f4)
if [ -z "$JWT_SECRET" ]; then
    print_error "Génération de secret JWT échouée"
    exit 1
fi
print_success "Secret JWT généré: ${JWT_SECRET:0:20}..."

print_test "Génération d'un JWT pour 'alice'..."
JWT_TOKEN=$($AUTH_TOOL generate-jwt \
    -secret "$JWT_SECRET" \
    -username "alice" \
    -roles "admin,user" \
    -expiration 1h \
    -format json 2>/dev/null | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$JWT_TOKEN" ]; then
    print_error "Génération de JWT échouée"
    exit 1
fi
print_success "JWT généré: ${JWT_TOKEN:0:30}..."

print_test "Validation du JWT..."
if $AUTH_TOOL validate -type jwt -token "$JWT_TOKEN" -secret "$JWT_SECRET" > /dev/null 2>&1; then
    print_success "Validation JWT OK"
else
    print_error "Validation JWT échouée"
fi

print_test "Démarrage du serveur avec JWT..."
export TSD_JWT_SECRET="$JWT_SECRET"
$SERVER -port $SERVER_PORT -auth jwt -jwt-expiration 1h > /tmp/tsd_server_jwt.log 2>&1 &
SERVER_PID=$!
sleep 2

if ! kill -0 $SERVER_PID 2>/dev/null; then
    print_error "Le serveur n'a pas démarré"
    cat /tmp/tsd_server_jwt.log
    exit 1
fi
print_success "Serveur démarré avec JWT"

print_test "Test health check SANS JWT (doit échouer)..."
if $CLIENT -server "$SERVER_URL" -health > /dev/null 2>&1; then
    print_error "Health check devrait échouer sans JWT"
else
    print_success "Health check échoué comme attendu"
fi

print_test "Test health check AVEC JWT..."
if $CLIENT -server "$SERVER_URL" -token "$JWT_TOKEN" -health > /dev/null 2>&1; then
    print_success "Health check OK avec JWT"
else
    print_error "Health check échoué avec JWT"
fi

print_test "Test exécution avec JWT..."
if $CLIENT -server "$SERVER_URL" -token "$JWT_TOKEN" "$TEST_FILE" > /dev/null 2>&1; then
    print_success "Exécution OK avec JWT"
else
    print_error "Exécution échouée avec JWT"
fi

print_test "Test avec JWT invalide (doit échouer)..."
if $CLIENT -server "$SERVER_URL" -token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid" -health > /dev/null 2>&1; then
    print_error "Devrait échouer avec un JWT invalide"
else
    print_success "Rejet du JWT invalide comme attendu"
fi

print_test "Génération d'un JWT pour 'bob'..."
JWT_TOKEN_BOB=$($AUTH_TOOL generate-jwt \
    -secret "$JWT_SECRET" \
    -username "bob" \
    -roles "developer" \
    -expiration 30m \
    -format json 2>/dev/null | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$JWT_TOKEN_BOB" ]; then
    print_success "JWT pour 'bob' généré"

    print_test "Test avec JWT de Bob..."
    if $CLIENT -server "$SERVER_URL" -token "$JWT_TOKEN_BOB" -health > /dev/null 2>&1; then
        print_success "Authentification OK avec JWT de Bob"
    else
        print_error "Authentification échouée avec JWT de Bob"
    fi
else
    print_error "Génération de JWT pour 'bob' échouée"
fi

print_test "Arrêt du serveur..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
SERVER_PID=""
sleep 1
print_success "Serveur arrêté"

# Test 4: Scénarios d'erreur
print_header "Test 4: Scénarios d'erreur"

print_test "Test clé API trop courte..."
if $AUTH_TOOL validate -type key -token "short" -keys "short" > /dev/null 2>&1; then
    print_error "Devrait rejeter les clés courtes"
else
    print_success "Rejet des clés courtes OK"
fi

print_test "Test secret JWT trop court..."
export TSD_JWT_SECRET="court"
if $SERVER -port $SERVER_PORT -auth jwt > /dev/null 2>&1 &; then
    SERVER_PID=$!
    sleep 1
    if kill -0 $SERVER_PID 2>/dev/null; then
        print_error "Le serveur ne devrait pas démarrer avec un secret court"
        kill $SERVER_PID 2>/dev/null || true
    else
        print_success "Rejet du secret court OK"
    fi
    SERVER_PID=""
else
    print_success "Rejet du secret court OK"
fi

print_test "Test serveur sans clés configurées..."
unset TSD_AUTH_KEYS
if $SERVER -port $SERVER_PORT -auth key > /dev/null 2>&1 &; then
    SERVER_PID=$!
    sleep 1
    if kill -0 $SERVER_PID 2>/dev/null; then
        print_error "Le serveur ne devrait pas démarrer sans clés"
        kill $SERVER_PID 2>/dev/null || true
    else
        print_success "Rejet de la config sans clés OK"
    fi
    SERVER_PID=""
else
    print_success "Rejet de la config sans clés OK"
fi

# Test 5: Curl
print_header "Test 5: Test avec curl"

if command -v curl &> /dev/null; then
    print_test "Démarrage du serveur pour tests curl..."
    export TSD_AUTH_KEYS="$API_KEY"
    $SERVER -port $SERVER_PORT -auth key > /tmp/tsd_server_curl.log 2>&1 &
    SERVER_PID=$!
    sleep 2

    print_test "Test curl health check avec token..."
    if curl -s -H "Authorization: Bearer $API_KEY" "$SERVER_URL/health" | grep -q "ok"; then
        print_success "Curl health check OK"
    else
        print_error "Curl health check échoué"
    fi

    print_test "Test curl execute avec token..."
    CURL_RESULT=$(curl -s -H "Authorization: Bearer $API_KEY" \
        -H "Content-Type: application/json" \
        -d '{"source":"type Test : <id: string>\nTest(\"t1\")","source_name":"curl_test"}' \
        "$SERVER_URL/api/v1/execute")

    if echo "$CURL_RESULT" | grep -q '"success":true'; then
        print_success "Curl execute OK"
    else
        print_error "Curl execute échoué"
    fi

    print_test "Arrêt du serveur..."
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
    SERVER_PID=""
else
    print_info "curl non installé, tests curl ignorés"
fi

# Résumé
print_header "Résumé des tests"

TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))

echo ""
echo -e "${BLUE}Tests exécutés: $TOTAL_TESTS${NC}"
echo -e "${GREEN}Tests réussis:  $TESTS_PASSED${NC}"
echo -e "${RED}Tests échoués:  $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ Tous les tests sont passés!${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}❌ Certains tests ont échoué.${NC}"
    echo ""
    exit 1
fi
