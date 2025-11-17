#!/usr/bin/env python3
import os
import re
import glob

def fix_facts_file(filepath):
    """Convertit le format {id: "value", field: value} vers Format[id=value, field=value]"""
    with open(filepath, 'r') as f:
        content = f.read()
    
    print(f"Traitement: {filepath}")
    print(f"Avant: {repr(content[:100])}")
    
    lines = content.strip().split('\n')
    fixed_lines = []
    
    for line in lines:
        if not line.strip():
            continue
            
        # Pattern pour Type{id: "value", ...}
        match = re.match(r'(\w+)\{id:\s*"([^"]*)",?\s*(.*)\}', line)
        if match:
            type_name = match.group(1)
            id_value = match.group(2)
            rest = match.group(3)
            
            # Nettoie les autres champs
            # Enlève quotes autour des valeurs string et convertit : vers =
            rest = re.sub(r'"([^"]*)"', r'\1', rest)  # Enlève quotes
            rest = re.sub(r':\s*', '=', rest)          # : vers =
            rest = rest.replace(', ', ', ')          # Nettoie espaces
            
            # Reconstructs au format TSD
            new_line = f"{type_name}[id={id_value}, {rest}]"
            fixed_lines.append(new_line)
            print(f"  Converti: {line} -> {new_line}")
        else:
            # Garde la ligne telle quelle si pas de match
            fixed_lines.append(line)
    
    new_content = '\n'.join(fixed_lines) + '\n'
    print(f"Après: {repr(new_content[:100])}")
    
    with open(filepath, 'w') as f:
        f.write(new_content)

# Traite tous les fichiers .facts dans le dossier extended
facts_files = glob.glob('/home/resinsec/dev/tsd/alpha_coverage_tests_extended/*.facts')

for facts_file in facts_files:
    fix_facts_file(facts_file)
    print("---")

print(f"Traitement terminé pour {len(facts_files)} fichiers")