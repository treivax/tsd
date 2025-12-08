# 🐛 Débugger un Test qui Échoue

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Un test échoue et tu as besoin d'identifier la cause racine du problème et de le corriger.

## Objectif

Analyser un test qui échoue, identifier la cause du problème, et proposer/implémenter une correction.

## ⚠️ RÈGLES STRICTES - TESTS RETE

### 🚫 INTERDICTIONS ABSOLUES POUR TESTS RETE

1. **AUCUNE SIMULATION DE RÉSULTATS** :
   - ❌ Pas de résultats hardcodés ou simulés
   - ❌ Pas de mock des résultats du réseau RETE
   - ❌ Pas de calcul manuel des tokens attendus
   - ✅ **TOUJOURS** extraire les résultats du réseau RETE réel
   - ✅ **TOUJOURS** interroger les TerminalNodes
   - ✅ **TOUJOURS** inspecter les mémoires (Left/Right/Result)

2. **EXTRACTION OBLIGATOIRE DEPUIS LE RÉSEAU** :
   ```go
   // ✅ BON - Extraction depuis le réseau
   terminalCount := 0
   for _, terminal := range network.TerminalNodes {
       terminalCount += len(terminal.Memory.GetTokens())
   }
   
   // ✅ BON - Inspection des tokens réels
   for _, token := range terminal.Memory.GetTokens() {
       for varName, fact := range token.Bindings {
           // Vérifier les données réelles du réseau
       }
   }
   
   // ❌ MAUVAIS - Simulation
   expectedTokens := 5  // Calculé manuellement !
   ```

3. **VALIDATION AVEC DONNÉES RÉSEAU RÉELLES** :
   - ✅ Compter les tokens dans les TerminalNodes
   - ✅ Vérifier les bindings dans les tokens
   - ✅ Inspecter les mémoires des JoinNodes
   - ✅ Tracer la propagation réelle
   - ❌ Ne jamais supposer le nombre de tokens
   - ❌ Ne jamais simuler les résultats

### ✅ BONNES PRATIQUES OBLIGATOIRES

1. **Code Golang** (si correction nécessaire) :
   - ❌ Aucun hardcoding de valeurs
   - ✅ Code générique avec paramètres
   - ✅ Constantes nommées pour toutes les valeurs
   - ✅ Respect des conventions Go (Effective Go)
   - ✅ go vet et golangci-lint sans erreur

2. **Tests** :
   - ✅ Extraction réelle depuis le réseau RETE
   - ✅ Validation des structures de données réelles
   - ✅ Messages d'assertion explicites
   - ✅ Tests déterministes et isolés

**Exemples** :

❌ **MAUVAIS - Résultats simulés** :
```go
// Ne JAMAIS faire ça !
expectedTokens := 3  // Simulé manuellement
if actualTokens != expectedTokens {
    t.Errorf("Attendu %d tokens", expectedTokens)
}
```

✅ **BON - Extraction depuis le réseau** :
```go
// Extraire depuis le réseau RETE réel
actualTokens := 0
for _, terminal := range network.TerminalNodes {
    actualTokens += len(terminal.Memory.GetTokens())
}

// Vérifier en inspectant les tokens réels
for _, terminal := range network.TerminalNodes {
    for _, token := range terminal.Memory.GetTokens() {
        t.Logf("Token trouvé: %d faits", len(token.Facts))
        // Validation basée sur les données réelles
    }
}
```

## Instructions

### 1. Identifier le Test qui Échoue

Précise :
- **Nom du test** : `TestNomDuTest`
- **Module** : `rete/`, `constraint/`, `test/integration/`, etc.
- **Message d'erreur** : Copie l'erreur complète

### 2. Analyser le Test

1. **Lire le code du test** :
   - Comprendre ce que le test essaie de valider
   - Identifier les assertions qui échouent
   - Examiner les données de test utilisées

2. **Examiner le contexte** :
   - Fichiers de contraintes utilisés (`.constraint`)
   - Fichiers de faits utilisés (`.facts`)
   - Configuration du test

3. **Tracer l'exécution** :
   - Activer le mode verbose : `go test -v -run TestNomDuTest`
   - Examiner les logs de propagation RETE
   - Identifier où l'exécution diverge de l'attendu

### 3. Identifier la Cause Racine

Poser les questions :
- **Quoi** : Quelle assertion échoue exactement ?
- **Où** : Dans quel module/fichier se produit le problème ?
- **Quand** : À quel moment de l'exécution (parsing, construction réseau, propagation) ?
- **Pourquoi** : Quelle est la cause sous-jacente ?

4. Proposer et Implémenter une Correction

1. **Analyser l'impact** :
   - Quels autres tests/modules sont affectés ?
   - Y a-t-il des effets de bord ?
   - La correction peut-elle introduire des race conditions ?

2. **Implémenter la correction** :
   - Modifier le code nécessaire
   - Ajouter des tests si nécessaire
   - Documenter les changements

