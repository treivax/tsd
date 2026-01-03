# Inventaire des Tests - Migration IDs

Date: 2025-12-19

## Résumé Exécutif

Analyse du projet TSD pour la migration des tests unitaires vers la nouvelle gestion des identifiants internes (`_id_`).

## État Actuel

### Constantes et Structures

✅ **Déjà implémenté** :
- `FieldNameInternalID = "_id_"` dans `constraint/constraint_constants.go`
- Structure `Fact` dans RETE utilise `json:"_id_"` 
- `FieldResolver` interdit l'accès à `_id_`
- `FactContext` pour résolution de variables implémenté
- `GenerateFactID` supporte les références de faits

### Tests Existants

Total de fichiers de tests trouvés :

### Tests par Module

#### constraint/
61
fichiers de tests
451
fonctions de test

#### rete/
145
fichiers de tests
1258
fonctions de test

#### api/
5
fichiers de tests
30
fonctions de test

#### tsdio/
2
fichiers de tests
42
fonctions de test
