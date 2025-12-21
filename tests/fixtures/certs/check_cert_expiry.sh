#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script pour vérifier l'expiration des certificats de test TLS
# Utilisé en CI/CD pour s'assurer que les certificats sont toujours valides

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CERT_FILE="$SCRIPT_DIR/test-server.crt"
WARNING_DAYS=30  # Avertir si expiration < 30 jours
ERROR_DAYS=7     # Erreur si expiration < 7 jours

echo "🔍 Vérification de l'expiration des certificats de test..."
echo "📁 Répertoire: $SCRIPT_DIR"
echo ""

# Vérifier que le certificat existe
if [ ! -f "$CERT_FILE" ]; then
    echo "⚠️  Certificat non trouvé: $CERT_FILE"
    echo "💡 Génération automatique..."
    bash "$SCRIPT_DIR/generate_certs.sh"
    exit 0
fi

# Extraire la date d'expiration
EXPIRY_DATE=$(openssl x509 -enddate -noout -in "$CERT_FILE" | cut -d= -f2)
EXPIRY_EPOCH=$(date -d "$EXPIRY_DATE" +%s 2>/dev/null || date -j -f "%b %d %T %Y %Z" "$EXPIRY_DATE" +%s 2>/dev/null)
CURRENT_EPOCH=$(date +%s)
DAYS_LEFT=$(( ($EXPIRY_EPOCH - $CURRENT_EPOCH) / 86400 ))

echo "📅 Date d'expiration: $EXPIRY_DATE"
echo "⏳ Jours restants: $DAYS_LEFT jours"
echo ""

# Vérifier le statut
if [ $DAYS_LEFT -lt 0 ]; then
    echo "❌ ERREUR: Le certificat a expiré il y a $((-$DAYS_LEFT)) jours!"
    echo "💡 Régénération automatique..."
    bash "$SCRIPT_DIR/generate_certs.sh"
    exit 0
elif [ $DAYS_LEFT -lt $ERROR_DAYS ]; then
    echo "❌ ERREUR: Le certificat expire dans $DAYS_LEFT jours (seuil: $ERROR_DAYS jours)"
    echo "💡 Régénération automatique..."
    bash "$SCRIPT_DIR/generate_certs.sh"
    exit 0
elif [ $DAYS_LEFT -lt $WARNING_DAYS ]; then
    echo "⚠️  AVERTISSEMENT: Le certificat expire dans $DAYS_LEFT jours (seuil: $WARNING_DAYS jours)"
    echo "💡 Considérez une régénération prochainement"
    echo ""
    echo "Pour régénérer:"
    echo "  cd $SCRIPT_DIR"
    echo "  ./generate_certs.sh"
    exit 0
else
    echo "✅ Certificat valide pour encore $DAYS_LEFT jours"
fi

# Vérifier également la clé privée
KEY_FILE="$SCRIPT_DIR/test-server.key"
if [ ! -f "$KEY_FILE" ]; then
    echo ""
    echo "⚠️  Clé privée manquante: $KEY_FILE"
    echo "💡 Régénération complète..."
    bash "$SCRIPT_DIR/generate_certs.sh"
    exit 0
fi

echo "✅ Clé privée présente"
echo ""
echo "🎯 Vérification terminée avec succès"
exit 0
