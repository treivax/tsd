#!/usr/bin/env python3

import os
import re

def fix_facts_format_extended():
    """Corrige le format des fichiers .facts dans les tests étendus"""
    
    # Répertoire des tests étendus
    extended_dir = "alpha_coverage_tests_extended"
    
    # Pattern pour détecter le mauvais format
    old_pattern = re.compile(r'(\w+)\{([^}]+)\}')
    
    # Fonction pour convertir une ligne
    def convert_line(line):
        match = old_pattern.match(line.strip())
        if not match:
            return line
            
        type_name = match.group(1)
        fields_str = match.group(2)
        
        # Parser les champs
        fields = []
        # Découper par virgules en gérant les guillemets
        current_field = ""
        in_quotes = False
        for char in fields_str + ",":
            if char == '"' and (not current_field or current_field[-1] != '\\'):
                in_quotes = not in_quotes
                current_field += char
            elif char == ',' and not in_quotes:
                if current_field.strip():
                    fields.append(current_field.strip())
                current_field = ""
            else:
                current_field += char
        
        # Convertir chaque champ du format "key: value" vers "key=value"
        converted_fields = []
        for field in fields:
            # Remplacer ": " par "="
            field = field.replace(': ', '=')
            converted_fields.append(field)
        
        # Reconstruire au format attendu
        return f"{type_name}[{', '.join(converted_fields)}]"
    
    # Parcourir tous les fichiers .facts
    files_corrected = 0
    
    for filename in os.listdir(extended_dir):
        if filename.endswith('.facts'):
            filepath = os.path.join(extended_dir, filename)
            
            # Lire le fichier
            with open(filepath, 'r') as f:
                lines = f.readlines()
            
            # Convertir les lignes
            new_lines = []
            changed = False
            for line in lines:
                new_line = convert_line(line)
                if new_line != line:
                    changed = True
                new_lines.append(new_line)
            
            # Réécrire le fichier si changé
            if changed:
                with open(filepath, 'w') as f:
                    f.writelines(new_lines)
                print(f"✅ Corrigé: {filename}")
                files_corrected += 1
    
    return files_corrected

if __name__ == "__main__":
    print("🔧 CORRECTION FORMAT FICHIERS .facts ÉTENDUS")
    print("==========================================")
    
    corrected = fix_facts_format_extended()
    
    print(f"\n📊 Résultats:")
    print(f"   • {corrected} fichiers corrigés")
    print(f"   • Format: Type{{key: value}} → Type[key=value]")
    print("\n✅ Correction terminée !")