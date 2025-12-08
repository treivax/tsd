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

2. 🏁 **Lancer les tests avec race detector (OBLIGATOIRE)** :
   - Exécuter `make test-race` ou `go test -race ./...`
   - ⚠️ **CRITIQUE** : Ce test est OBLIGATOIRE pour détecter les race conditions
   - Les race conditions ne sont détectées QUE avec le flag `-race`
   - Vérifier qu'aucune race condition n'est détectée
   - **Ne JAMAIS skip cette étape**, même si plus lente (~10x)

3. **Lancer le runner universel RETE** :
   - Exécuter `make rete-unified`
   - Vérifier que les 58 tests passent

4. **Vérifier l'absence d'erreurs critiques** :
   - Pas d'erreur "variable non liée"
   - Pas d'erreur de parsing
   - Pas d'erreur de réseau RETE
   - Pas de race condition détectée

5. **Générer un rapport de synthèse** :
   - Nombre de tests passés/échoués
   - Temps d'exécution
   - Modules testés
   - Race conditions détectées (doit être 0)
   - Erreurs éventuelles

## Critères de Succès

✅ Tous les tests unitaires Go passent (PASS)
🏁 **✅ `go test -race ./...` passe sans race condition (OBLIGATOIRE)**
✅ Les 58 tests du runner universel passent
✅ Aucune erreur critique détectée
✅ Rapport de synthèse généré

## Commandes Make Disponibles

```bash
make test                 # Tests unitaires Go
make test-race            # 🏁 Tests avec race detector (OBLIGATOIRE)
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

2. 🏁 Tests Race Detector (OBLIGATOIRE) : [STATUT]
   - Commande : go test -race ./...
   - Race conditions détectées : [OUI/NON]
   - Détails si race détectée

3. Runner Universel RETE : [STATUT]
   - Tests exécutés : X
   - Tests réussis : X
   - Tests échoués : X

4. Erreurs Critiques : [OUI/NON]
   - Détails des erreurs le cas échéant

5. Conclusion : [SUCCÈS/ÉCHEC]
   - ⚠️ Note : Échec si race conditions détectées
```

## Exemple d'Utilisation

```
Relance moi l'ensemble des tests, dont le runner universel
```

ou plus simplement :

```
Utilise le prompt "run-tests"
```
