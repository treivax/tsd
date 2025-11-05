#!/usr/bin/env bash

echo "=== MODIFICATIONS APPORTÉES ==="
echo ""

echo "✅ 1. GRAMMAIRE PEG (constraint.peg)"
echo "   - Modifiée pour supporter PLUSIEURS expressions"
echo "   - Règle Start: types + ExpressionList (au lieu d'une seule Expression)"
echo "   - Nouvelle règle ExpressionList pour gérer plusieurs expressions"
echo ""

echo "✅ 2. STRUCTURES GO (constraint_types.go)"
echo "   - Program.Expression → Program.Expressions []Expression"
echo "   - Support de multiples expressions dans l'AST"
echo ""

echo "✅ 3. FONCTIONS UTILITAIRES (constraint_utils.go)"
echo "   - ValidateTypes: valide toutes les expressions"
echo "   - ValidateFieldAccess: prend un index d'expression en paramètre"  
echo "   - ValidateProgram: affiche le nombre d'expressions"
echo ""

echo "✅ 4. NOUVEAUX FICHIERS DE TEST"
echo "   - test_multi_expressions.txt: exemple avec 3 expressions"
echo ""

echo "✅ 5. DOCUMENTATION MISE À JOUR (PARSER_README.md)"
echo "   - Exemples avec plusieurs expressions"
echo "   - Format du langage étendu"
echo ""

echo "=== FORMAT SUPPORTÉ MAINTENANT ==="
echo ""
echo "type Type1 : < champ1: type1, champ2: type2 >"
echo "type Type2 : < champA: typeA, champB: typeB >"
echo ""
echo "{ var1: Type1, var2: Type2 } / contraintes1"
echo ""
echo "{ var3: Type1 } / contraintes2"
echo ""  
echo "{ var4: Type2, var5: Type1 } / contraintes3"
echo ""

echo "=== POUR TESTER ==="
echo ""
echo "1. Générer le parser:"
echo "   ./build.sh"
echo ""
echo "2. Tester avec plusieurs expressions:"
echo "   go run parser.go constraint_types.go constraint_utils.go constraint_main.go test_multi_expressions.txt"
echo ""

echo "✨ Le parser supporte maintenant UN OU PLUSIEURS types suivis de UNE OU PLUSIEURS expressions !"