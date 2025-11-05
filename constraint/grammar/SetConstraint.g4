grammar SetConstraint;

// Point d'entrée principal
program: (typeDefinition)* expression EOF ;

// Définition de type
typeDefinition: 'type' IDENTIFIER ':' '<' fieldList '>' ;

// Liste de champs séparés par des virgules
fieldList: field (',' field)* ;

// Un champ avec son nom et son type
field: IDENTIFIER ':' atomicType ;

// Types atomiques supportés
atomicType: 'string' | 'number' | 'bool' ;

// Expression complète : ensemble / contraintes
expression: set '/' constraints ;

// Définition d'un ensemble
set: '{' typedVariableList '}' ;

// Liste de variables typées séparées par des virgules
typedVariableList: typedVariable (',' typedVariable)* ;

// Une variable avec son type
typedVariable: IDENTIFIER ':' IDENTIFIER ;

// Une variable (identifiant) - pour les expressions arithmétiques
variable: IDENTIFIER ;

// Liste de contraintes séparées par des opérateurs logiques
constraints: constraint (logicalOp constraint)* ;

// Une contrainte individuelle
constraint: 
    | arithmeticExpr comparisonOp arithmeticExpr    // x + y = 34
    | '(' constraints ')'                           // Parenthèses pour grouper
    ;

// Expression arithmétique
arithmeticExpr: 
    | arithmeticExpr ('*' | '/') arithmeticExpr    // Multiplication, division
    | arithmeticExpr ('+' | '-') arithmeticExpr    // Addition, soustraction  
    | '(' arithmeticExpr ')'                       // Parenthèses
    | fieldAccess                                  // Accès aux champs d'un objet
    | variable                                     // Variable
    | NUMBER                                       // Nombre
    | STRING_LITERAL                               // Chaîne de caractères
    | booleanLiteral                               // Valeur booléenne
    ;

// Accès aux champs d'un objet
fieldAccess: variable '.' IDENTIFIER ;

// Valeurs booléennes
booleanLiteral: 'true' | 'false' ;

// Opérateurs de comparaison
comparisonOp: '=' | '!=' | '<' | '>' | '<=' | '>=' ;

// Opérateurs logiques
logicalOp: 'AND' | 'OR' | '&' | '|' | '&&' | '||' ;

// Tokens lexicaux
IDENTIFIER: [a-zA-Z][a-zA-Z0-9_]* ;
NUMBER: [0-9]+ ('.' [0-9]+)? ;
STRING_LITERAL: '"' (~["])* '"' | '\'' (~['])* '\'' ;

// Ignorer les espaces et commentaires
WS: [ \t\r\n]+ -> skip ;
COMMENT: '//' ~[\r\n]* -> skip ;
BLOCK_COMMENT: '/*' .*? '*/' -> skip ;