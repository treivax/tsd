# 🔍 Analyser une Erreur ou un Problème

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Une erreur s'est produite dans le système (compilation, exécution, tests) et tu as besoin d'aide pour la comprendre et la résoudre.

## Objectif

Analyser une erreur, identifier sa cause, et proposer une solution adaptée.

## ⚠️ RÈGLES STRICTES - CODE GOLANG ET TESTS RETE

### 🚫 INTERDICTIONS ABSOLUES

1. **CODE GOLANG - AUCUN HARDCODING** :
   - ❌ Pas de valeurs en dur dans le code de correction
   - ❌ Pas de "magic numbers" ou "magic strings"
   - ❌ Pas de chemins de fichiers hardcodés
   - ❌ Pas de configurations hardcodées
   - ✅ Utiliser des constantes nommées
   - ✅ Utiliser des variables de configuration
   - ✅ Utiliser des paramètres de fonction

2. **CODE TOUJOURS GÉNÉRIQUE** :
   - ✅ Fonctions réutilisables avec paramètres
   - ✅ Interfaces pour abstraction
   - ✅ Code extensible sans modification
   - ❌ Pas de code spécifique à un cas d'usage

3. **TESTS RETE - AUCUNE SIMULATION** :
   - ❌ Pas de résultats hardcodés ou simulés
   - ❌ Pas de mock des résultats du réseau RETE
   - ❌ Pas de calcul manuel des tokens attendus
   - ✅ **TOUJOURS** extraire les résultats du réseau RETE réel
   - ✅ **TOUJOURS** interroger les TerminalNodes
   - ✅ **TOUJOURS** inspecter les mémoires (Left/Right/Result)

### ✅ BONNES PRATIQUES GO OBLIGATOIRES

1. **Conventions Go** :
   - Respect de Effective Go
   - Nommage idiomatique (MixedCaps pour export)
   - Gestion explicite des erreurs (pas de panic sauf critique)
   - go fmt et goimports appliqués
   - Commentaires GoDoc pour exports

2. **Architecture** :
   - Single Responsibility Principle
   - Interfaces petites et focalisées
   - Composition over inheritance
   - Dependency injection
   - Découplage fort

3. **Qualité** :
   - Code auto-documenté
   - Complexité cyclomatique < 15
   - Fonctions < 50 lignes (sauf justification)
   - Pas de duplication (DRY)
   - go vet et golangci-lint sans erreur

**Exemples** :

❌ **MAUVAIS - Hardcodé** :
```go
func FixError() error {
    timeout := 30  // Magic number !
    if userID == "special-123" {  // Hardcodé !
        // correction spécifique
    }
}
```

✅ **BON - Générique** :
```go
const DefaultTimeout = 30 * time.Second

type ErrorHandler interface {
    Handle(userID string) error
}

func FixError(timeout time.Duration, handler ErrorHandler) error {
    // Code générique et réutilisable
}
```

❌ **MAUVAIS - Test RETE simulé** :
```go
// Ne JAMAIS faire ça !
expectedTokens := 3  // Simulé manuellement
```

✅ **BON - Test RETE avec extraction** :
```go
// Extraire depuis le réseau RETE réel
actualTokens := 0
for _, terminal := range network.TerminalNodes {
    actualTokens += len(terminal.Memory.GetTokens())
}
```

## Instructions

### 1. Fournir l'Erreur Complète

**Partage** :
- **Message d'erreur complet** : Copie tout le stack trace
- **Contexte d'exécution** : Quelle commande a produit l'erreur ?
- **Environnement** : Version Go, OS, etc. si pertinent
- **Moment** : Quand survient l'erreur (compilation, runtime, tests) ?

**Exemple** :
```
Commande : make test
Erreur :
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x...]
```

### 2. Décrire le Comportement

**Précise** :
- **Comportement attendu** : Ce qui devrait se passer
- **Comportement observé** : Ce qui se passe réellement
- **Reproductibilité** : L'erreur survient-elle toujours ?
- **Changements récents** : Qu'est-ce qui a été modifié récemment ?

### 3. Contexte Additionnel

