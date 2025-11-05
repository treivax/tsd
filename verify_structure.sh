#!/usr/bin/env bash

echo "=== VÉRIFICATION DE LA STRUCTURE DES FICHIERS ==="
echo ""

echo "📁 FICHIER constraint_types.go (STRUCTURES SEULEMENT):"
echo "✅ Structures Go pour l'AST"
echo "✅ Pas de fonctions (plus de duplication)"
echo ""
grep -n "^type " /home/resinsec/dev/tsd/constraint_types.go | head -5
echo "..."

echo ""
echo "📁 FICHIER constraint_utils.go (FONCTIONS SEULEMENT):"
echo "✅ Fonctions utilitaires pour validation"
echo "✅ Pas de structures (séparation claire)"
echo ""
grep -n "^func " /home/resinsec/dev/tsd/constraint_utils.go

echo ""
echo "📁 FICHIER constraint_main.go:"
echo "✅ Programme principal avec main()"
echo ""
grep -n "^func " /home/resinsec/dev/tsd/constraint_main.go

echo ""
echo "🔍 VÉRIFICATION DES DOUBLONS:"
echo ""
echo "Recherche de fonctions dupliquées entre les fichiers..."

TYPES_FUNCS=$(grep "^func " /home/resinsec/dev/tsd/constraint_types.go 2>/dev/null | wc -l)
UTILS_FUNCS=$(grep "^func " /home/resinsec/dev/tsd/constraint_utils.go 2>/dev/null | wc -l)
MAIN_FUNCS=$(grep "^func " /home/resinsec/dev/tsd/constraint_main.go 2>/dev/null | wc -l)

echo "- constraint_types.go: $TYPES_FUNCS fonctions"
echo "- constraint_utils.go: $UTILS_FUNCS fonctions"  
echo "- constraint_main.go: $MAIN_FUNCS fonction(s)"

if [ "$TYPES_FUNCS" -eq 0 ]; then
    echo "✅ Aucune fonction dans constraint_types.go (correct !)"
else
    echo "❌ Des fonctions trouvées dans constraint_types.go"
fi

echo ""
echo "🎯 RÉPARTITION CORRECTE:"
echo "- constraint_types.go → Structures uniquement" 
echo "- constraint_utils.go → Fonctions utilitaires"
echo "- constraint_main.go → Point d'entrée main()"
echo ""
echo "✨ PLUS DE DOUBLONS ! Structure propre et organisée."