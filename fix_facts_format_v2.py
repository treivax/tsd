#!/usr/bin/env python3
import os
import re
import glob

def fix_facts_file(filepath):
    """Corrige le format des fichiers .facts pour TSD"""
    with open(filepath, 'r') as f:
        content = f.read()
    
    print(f"Traitement: {os.path.basename(filepath)}")
    
    lines = content.strip().split('\n')
    fixed_lines = []
    
    for line in lines:
        if not line.strip():
            continue
            
        # Pattern pour Type[XXXX, field=value=field2=value2]
        match = re.match(r'(\w+)\[([^,\[\]]+),\s*(.*)\]', line)
        if match:
            type_name = match.group(1)
            first_value = match.group(2).strip()
            rest = match.group(3).strip()
            
            # Si le premier champ n'a pas de "=" c'est probablement l'id
            if '=' not in first_value:
                id_field = f"id={first_value}"
            else:
                id_field = first_value
            
            # Corrige les champs restants qui ont des = multiples
            if rest:
                # Divise par = et regroupe en pairs field=value
                parts = rest.split('=')
                corrected_fields = []
                
                i = 0
                while i < len(parts):
                    if i + 1 < len(parts):
                        field = parts[i].strip()
                        value = parts[i + 1].strip()
                        # Si value contient encore un field name, le divise
                        if i + 2 < len(parts) and not value.replace('.', '').replace('-', '').isdigit() and not value.lower() in ['true', 'false']:
                            # Trouve où commence le prochain field
                            next_field_pos = value.rfind(' ')
                            if next_field_pos > 0:
                                actual_value = value[:next_field_pos]
                                next_field = value[next_field_pos+1:]
                                corrected_fields.append(f"{field}={actual_value}")
                                parts[i + 1] = next_field  # Prépare pour le prochain
                                i += 1
                            else:
                                corrected_fields.append(f"{field}={value}")
                                i += 2
                        else:
                            corrected_fields.append(f"{field}={value}")
                            i += 2
                    else:
                        break
                
                rest_corrected = ', '.join(corrected_fields)
            else:
                rest_corrected = ""
            
            if rest_corrected:
                new_line = f"{type_name}[{id_field}, {rest_corrected}]"
            else:
                new_line = f"{type_name}[{id_field}]"
            
            fixed_lines.append(new_line)
            print(f"  {line} -> {new_line}")
        else:
            fixed_lines.append(line)
    
    new_content = '\n'.join(fixed_lines) + '\n'
    
    with open(filepath, 'w') as f:
        f.write(new_content)

# Traite tous les fichiers .facts dans le dossier extended
facts_files = glob.glob('/home/resinsec/dev/tsd/alpha_coverage_tests_extended/*.facts')

for facts_file in facts_files:
    fix_facts_file(facts_file)
    print("---")

print(f"Correction terminée pour {len(facts_files)} fichiers")