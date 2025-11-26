# 🧪 Lancer l'ensemble des tests TSD

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Le projet contient plusieurs types de tests :
- Tests unitaires Go (modules `rete/`, `constraint/`, `test/`)
- Tests d'intégration (`test/integration/`)
- Runner universel RETE (58 tests Beta + Alpha + Intégration)

## Objectif

Exécuter tous les tests du projet pour valider que le système fonctionne correctement.

## Instructions

1. **Lancer les tests unitaires Go** :
   - Exécuter `make test`
   - Vérifier que tous les modules passent (rete, constraint, test, integration)

2. **Lancer le runner universel RETE** :
   - Exécuter `make rete-unified`
   - Vérifier que les 58 tests passent

3. **Vérifier l'absence d'erreurs critiques** :
   - Pas d'erreur "variable non liée"
   - Pas d'erreur de parsing
   - Pas d'erreur de réseau RETE

4. **Générer un rapport de synthèse** :
   - Nombre de tests passés/échoués
   - Temps d'exécution
   - Modules testés
   - Erreurs éventuelles

## Critères de Succès

✅ Tous les tests unitaires Go passent (PASS)
✅ Les 58 tests du runner universel passent
✅ Aucune erreur critique détectée
✅ Rapport de synthèse généré

## Commandes Make Disponibles

```bash
make test                 # Tests unitaires Go
make test-coverage        # Tests avec couverture
make test-integration     # Tests d'intégration uniquement
make rete-unified         # Runner universel (tous les tests RETE)
make rete-all             # Tous les tests beta RETE
make validate             # Validation complète (format + lint + build + test)
```

## Format de Réponse Attendu

```
=== RÉCAPITULATIF DES TESTS ===

1. Tests Unitaires Go : [STATUT]
   - constraint : [OK/FAIL]
   - rete : [OK/FAIL]
   - test : [OK/FAIL]
   - integration : [OK/FAIL]

2. Runner Universel RETE : [STATUT]
   - Tests exécutés : X
   - Tests réussis : X
   - Tests échoués : X

3. Erreurs Critiques : [OUI/NON]
   - Détails des erreurs le cas échéant

4. Conclusion : [SUCCÈS/ÉCHEC]
```

## Exemple d'Utilisation

```
Relance moi l'ensemble des tests, dont le runner universel
```

ou plus simplement :

```
Utilise le prompt "run-tests"
```
