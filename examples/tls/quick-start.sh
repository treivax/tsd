#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script de démarrage rapide pour tester TLS/HTTPS avec TSD
# Ce script démontre la configuration complète d'un serveur HTTPS et client

set -e

# Couleurs pour la sortie
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TSD_BIN="${TSD_BIN:-../../bin/tsd}"
CERTS_DIR="./test-certs"
TEST_FILE="test-program.tsd"
SERVER_PORT="${SERVER_PORT:-8080}"

echo -e "${BLUE}🔐 Script de Démarrage Rapide TLS/HTTPS pour TSD${NC}"
echo "=================================================="
echo ""

# Vérifier que le binaire tsd existe
if [ ! -f "$TSD_BIN" ]; then
    echo -e "${RED}❌ Binaire TSD non trouvé: $TSD_BIN${NC}"
    echo ""
    echo "Veuillez compiler TSD d'abord :"
    echo "  cd ../../"
    echo "  make build"
    exit 1
fi

echo -e "${GREEN}✅ Binaire TSD trouvé: $TSD_BIN${NC}"
echo ""

# Étape 1 : Générer les certificats
echo -e "${BLUE}📋 Étape 1/5 : Génération des certificats TLS${NC}"
echo "=============================================="

if [ -d "$CERTS_DIR" ]; then
    echo -e "${YELLOW}⚠️  Le répertoire $CERTS_DIR existe déjà${NC}"
    read -p "Voulez-vous le supprimer et régénérer les certificats ? (o/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Oo]$ ]]; then
        rm -rf "$CERTS_DIR"
        echo -e "${GREEN}✅ Répertoire supprimé${NC}"
    else
        echo -e "${YELLOW}⏭️  Utilisation des certificats existants${NC}"
    fi
fi

if [ ! -d "$CERTS_DIR" ]; then
    echo -e "${BLUE}🔑 Génération des certificats auto-signés...${NC}"
    $TSD_BIN auth generate-cert \
        -output-dir "$CERTS_DIR" \
        -hosts "localhost,127.0.0.1" \
        -valid-days 365 \
        -org "TSD Test"
    echo ""
fi

# Vérifier que les certificats sont générés
if [ ! -f "$CERTS_DIR/server.crt" ] || [ ! -f "$CERTS_DIR/server.key" ]; then
    echo -e "${RED}❌ Échec de la génération des certificats${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Certificats générés avec succès${NC}"
echo -e "   📄 Certificat serveur: $CERTS_DIR/server.crt"
echo -e "   🔑 Clé privée serveur: $CERTS_DIR/server.key"
echo -e "   📄 Certificat CA: $CERTS_DIR/ca.crt"
echo ""

# Étape 2 : Créer un programme TSD de test
echo -e "${BLUE}📋 Étape 2/5 : Création d'un programme TSD de test${NC}"
echo "=================================================="

cat > "$TEST_FILE" << 'EOF'
# Programme de test TSD - Vérification d'âge
type Person : <id: string, name: string, age: int>

# Faits de test
fact p1 : Person <id: "1", name: "Alice", age: 30>
fact p2 : Person <id: "2", name: "Bob", age: 25>
fact p3 : Person <id: "3", name: "Charlie", age: 17>

# Règles
rule check_adult : {p: Person} / p.age >= 18 ==> adult(p.id, p.name)
rule check_minor : {p: Person} / p.age < 18 ==> minor(p.id, p.name)
EOF

echo -e "${GREEN}✅ Programme de test créé: $TEST_FILE${NC}"
echo ""
echo "Contenu du programme :"
echo "----------------------"
cat "$TEST_FILE"
echo "----------------------"
echo ""

# Étape 3 : Démarrer le serveur HTTPS
echo -e "${BLUE}📋 Étape 3/5 : Démarrage du serveur HTTPS${NC}"
echo "=========================================="

