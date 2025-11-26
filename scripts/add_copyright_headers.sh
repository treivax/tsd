#!/bin/bash

# Script pour ajouter les en-têtes de copyright dans tous les fichiers .go
# Usage: ./scripts/add_copyright_headers.sh

set -e

# Couleurs pour l'output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# En-tête de copyright à ajouter
read -r -d '' COPYRIGHT_HEADER << 'EOF' || true
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

EOF

# Compteurs
COUNT_ADDED=0
COUNT_SKIPPED=0
COUNT_ALREADY=0

echo "🔍 Recherche de tous les fichiers .go..."

# Trouver tous les fichiers .go sauf parser.go (code généré)
while IFS= read -r file; do
    # Vérifier si le fichier contient déjà l'en-tête de copyright
    if head -3 "$file" | grep -q "Copyright (c) 2025 TSD Contributors"; then
        echo "  ⏭️  Déjà présent: $file"
        ((COUNT_ALREADY++))
        continue
    fi

    # Vérifier si c'est du code généré
    if head -1 "$file" | grep -q "Code generated"; then
        echo "  🔧 Code généré (ignoré): $file"
        ((COUNT_SKIPPED++))
        continue
    fi

    # Créer un fichier temporaire avec l'en-tête
    {
        echo "$COPYRIGHT_HEADER"
        cat "$file"
    } > "${file}.tmp"

    # Remplacer le fichier original
    mv "${file}.tmp" "$file"

    echo -e "  ${GREEN}✅${NC} Ajouté: $file"
    ((COUNT_ADDED++))

done < <(find . -name "*.go" -type f ! -path "./.git/*" ! -path "./vendor/*")

echo ""
echo "================================"
echo -e "${GREEN}✅ Opération terminée!${NC}"
echo "================================"
echo "  Fichiers modifiés: $COUNT_ADDED"
echo "  Déjà présents: $COUNT_ALREADY"
echo "  Ignorés (code généré): $COUNT_SKIPPED"
echo "  Total traité: $((COUNT_ADDED + COUNT_ALREADY + COUNT_SKIPPED))"
echo ""
