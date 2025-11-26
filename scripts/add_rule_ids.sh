#!/bin/bash

# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script pour ajouter des identifiants de règles aux fichiers .constraint
# Ce script transforme les règles de format:
#   {var: Type} / condition ==> action
# en format:
#   rule r1 : {var: Type} / condition ==> action

set -e

echo "🔄 Ajout des identifiants de règles aux fichiers .constraint"
echo "============================================================"

# Compteurs
total_files=0
total_rules=0
total_updated=0

# Trouver tous les fichiers .constraint
while IFS= read -r file; do
    if [ ! -f "$file" ]; then
        continue
    fi

    total_files=$((total_files + 1))
    echo ""
    echo "📄 Traitement: $file"

    # Créer un fichier temporaire
    tmp_file="${file}.tmp"

    # Compteur de règles pour ce fichier
    rule_counter=0
    file_updated=false

    # Lire le fichier ligne par ligne
    while IFS= read -r line || [ -n "$line" ]; do
        # Détecter si la ligne contient une règle (commence par { et contient / et ==>)
        if [[ "$line" =~ ^\{.*\}.*/.*(==>|→) ]]; then
            rule_counter=$((rule_counter + 1))
            total_rules=$((total_rules + 1))

            # Générer l'identifiant de règle
            filename=$(basename "$file" .constraint)
            rule_id="r${rule_counter}"

            # Vérifier si la règle a déjà un identifiant (commence par "rule ")
            if [[ "$line" =~ ^rule[[:space:]] ]]; then
                echo "  ⏭️  Règle $rule_counter déjà migrée"
                echo "$line" >> "$tmp_file"
            else
                # Ajouter l'identifiant de règle
                updated_line="rule ${rule_id} : ${line}"
                echo "  ✅ Règle $rule_counter: ajout de l'identifiant '${rule_id}'"
                echo "$updated_line" >> "$tmp_file"
                file_updated=true
            fi
        else
            # Ligne normale (commentaire, type, fact, etc.)
            echo "$line" >> "$tmp_file"
        fi
    done < "$file"

    # Remplacer le fichier original si des modifications ont été faites
    if [ "$file_updated" = true ]; then
        mv "$tmp_file" "$file"
        total_updated=$((total_updated + 1))
        echo "  💾 Fichier mis à jour avec $rule_counter règles"
    else
        rm -f "$tmp_file"
        echo "  ℹ️  Aucune modification nécessaire"
    fi

done < <(find . -name "*.constraint" -type f | sort)

echo ""
echo "============================================================"
echo "✨ Migration terminée !"
echo "📊 Statistiques:"
echo "   - Fichiers traités: $total_files"
echo "   - Fichiers mis à jour: $total_updated"
echo "   - Règles totales: $total_rules"
echo "============================================================"