3. **Valider la correction** :
   - Relancer le test spécifique
   - 🏁 **Relancer avec race detector : `go test -race -run TestNomDuTest` (OBLIGATOIRE)**
   - Relancer tous les tests pour éviter les régressions
   - 🏁 **Vérifier race detector global : `make test-race` (OBLIGATOIRE)**
   - Vérifier le runner universel

## Critères de Succès

✅ La cause racine est identifiée et documentée
✅ Une correction est proposée et implémentée
✅ Le test qui échouait passe maintenant
🏁 **✅ `go test -race` passe sans race condition (OBLIGATOIRE)**
✅ Aucune régression sur les autres tests
🏁 **✅ `make test-race` passe sans erreur (OBLIGATOIRE)**
✅ Le runner universel passe toujours (58/58)

## Commandes Utiles

```bash
# Lancer un test spécifique en mode verbose
go test -v -run TestNomDuTest ./rete

# Lancer avec timeout plus long
go test -v -timeout 5m -run TestNomDuTest ./rete

# Afficher seulement les échecs
go test -v -run TestNomDuTest ./rete 2>&1 | grep -A10 "FAIL"

# 🏁 OBLIGATOIRE : Lancer avec race detector (détecte race conditions)
go test -race -run TestNomDuTest ./rete
# ⚠️ CRITICAL: Toujours exécuter avec -race pour détecter les race conditions
# Les race conditions ne sont détectées QUE par le flag -race
# Ne JAMAIS skip cette étape, même si plus lent (~10x)

# Lancer tous les tests du module
go test -v ./rete

# 🏁 OBLIGATOIRE : Vérifier qu'on n'a pas de régression (avec race detector)
make test && make test-race && make rete-unified
```

## Format de Réponse Attendu

```
=== ANALYSE DU TEST ÉCHOUÉ ===

1. Identification
   - Test : TestNomDuTest
   - Module : rete/
   - Erreur : [message d'erreur complet]

2. Cause Racine
   - Description du problème
   - Fichier/fonction concernée
   - Pourquoi ça échoue

3. Solution Proposée
   - Modifications à apporter
   - Fichiers à modifier
   - Impact sur le reste du code

4. Implémentation
   - [Code modifié]
   
5. Validation
   - Test spécifique : [PASS/FAIL]
   - Suite de tests : [X/Y passent]
   - Runner universel : [58/58 passent]

6. Documentation
   - Changements apportés
   - Raison des modifications
```

## Exemple d'Utilisation

```
Le test TestIncrementalPropagation échoue avec l'erreur 
"variable non liée: p". Peux-tu utiliser le prompt "debug-test" 
pour identifier et corriger le problème ?
```

## Checklist de Debugging

- [ ] J'ai lu le code du test
- [ ] J'ai compris ce qu'il teste
- [ ] J'ai examiné le message d'erreur complet
- [ ] J'ai tracé l'exécution en mode verbose
- [ ] J'ai identifié la cause racine
- [ ] **TESTS RETE** : Résultats extraits du réseau (pas simulés)
- [ ] **CODE GO** : Aucun hardcoding introduit
- [ ] **CODE GO** : Code générique avec paramètres
- [ ] 🏁 **`go test -race` exécuté sur le test corrigé (OBLIGATOIRE)**
- [ ] **Aucune race condition détectée**
- [ ] J'ai vérifié l'impact de ma correction
- [ ] J'ai testé la correction localement
- [ ] 🏁 **`make test-race` passé sans erreur (OBLIGATOIRE)**
- [ ] Aucune régression n'a été introduite
- [ ] La documentation est à jour si nécessaire

## Cas Courants d'Erreurs

### Erreur de Variables Non Liées
- **Symptôme** : `variable non liée: x`
- **Cause** : Évaluation de condition avant que toutes les variables soient disponibles
- **Solution** : Évaluation partielle ou vérification des variables disponibles

### Erreur de Parsing
- **Symptôme** : `erreur de parsing` ou `unexpected token`
- **Cause** : Syntaxe incorrecte dans fichier `.constraint`
- **Solution** : Vérifier la grammaire PEG et le fichier de contraintes

### Erreur de Propagation
- **Symptôme** : Tokens attendus non créés
- **Cause** : Conditions mal évaluées ou nœuds mal connectés
- **Solution** : Tracer la propagation et vérifier les conditions

### Erreur de Mémoire
- **Symptôme** : Tokens en double ou perdus
- **Cause** : Gestion incorrecte des mémoires (Left/Right/Result)
- **Solution** : Vérifier la logique de stockage dans les nœuds
- **⚠️ Important** : Toujours extraire les tokens réels du réseau, ne jamais simuler

## Notes

- **CRITIQUE** : Pour tests RETE, toujours extraire résultats du réseau réel
- **CRITIQUE** : Aucun hardcoding dans le code de correction
- **CRITIQUE** : Code générique et réutilisable uniquement
- Toujours vérifier que le problème n'existe pas déjà dans les issues GitHub
- Documenter les corrections non-évidentes
- Ajouter des tests de régression si nécessaire
- Mettre à jour les commentaires dans le code
- Valider avec go vet et golangci-lint