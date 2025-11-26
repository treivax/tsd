#!/bin/bash

# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script pour mettre à jour les règles dans les fichiers de test Go
# Ce script transforme les règles dans les chaînes de caractères Go de format:
#   {var: Type} / condition ==> action
# en format:
#   rule r1 : {var: Type} / condition ==> action

set -e

echo "🔄 Mise à jour des règles dans les fichiers de test Go"
echo "========================================================"

# Compteurs
total_files=0
total_updated=0

# Trouver tous les fichiers *_test.go
while IFS= read -r file; do
    if [ ! -f "$file" ]; then
        continue
    fi

    total_files=$((total_files + 1))

    # Créer un fichier temporaire
    tmp_file="${file}.tmp"
    file_updated=false

    # Utiliser awk pour traiter le fichier ligne par ligne avec contexte multi-ligne
    awk '
    BEGIN {
        in_string = 0
        rule_counter = 0
    }
    {
        line = $0

        # Détecter les backticks de chaînes multi-lignes
        if (line ~ /`/) {
            if (in_string == 0) {
                in_string = 1
                rule_counter = 0
            } else {
                in_string = 0
            }
        }

        # Si on est dans une chaîne et que la ligne contient une règle non migrée
        if (in_string == 1 && line ~ /^[[:space:]]*\{[a-z]+:[[:space:]]*[A-Z]/ && line ~ /\/.*==>/ && line !~ /^[[:space:]]*rule[[:space:]]/) {
            rule_counter++
            # Extraire l'\''indentation
            match(line, /^[[:space:]]*/)
            indent = substr(line, RSTART, RLENGTH)
            # Extraire la règle (sans l'\''indentation)
            rule = substr(line, RLENGTH + 1)
            # Ajouter l'\''identifiant
            print indent "rule r" rule_counter " : " rule
            updated = 1
        } else {
            print line
        }
    }
    END {
        if (updated == 1) {
            exit 0
        } else {
            exit 1
        }
    }
    ' "$file" > "$tmp_file"

    # Vérifier si awk a détecté des changements (exit code 0 = changements, 1 = pas de changements)
    if [ $? -eq 0 ]; then
        mv "$tmp_file" "$file"
        total_updated=$((total_updated + 1))
        echo "  ✅ Mis à jour: $file"
        file_updated=true
    else
        rm -f "$tmp_file"
    fi

done < <(find . -name "*_test.go" -type f | sort)

echo ""
echo "========================================================"
echo "✨ Migration des tests Go terminée !"
echo "📊 Statistiques:"
echo "   - Fichiers traités: $total_files"
echo "   - Fichiers mis à jour: $total_updated"
echo "========================================================"
