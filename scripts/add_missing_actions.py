#!/usr/bin/env python3
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

"""
Script pour ajouter automatiquement les définitions d'actions manquantes dans les fichiers .tsd
"""

import re
import sys
from collections import defaultdict
from pathlib import Path


def extract_action_calls(content):
    """Extrait tous les appels d'actions dans le contenu"""
    # Pattern pour les appels d'action: ==> action_name(args)
    pattern = r"==>\s*(\w+)\s*\("
    actions = set()
    for match in re.finditer(pattern, content):
        actions.add(match.group(1))
    return actions


def extract_action_definitions(content):
    """Extrait toutes les définitions d'actions existantes"""
    pattern = r"^action\s+(\w+)\s*\("
    actions = set()
    for match in re.finditer(pattern, content, re.MULTILINE):
        actions.add(match.group(1))
    return actions


def extract_action_signature(content, action_name):
    """Extrait la signature d'une action à partir de ses appels"""
    # Trouve tous les appels de cette action
    pattern = rf"==>\s*{re.escape(action_name)}\s*\((.*?)\)"
    calls = re.findall(pattern, content)

    if not calls:
        return f"action {action_name}()"

    # Prend le premier appel pour déterminer le nombre d'arguments
    first_call = calls[0]
    if not first_call.strip():
        return f"action {action_name}()"

    # Compte le nombre d'arguments (simplifié)
    args = [arg.strip() for arg in first_call.split(",")]

    # Génère des paramètres génériques
    params = []
    for i, arg in enumerate(args):
        # Devine le type basé sur le contenu de l'argument
        if arg.startswith('"') or arg.startswith("'"):
            param_type = "string"
        elif arg.replace(".", "", 1).isdigit() or arg.isdigit():
            param_type = "number"
        elif arg in ["true", "false"]:
            param_type = "bool"
        elif "." in arg:
            # Accès à un champ - utilise string par défaut
            param_type = "string"
        else:
            # Variable - utilise string par défaut
            param_type = "string"

        params.append(f"arg{i + 1}: {param_type}")

    return f"action {action_name}({', '.join(params)})"


def add_missing_actions(file_path):
    """Ajoute les définitions d'actions manquantes dans un fichier"""
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Extrait les actions appelées et définies
    called_actions = extract_action_calls(content)
    defined_actions = extract_action_definitions(content)

    # Trouve les actions manquantes
    missing_actions = called_actions - defined_actions

    if not missing_actions:
        return False, content

    # Génère les définitions manquantes
    definitions = []
    for action in sorted(missing_actions):
        signature = extract_action_signature(content, action)
        definitions.append(signature)

    # Trouve où insérer les définitions (après les types, avant les règles)
    # Cherche la première règle
    rule_match = re.search(r"^(rule\s+\w+\s*:)", content, re.MULTILINE)

    if rule_match:
        # Insère avant la première règle
        insert_pos = rule_match.start()

        # Ajoute les définitions avec commentaire
        actions_block = "\n// Actions auto-générées\n"
        actions_block += "\n".join(definitions)
        actions_block += "\n\n"

        new_content = content[:insert_pos] + actions_block + content[insert_pos:]
    else:
        # Pas de règle trouvée, ajoute à la fin
        new_content = content + "\n\n// Actions auto-générées\n"
        new_content += "\n".join(definitions)
        new_content += "\n"

    return True, new_content


def process_directory(directory):
    """Traite tous les fichiers .tsd dans un répertoire"""
    modified_files = []
    error_files = []

    for tsd_file in Path(directory).rglob("*.tsd"):
        try:
            modified, new_content = add_missing_actions(tsd_file)
            if modified:
                # Crée un backup
                backup_file = tsd_file.with_suffix(".tsd.bak")
                tsd_file.rename(backup_file)

                # Écrit le nouveau contenu
                with open(tsd_file, "w", encoding="utf-8") as f:
                    f.write(new_content)

                modified_files.append(str(tsd_file))
                print(f"✓ Modifié: {tsd_file}")
            else:
                print(f"○ Aucune action manquante: {tsd_file}")
        except Exception as e:
            error_files.append((str(tsd_file), str(e)))
            print(f"✗ Erreur: {tsd_file} - {e}")

    return modified_files, error_files


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 add_missing_actions.py <directory>")
        sys.exit(1)

    directory = sys.argv[1]

    if not Path(directory).exists():
        print(f"Erreur: Le répertoire '{directory}' n'existe pas")
        sys.exit(1)

    print(f"🔍 Recherche de fichiers .tsd dans {directory}...")
    print()

    modified_files, error_files = process_directory(directory)

    print()
    print("=" * 60)
    print("Résumé:")
    print(f"  Fichiers modifiés: {len(modified_files)}")
    print(f"  Erreurs: {len(error_files)}")

    if error_files:
        print()
        print("Fichiers avec erreurs:")
        for file, error in error_files:
            print(f"  - {file}: {error}")

    if modified_files:
        print()
        print("✅ Les actions manquantes ont été ajoutées!")
        print("⚠️  Des backups ont été créés (.tsd.bak)")


if __name__ == "__main__":
    main()