**Si disponible, fournis** :
- Fichiers `.constraint` utilisés
- Fichiers `.facts` utilisés
- Configuration du réseau RETE
- Logs de propagation (si mode verbose)
- Code récemment modifié

## Analyse Structurée

### Phase 1 : Classification de l'Erreur

**Type d'erreur** :
- [ ] Erreur de compilation
- [ ] Erreur de runtime (panic)
- [ ] Erreur logique (résultat incorrect)
- [ ] Erreur de test
- [ ] Erreur de performance
- [ ] Erreur de mémoire

**Criticité** :
- [ ] Bloquante (empêche l'exécution)
- [ ] Majeure (fonctionnalité cassée)
- [ ] Mineure (comportement dégradé)
- [ ] Cosmétique (affichage/format)

### Phase 2 : Localisation

**Où se produit l'erreur** :
- Module : `rete/`, `constraint/`, `test/`, etc.
- Fichier : Nom du fichier concerné
- Fonction : Fonction/méthode où survient l'erreur
- Ligne : Numéro de ligne si disponible

### Phase 3 : Investigation

**Questions à explorer** :
1. **Stack trace** : Quelle est la chaîne d'appels ?
2. **État du système** : Valeurs des variables au moment de l'erreur ?
3. **Données d'entrée** : Quelles données causent le problème ?
4. **Conditions** : Quelles conditions déclenchent l'erreur ?
5. **Historique** : L'erreur existait-elle avant ?

### Phase 4 : Hypothèses

**Causes possibles** :
- Pointeur nil
- Index hors limites
- Race condition
- Variable non initialisée
- Type incompatible
- Logique incorrecte
- Dépendance manquante

## Catégories d'Erreurs Communes

### 1. Erreurs de Variables Non Liées

**Symptômes** :
```
❌ Erreur: variable non liée: p (variables disponibles: [u o])
```

**Causes courantes** :
- Évaluation de condition avant disponibilité des variables
- Jointure multi-variables incomplète
- Ordre de propagation incorrect

**Solution type** :
- Vérifier variables disponibles avant évaluation
- Implémenter évaluation partielle
- Ajuster l'ordre de propagation

### 2. Erreurs de Parsing

**Symptômes** :
```
❌ Erreur de parsing ligne 5: unexpected token "{"
```

**Causes courantes** :
- Syntaxe incorrecte dans `.constraint`
- Grammaire PEG non à jour
- Caractères spéciaux non échappés

**Solution type** :
- Valider la syntaxe du fichier `.constraint`
- Vérifier la grammaire PEG
- Tester avec un fichier minimal

### 3. Panics / Nil Pointer

**Symptômes** :
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Causes courantes** :
- Accès à un pointeur nil
- Map/slice non initialisé
- Retour de fonction nil non vérifié

**Solution type** :
- Vérifier initialisation des variables
- Ajouter des gardes nil checks
- Utiliser les valeurs par défaut

### 4. Erreurs de Propagation RETE

**Symptômes** :
```
✅ Fait soumis mais aucun token terminal créé
```

**Causes courantes** :
- Conditions mal évaluées
- Nœuds mal connectés
- Mémoire (Left/Right) incorrecte

**Solution type** :
- Tracer la propagation en verbose
- Vérifier les conditions de jointure
- Valider la construction du réseau

### 5. Erreurs de Tests

**Symptômes** :
```
--- FAIL: TestNom (0.00s)
    test.go:42: attendu 5, reçu 3
```

**Causes courantes** :
- Données de test incorrectes
- Assertion trop stricte
- Race condition dans le test
- Test non isolé

**Solution type** :
- Vérifier les données d'entrée
- Revoir les assertions
- Isoler le test
- Utiliser -race flag

## Format de Réponse Attendu

```
=== ANALYSE DE L'ERREUR ===

1. Classification
   - Type : [Type d'erreur]
   - Criticité : [Bloquante/Majeure/Mineure]
   - Module : [rete/constraint/test/...]

2. Localisation
   - Fichier : [nom du fichier]
   - Fonction : [nom de la fonction]
   - Ligne : [numéro si disponible]

3. Cause Racine
   - Description détaillée du problème
   - Pourquoi ça se produit
   - Conditions de déclenchement

4. Impact
   - Fonctionnalités affectées
   - Portée du problème
   - Urgence de la correction

5. Solution Proposée
   - Approche de correction
   - Fichiers à modifier
   - Tests à ajouter/modifier
   - **⚠️ VÉRIFICATION** : Aucun hardcoding introduit
   - **⚠️ VÉRIFICATION** : Code générique avec paramètres
   - **⚠️ VÉRIFICATION** : Tests RETE avec extraction réelle

6. Plan d'Action
   - Étape 1 : ...
   - Étape 2 : ...
   - Étape 3 : ...

7. Prévention Future
   - Comment éviter ce problème à l'avenir
   - Tests à ajouter
   - Documentation à améliorer
```

## Commandes de Diagnostic

```bash
# Exécuter avec stack trace complet
go test -v -run TestNom ./rete 2>&1 | tee error.log

# Exécuter avec race detector
go test -race -v -run TestNom ./rete

# Exécuter avec couverture
go test -cover -coverprofile=coverage.out ./rete

# Analyser la couverture
go tool cover -html=coverage.out

# Vérifier les erreurs statiques
go vet ./...
golangci-lint run

# Profiling mémoire
go test -memprofile mem.prof -run TestNom ./rete
go tool pprof mem.prof

# Profiling CPU
go test -cpuprofile cpu.prof -run TestNom ./rete
go tool pprof cpu.prof
```

## Exemple d'Utilisation

```
J'ai cette erreur quand je lance make test:

panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x5a8c73]

goroutine 1 [running]:
github.com/treivax/tsd/rete.(*JoinNode).evaluateJoinConditions(...)
    /home/user/tsd/rete/node_join.go:265

Ça se produit dans le test TestIncrementalPropagation.

Peux-tu utiliser le prompt "analyze-error" pour m'aider ?
```

## Checklist d'Analyse

- [ ] Message d'erreur complet copié
- [ ] Commande qui a produit l'erreur identifiée
- [ ] Comportement attendu vs observé décrit
- [ ] Contexte d'exécution fourni
- [ ] Stack trace analysé
- [ ] Variables et état examinés
- [ ] Cause racine identifiée
- [ ] Solution proposée
- [ ] **AUCUN HARDCODING** dans le code de correction
- [ ] **CODE GÉNÉRIQUE** avec paramètres/interfaces
- [ ] **TESTS RETE** : Résultats extraits du réseau (pas simulés)
- [ ] **go vet et golangci-lint** sans erreur
- [ ] Tests de validation prévus

## Outils d'Aide

### Logging Verbeux

Activer les logs détaillés dans le code :
```go
fmt.Printf("🔍 DEBUG: variable=%v\n", variable)
fmt.Printf("🔍 DEBUG: état=%+v\n", structure)
```

### Breakpoints avec Delve

```bash
# Installer delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Débugger un test
dlv test ./rete -- -test.run TestNom

# Dans delve:
(dlv) break node_join.go:265
(dlv) continue
(dlv) print bindings
(dlv) next
```

### Diagnostiques RETE

```go
// Afficher l'état du réseau
network.PrintDiagnostics()

// Afficher les tokens dans un nœud
node.Memory.PrintTokens()

// Tracer la propagation
network.EnableVerboseMode(true)
```

## Notes Importantes

- **CRITIQUE** : Aucun hardcoding dans les corrections de code
- **CRITIQUE** : Code générique et réutilisable uniquement
- **CRITIQUE** : Tests RETE avec extraction réseau réel (pas de simulation)
- **Ne pas paniquer** : Toute erreur a une cause et une solution
- **Diviser pour régner** : Isoler le problème par élimination
- **Tester les hypothèses** : Valider chaque hypothèse une par une
- **Documenter** : Noter les découvertes pour référence future
- **Demander de l'aide** : N'hésite pas si tu bloques
- **Valider** : go vet et golangci-lint sur tout code modifié

## Ressources

- [Go Error Handling](https://go.dev/blog/error-handling-and-go)
- [Effective Go](https://go.dev/doc/effective_go)
- [Debugging with Delve](https://github.com/go-delve/delve)
- [Tests du projet](../../test/)

---

**Rappel** : Une erreur bien analysée est à moitié résolue !