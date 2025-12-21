#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script pour générer des certificats auto-signés pour les tests TLS
# Ce script génère une paire certificat/clé pour les tests uniquement

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CERT_FILE="$SCRIPT_DIR/test-server.crt"
KEY_FILE="$SCRIPT_DIR/test-server.key"

echo "🔐 Génération de certificats de test auto-signés..."
echo "📁 Répertoire: $SCRIPT_DIR"

# Supprimer les anciens certificats s'ils existent
rm -f "$CERT_FILE" "$KEY_FILE"

# Générer une clé privée RSA 2048 bits
openssl genrsa -out "$KEY_FILE" 2048 2>/dev/null

# Générer un certificat auto-signé valide 365 jours
openssl req -new -x509 -sha256 \
    -key "$KEY_FILE" \
    -out "$CERT_FILE" \
    -days 365 \
    -subj "/C=FR/ST=Test/L=Test/O=TSD Test/OU=Testing/CN=localhost" \
    2>/dev/null

# Vérifier que les fichiers ont été créés
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo "✅ Certificats générés avec succès:"
    echo "   📄 Certificat: $CERT_FILE"
    echo "   🔑 Clé privée: $KEY_FILE"
    echo ""
    echo "⚠️  ATTENTION: Ces certificats sont UNIQUEMENT pour les tests!"
    echo "   Ne JAMAIS utiliser en production."
    exit 0
else
    echo "❌ Erreur lors de la génération des certificats"
    exit 1
fi