echo -e "${BLUE}🚀 Démarrage du serveur TSD en HTTPS sur le port $SERVER_PORT...${NC}"
$TSD_BIN server \
    --port "$SERVER_PORT" \
    --tls-cert "$CERTS_DIR/server.crt" \
    --tls-key "$CERTS_DIR/server.key" \
    -v &

SERVER_PID=$!

# Attendre que le serveur démarre
echo -e "${YELLOW}⏳ Attente du démarrage du serveur...${NC}"
sleep 3

# Vérifier que le serveur tourne
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo -e "${RED}❌ Le serveur n'a pas pu démarrer${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Serveur démarré (PID: $SERVER_PID)${NC}"
echo ""

# Fonction de nettoyage
cleanup() {
    echo ""
    echo -e "${YELLOW}🧹 Nettoyage...${NC}"
    if kill -0 $SERVER_PID 2>/dev/null; then
        echo -e "${BLUE}   Arrêt du serveur (PID: $SERVER_PID)...${NC}"
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
    fi
    echo -e "${GREEN}✅ Nettoyage terminé${NC}"
}

trap cleanup EXIT INT TERM

# Étape 4 : Tester avec le client (mode insecure)
echo -e "${BLUE}📋 Étape 4/5 : Test avec le client (mode insecure)${NC}"
echo "==================================================="

echo -e "${BLUE}🔧 Exécution du client en mode insecure (certificats auto-signés)...${NC}"
echo ""

$TSD_BIN client "$TEST_FILE" \
    -server "https://localhost:$SERVER_PORT" \
    -insecure \
    -v

echo ""
echo -e "${GREEN}✅ Test en mode insecure réussi${NC}"
echo ""

# Étape 5 : Tester avec vérification du CA
echo -e "${BLUE}📋 Étape 5/5 : Test avec vérification du CA${NC}"
echo "==========================================="

echo -e "${BLUE}🔒 Exécution du client avec vérification du CA...${NC}"
echo ""

$TSD_BIN client "$TEST_FILE" \
    -server "https://localhost:$SERVER_PORT" \
    -tls-ca "$CERTS_DIR/ca.crt" \
    -v

echo ""
echo -e "${GREEN}✅ Test avec CA réussi${NC}"
echo ""

# Résumé
echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}🎉 Tests TLS/HTTPS terminés avec succès !${NC}"
echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo ""
echo "Résumé des tests effectués :"
echo ""
echo "  1. ✅ Génération de certificats auto-signés"
echo "  2. ✅ Création d'un programme TSD de test"
echo "  3. ✅ Démarrage d'un serveur HTTPS"
echo "  4. ✅ Connexion client en mode insecure"
echo "  5. ✅ Connexion client avec vérification CA"
echo ""
echo "Fichiers générés :"
echo "  📁 $CERTS_DIR/        (certificats TLS)"
echo "  📄 $TEST_FILE         (programme de test)"
echo ""
echo "Pour tester manuellement :"
echo ""
echo "  # Démarrer le serveur"
echo "  $TSD_BIN server --tls-cert $CERTS_DIR/server.crt --tls-key $CERTS_DIR/server.key"
echo ""
echo "  # Client avec mode insecure"
echo "  $TSD_BIN client $TEST_FILE -server https://localhost:$SERVER_PORT -insecure"
echo ""
echo "  # Client avec vérification CA"
echo "  $TSD_BIN client $TEST_FILE -server https://localhost:$SERVER_PORT -tls-ca $CERTS_DIR/ca.crt"
echo ""
echo "  # Health check"
echo "  $TSD_BIN client -health -server https://localhost:$SERVER_PORT -insecure"
echo ""
echo -e "${YELLOW}⚠️  Note: Les certificats générés sont auto-signés et pour développement uniquement${NC}"
echo -e "${YELLOW}   En production, utilisez des certificats signés par une CA reconnue (Let's Encrypt, etc.)${NC}"
echo ""

# Le serveur sera arrêté automatiquement par le trap cleanup
echo -e "${BLUE}Appuyez sur Entrée pour arrêter le serveur et quitter...${NC}"
read -r
