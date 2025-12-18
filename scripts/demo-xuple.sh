#!/bin/bash

# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script de démonstration de l'action Xuple
# Ce script exécute le test d'intégration qui montre comment
# l'action Xuple crée des xuples dans des xuple-spaces

set -e

# Couleurs pour la sortie
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo ""
echo "╔══════════════════════════════════════════════════════════╗"
echo "║       TSD - Démonstration Action Xuple                  ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""

echo -e "${BLUE}📚 Cette démonstration montre:${NC}"
echo "   1. Création de xuple-spaces avec différentes politiques"
echo "   2. Création de xuples via l'action Xuple"
echo "   3. Inspection du contenu des xuple-spaces"
echo "   4. Test des politiques FIFO/LIFO et once/per-agent"
echo ""

echo -e "${YELLOW}📂 Fichiers de référence:${NC}"
echo "   - Exemple TSD: examples/xuples/xuple-action-example.tsd"
echo "   - Code test:   rete/actions/builtin_integration_test.go"
echo "   - Guide:       docs/ACTION_XUPLE_GUIDE.md"
echo ""

echo -e "${BLUE}🚀 Lancement du test d'intégration...${NC}"
echo ""

# Exécuter le test d'intégration
cd "$(dirname "$0")/.." || exit 1
go test -v ./rete/actions -run TestBuiltinActions_EndToEnd_XupleAction

echo ""
echo -e "${GREEN}✨ Démonstration terminée!${NC}"
echo ""
echo -e "${YELLOW}💡 Pour voir le code source de l'exemple TSD:${NC}"
echo "   cat examples/xuples/xuple-action-example.tsd"
echo ""
echo -e "${YELLOW}💡 Pour voir le code du test:${NC}"
echo "   cat rete/actions/builtin_integration_test.go"
echo ""
echo -e "${YELLOW}💡 Pour lire le guide complet:${NC}"
echo "   cat docs/ACTION_XUPLE_GUIDE.md"
echo ""
