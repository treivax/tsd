#!/usr/bin/env python3
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

"""
Script pour ajouter automatiquement les définitions d'actions manquantes dans les tests Go
qui génèrent dynamiquement du contenu TSD.
"""

import re
import sys
from pathlib import Path


def extract_actions_from_rules(content):
    """Extrait tous les appels d'actions dans les règles"""
    # Pattern pour les appels d'action: ==> action_name(args)
    # Supporte les noms en minuscules et majuscules
    pattern = r"==>\s*([A-Za-z_]\w*)\s*\("
    actions = set()
    for match in re.finditer(pattern, content):
        actions.add(match.group(1))
    return actions


def has_action_definition(content, action_name):
    """Vérifie si une action est déjà définie"""
    pattern = rf"^\s*action\s+{re.escape(action_name)}\s*\("
    return re.search(pattern, content, re.MULTILINE) is not None


def add_missing_action_definitions(content):
    """Ajoute les définitions d'actions manquantes dans un contenu TSD"""
    # Extraire les actions appelées
    called_actions = extract_actions_from_rules(content)

    # Trouver les actions déjà définies
    defined_actions = set()
    for action in called_actions:
        if has_action_definition(content, action):
            defined_actions.add(action)

    # Actions manquantes
    missing_actions = called_actions - defined_actions

    if not missing_actions:
        return content, False

    # Générer les définitions manquantes
    definitions = []
    for action in sorted(missing_actions):
        # Actions courantes avec signatures typiques
        action_lower = action.lower()
        if action_lower == "print":
            definitions.append(
                f"action {action}(arg1: string, arg2: string, arg3: string)"
            )
        elif action_lower.endswith("_detected") or action_lower.endswith("_found"):
            definitions.append(f"action {action}(id: string)")
        elif action_lower.startswith("log"):
            definitions.append(f"action {action}(message: string)")
        elif action_lower.startswith("notify"):
            definitions.append(f"action {action}(message: string)")
        elif action_lower.startswith("process") or action_lower.startswith("handle"):
            definitions.append(f"action {action}(id: string)")
        else:
            # Par défaut, une action simple avec un string
            definitions.append(f"action {action}(arg: string)")

    # Trouver où insérer les définitions (après les types, avant les règles)
    rule_match = re.search(r"^(rule\s+\w+\s*:)", content, re.MULTILINE)

    if rule_match:
        # Insère avant la première règle
        insert_pos = rule_match.start()
        actions_block = "\n" + "\n".join(definitions) + "\n\n"
        new_content = content[:insert_pos] + actions_block + content[insert_pos:]
    else:
        # Pas de règle trouvée, ajoute à la fin des types
        type_matches = list(re.finditer(r"^type\s+\w+\(.*?\)", content, re.MULTILINE))
        if type_matches:
            last_type = type_matches[-1]
            insert_pos = last_type.end()
            actions_block = "\n\n" + "\n".join(definitions) + "\n"
            new_content = content[:insert_pos] + actions_block + content[insert_pos:]
        else:
            # Ajoute au début
            actions_block = "\n".join(definitions) + "\n\n"
            new_content = actions_block + content

    return new_content, True


def fix_test_file(file_path):
    """Corrige un fichier de test Go en ajoutant les actions manquantes dans les content strings"""
    try:
        with open(file_path, "r", encoding="utf-8") as f:
            original_content = f.read()

        # Pattern pour trouver les content := `...` blocks
        pattern = r"(content\s*:=\s*`)(.*?)(`)"

        modified = False
        new_content = original_content

        def replace_content(match):
            nonlocal modified
            prefix = match.group(1)
            tsd_content = match.group(2)
            suffix = match.group(3)

            # Ajouter les actions manquantes
            new_tsd_content, changed = add_missing_action_definitions(tsd_content)

            if changed:
                modified = True
                return prefix + new_tsd_content + suffix
            return match.group(0)

        new_content = re.sub(
            pattern, replace_content, original_content, flags=re.DOTALL
        )

        if modified:
            # Créer un backup
            backup_path = str(file_path) + ".bak"
            with open(backup_path, "w", encoding="utf-8") as f:
                f.write(original_content)

            # Écrire le nouveau contenu
            with open(file_path, "w", encoding="utf-8") as f:
                f.write(new_content)

            return True, None

        return False, None

    except Exception as e:
        return False, str(e)


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 fix_test_actions.py <directory>")
        sys.exit(1)

    directory = Path(sys.argv[1])

    if not directory.exists():
        print(f"Erreur: Le répertoire '{directory}' n'existe pas")
        sys.exit(1)

    print(f"🔍 Recherche de fichiers *_test.go dans {directory}...")
    print()

    modified_files = []
    error_files = []
    unchanged_files = []

    for test_file in directory.rglob("*_test.go"):
        modified, error = fix_test_file(test_file)

        if error:
            error_files.append((str(test_file), error))
            print(f"✗ Erreur: {test_file} - {error}")
        elif modified:
            modified_files.append(str(test_file))
            print(f"✓ Modifié: {test_file}")
        else:
            unchanged_files.append(str(test_file))
            print(f"○ Inchangé: {test_file}")

    print()
    print("=" * 60)
    print("Résumé:")
    print(f"  Fichiers modifiés: {len(modified_files)}")
    print(f"  Fichiers inchangés: {len(unchanged_files)}")
    print(f"  Erreurs: {len(error_files)}")

    if error_files:
        print()
        print("Fichiers avec erreurs:")
        for file, error in error_files:
            print(f"  - {file}: {error}")

    if modified_files:
        print()
        print("✅ Les actions manquantes ont été ajoutées dans les tests!")
        print("⚠️  Des backups ont été créés (.bak)")


if __name__ == "__main__":
    main()
