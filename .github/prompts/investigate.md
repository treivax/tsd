# 🔍 Investiguer un Comportement (Investigate)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu observes un comportement étrange, inattendu ou inexpliqué dans le projet, mais **sans erreur explicite**. Contrairement à un bug identifié ou un test qui échoue, il s'agit d'une investigation exploratoire pour comprendre ce qui se passe.

## Objectif

Investiguer et comprendre un comportement inhabituel ou inexpliqué du système, en explorant méthodiquement les différentes hypothèses jusqu'à identifier la cause racine.

## Différence avec Autres Prompts

| Prompt | Quand l'Utiliser |
|--------|------------------|
| `analyze-error` | ❌ Quand il y a une **erreur explicite** (message d'erreur, stack trace) |
| `debug-test` | ❌ Quand un **test échoue** (échec identifié, assertion) |
| `fix-bug` | ❌ Quand le **bug est identifié** et qu'on veut le corriger |
| `investigate` | ✅ Quand le comportement est **étrange mais pas d'erreur** |

**Exemples d'utilisation de `investigate`** :
- "Les tokens se propagent bizarrement mais pas d'erreur"
- "La performance est anormalement lente sur certains cas"
- "Le réseau RETE a une structure inattendue"
- "Certaines règles ne génèrent pas de tokens alors qu'elles devraient"
- "Comportement différent entre deux exécutions similaires"

## Instructions

### PHASE 1 : OBSERVATION (Documenter le Comportement)

#### 1.1 Décrire le Comportement Observé

**Être très spécifique** :

```markdown
## Comportement Observé

**Quoi** : Les tokens ne se propagent pas aux TerminalNodes dans certains cas

**Quand** : Seulement avec des jointures 3-way et plus de 100 faits

**Où** : Fichier `rete/node_join.go`, fonction `propagateToChildren`

**Fréquence** :
- Toujours : ❌
- Parfois : ✅ (environ 30% des cas)
- Une seule fois : ❌

**Impact** :
- Bloquant : ❌
- Gênant : ✅
- Mineur : ❌

**Depuis quand** :
- Dernière modification de `evaluateJoinConditions` (commit abc123)
- Avant ça : comportement normal
```

#### 1.2 Définir le Comportement Attendu

**Contraste attendu vs observé** :

```markdown
## Attendu vs Observé

### Comportement Attendu
- 5 tokens devraient arriver au TerminalNode
- Propagation complète en < 100ms
- Toutes les jointures résolues

### Comportement Observé
- Seulement 3 tokens arrivent
- Propagation prend 500ms
- 2 jointures non résolues (raison inconnue)

### Différence
- **2 tokens manquants** (où sont-ils ?)
- **Performance dégradée** de 5x (pourquoi ?)
- **Jointures incomplètes** (quelle condition échoue ?)
```

#### 1.3 Collecter les Informations de Contexte

**Environnement** :
```bash
# Version Go
go version

# Commit actuel
git log -1 --oneline

# Fichiers modifiés récemment
git diff --name-only HEAD~5..HEAD

# Configuration
cat .env  # Si applicable
```

**Données du problème** :
- Fichiers `.constraint` et `.facts` concernés
- Logs pertinents
- Métriques (si disponibles)
- Captures d'écran (si interface)

### PHASE 2 : REPRODUCTION (Isoler le Problème)

#### 2.1 Créer un Cas de Reproduction Minimal

**Objectif** : Reproduire le comportement avec le minimum de données

**Méthode** :
1. Partir du cas complet qui montre le comportème
2. Supprimer progressivement des éléments
3. Identifier le minimum nécessaire pour reproduire

**Exemple** :
```constraint
# Cas complet (100 règles) → Comportement étrange
# ... réduction progressive ...
# Cas minimal (2 règles) → Comportement étrange toujours présent

# CAS MINIMAL DE REPRODUCTION :
{a: TypeA}, {b: TypeB}, {c: TypeC} /
    a.id == b.aId,
    b.id == c.bId
==> result(a, b, c)
```

```json
{
  "facts": [
    {"type": "TypeA", "data": {"id": 1}},
    {"type": "TypeB", "data": {"id": 2, "aId": 1}},
    {"type": "TypeC", "data": {"id": 3, "bId": 2}}
  ]
}
```

**Validation** :
```bash
# Le cas minimal reproduit-il le problème ?
make rete-run CONSTRAINT=minimal.constraint FACTS=minimal.facts

# Oui → Bon cas minimal ✅
# Non → Continuer à ajuster
```

#### 2.2 Tester les Variations

**Identifier ce qui change le comportement** :

```markdown
## Tests de Variations

### Variation 1 : Nombre de faits
- 10 faits : Comportement normal ✅
- 50 faits : Comportement normal ✅
- 100 faits : Comportement étrange ❌
- 200 faits : Comportement étrange ❌

**Hypothèse** : Seuil autour de 100 faits

### Variation 2 : Nombre de jointures
- 1 jointure : Normal ✅
- 2 jointures : Normal ✅
- 3 jointures : Étrange ❌

**Hypothèse** : Problème avec 3+ jointures

### Variation 3 : Type de données
- Integers : Étrange ❌
- Strings : Étrange ❌
- Mixed : Étrange ❌

**Hypothèse** : Pas lié au type de données

### Variation 4 : Ordre de soumission
- Ordre A→B→C : Étrange ❌
- Ordre C→B→A : Normal ✅
- Ordre B→A→C : Normal ✅

**Hypothèse** : ⭐ Problème avec l'ordre de soumission !
```

#### 2.3 Mesurer et Quantifier

**Collecter des métriques** :

```go
// Ajouter des logs de mesure
func investigateTokenPropagation() {
    start := time.Now()
    
    log.Printf("🔍 INVESTIGATION: Début propagation")
    log.Printf("🔍 Nombre de faits: %d", len(facts))
    log.Printf("🔍 Nombre de nœuds: %d", len(network.Nodes))
    
    // ... propagation ...
    
    elapsed := time.Since(start)
    log.Printf("🔍 Tokens créés: %d", tokenCount)
    log.Printf("🔍 Temps écoulé: %v", elapsed)
    log.Printf("🔍 Tokens attendus: %d", expectedCount)
    
    if tokenCount != expectedCount {
        log.Printf("⚠️  ANOMALIE: %d tokens manquants", expectedCount-tokenCount)
    }
}
```

**Analyser les patterns** :
```bash
# Exécuter plusieurs fois et collecter
for i in {1..10}; do
    echo "Run $i:"
    make rete-run CONSTRAINT=test.constraint FACTS=test.facts 2>&1 | grep "tokens"
done

# Analyser la variance
# Si variance élevée → problème de timing/race condition ?
# Si variance nulle → problème déterministe
```

### PHASE 3 : HYPOTHÈSES (Formuler des Théories)

#### 3.1 Brainstorming des Causes Possibles

**Catégories de causes** :

**1. Logique Métier** :
- Conditions mal évaluées
- Jointures incorrectes
- Variables non liées
- Ordre d'évaluation

**2. État du Système** :
- Mémoire corrompue
- État partagé entre appels
- Cache invalide
- Structures de données inconsistantes

**3. Concurrence** :
- Race conditions
- Deadlocks
- Ordre non déterministe
- Synchronisation manquante

**4. Performance** :
- Algorithme inefficace
- Boucles infinies cachées
- Allocations excessives
- Garbage collection

**5. Dépendances** :
- Version incompatible
- Bug dans bibliothèque externe
- Configuration incorrecte

**6. Environnement** :
- Différence OS
- Variables d'environnement
- Permissions fichiers
- Ressources limitées

#### 3.2 Prioriser les Hypothèses

**Critères de priorisation** :
1. **Probabilité** : Quelle hypothèse est la plus probable ?
2. **Impact** : Quel serait l'impact si c'est ça ?
3. **Facilité de test** : Peut-on tester rapidement ?

**Template** :
```markdown
## Hypothèses Priorisées

### Hypothèse 1 : Ordre de soumission des faits ⭐⭐⭐
**Probabilité** : Haute (variations le confirment)
**Impact** : Moyen (workaround possible)
**Test** : Facile (changer ordre)
**Statut** : À tester en priorité

### Hypothèse 2 : Buffer overflow à 100+ faits ⭐⭐
**Probabilité** : Moyenne (seuil observé à 100)
**Impact** : Élevé (limite d'utilisation)
**Test** : Facile (tester différentes tailles)
**Statut** : À tester ensuite

### Hypothèse 3 : Race condition dans propagation ⭐
**Probabilité** : Faible (comportement déterministe)
**Impact** : Élevé (bugs intermittents)
**Test** : Difficile (tests race, profiling)
**Statut** : Si les autres échouent
```

### PHASE 4 : EXPÉRIMENTATION (Tester les Hypothèses)

#### 4.1 Tester Chaque Hypothèse Méthodiquement

**Pour chaque hypothèse** :

```markdown
## Test de l'Hypothèse 1 : Ordre de soumission

### Setup
- Créer 3 fichiers .facts avec ordres différents
- Même contenu, juste ordre changé

### Protocole
1. Exécuter avec ordre A→B→C
2. Exécuter avec ordre C→B→A
3. Exécuter avec ordre B→A→C
4. Comparer résultats

### Résultats
| Ordre | Tokens | Temps | Comportement |
|-------|--------|-------|--------------|
| A→B→C | 3 | 500ms | ❌ Étrange |
| C→B→A | 5 | 80ms | ✅ Normal |
| B→A→C | 5 | 85ms | ✅ Normal |

### Conclusion
✅ **HYPOTHÈSE CONFIRMÉE** : L'ordre A→B→C cause le problème

### Prochaine Étape
Investiguer pourquoi cet ordre spécifique pose problème
```

#### 4.2 Ajouter de l'Instrumentation

**Logs stratégiques** :

```go
// Dans les points clés du code
func (j *JoinNode) Activate(token *Token) {
    log.Printf("🔍 JoinNode.Activate: token=%v, leftMemory=%d, rightMemory=%d",
        token, len(j.LeftMemory), len(j.RightMemory))
    
    // ... logique ...
    
    log.Printf("🔍 JoinNode.Activate: produced %d new tokens", newTokens)
}
```

**Dumps d'état** :

```go
func dumpNetworkState(network *Network) {
    log.Println("🔍 ========== NETWORK STATE ==========")
    for i, node := range network.Nodes {
        log.Printf("🔍 Node %d: type=%T, children=%d", i, node, len(node.Children))
        if jn, ok := node.(*JoinNode); ok {
            log.Printf("🔍   - LeftMemory: %d tokens", len(jn.LeftMemory))
            log.Printf("🔍   - RightMemory: %d tokens", len(jn.RightMemory))
        }
    }
    log.Println("🔍 ===================================")
}
```

**Traces de propagation** :

```go
var propagationTrace []string

func tracePropagation(from, to Node, token *Token) {
    trace := fmt.Sprintf("%T → %T: token=%v", from, to, token)
    propagationTrace = append(propagationTrace, trace)
}

func dumpTrace() {
    log.Println("🔍 ========== PROPAGATION TRACE ==========")
    for i, t := range propagationTrace {
        log.Printf("🔍 %d: %s", i, t)
    }
    log.Println("🔍 ========================================")
}
```

#### 4.3 Utiliser les Outils de Debug

**Profiling CPU** :
```bash
# Profiling CPU
go test -cpuprofile=cpu.prof -bench=. ./rete
go tool pprof cpu.prof
# (pprof) top10
# (pprof) list functionName
```

**Profiling Mémoire** :
```bash
# Profiling mémoire
go test -memprofile=mem.prof -bench=. ./rete
go tool pprof mem.prof
# (pprof) top10
# (pprof) list functionName
```

**Race Detector** :
```bash
# Détecter race conditions
go test -race ./rete
go build -race ./cmd/rete-runner
```

**Debugger** :
```bash
# Utiliser delve
dlv test ./rete -- -test.run TestProblematic
(dlv) break rete/node_join.go:123
(dlv) continue
(dlv) print token
(dlv) next
```

### PHASE 5 : ANALYSE (Comprendre la Cause)

#### 5.1 Identifier la Cause Racine

**Analyser les résultats des tests** :

```markdown
## Cause Racine Identifiée

### Symptôme
Tokens manquants avec ordre de soumission A→B→C

### Cause Directe
Les tokens de TypeA arrivent avant que le JoinNode ait reçu
les tokens de TypeB dans sa RightMemory

### Cause Racine
Le réseau RETE est construit avec l'hypothèse que les faits
arrivent dans un certain ordre (B avant A), mais ce n'est pas garanti

### Mécanisme
1. Fait A arrive → propagé à AlphaNode[A] → JoinNode
2. JoinNode cherche match dans RightMemory (TypeB)
3. RightMemory est vide → Pas de match
4. Token A stocké dans LeftMemory
5. Fait B arrive → propagé à AlphaNode[B] → RightMemory du JoinNode
6. ❌ PROBLÈME : La logique ne re-teste pas les tokens existants 
   dans LeftMemory après ajout dans RightMemory

### Preuve
Ajout de logs confirme que:
- LeftMemory contient bien les tokens A
- RightMemory reçoit bien les tokens B
- Mais la propagation ne se fait pas rétroactivement
```

#### 5.2 Comprendre l'Impact

**Évaluer la portée** :

```markdown
## Impact du Problème

### Scope
- **Affecté** : Toutes les jointures 2+ où l'ordre des faits n'est pas contrôlé
- **Non affecté** : AlphaNodes simples, règles sans jointure
- **Sévérité** : Haute (résultats incorrects)

### Cas d'Usage Impactés
1. Runner universel avec fichiers .facts non ordonnés : ✅ Affecté
2. Soumission interactive de faits : ✅ Affecté
3. Tests avec ordre garanti : ❌ Non affecté (d'où absence de détection)

### Workarounds Possibles
1. **Court terme** : Ordonner les faits dans les fichiers .facts
2. **Moyen terme** : Ajouter un tri automatique avant soumission
3. **Long terme** : Corriger la logique de propagation rétroactive

### Risque de Régression
- Risque élevé si correction mal faite
- Tests de régression nécessaires
- Valider avec runner universel (58 tests)
```

### PHASE 6 : DOCUMENTATION (Partager les Découvertes)

#### 6.1 Documenter l'Investigation

**Template de rapport** :

```markdown
# 🔍 Rapport d'Investigation : Tokens Manquants dans Jointures

**Date** : 2025-11-26  
**Investigateur** : [Nom]  
**Durée** : 4 heures  
**Statut** : ✅ Cause identifiée

---

## 📋 Résumé Exécutif

**Problème** : Tokens manquants lors de jointures avec certains ordres de faits

**Cause Racine** : Absence de propagation rétroactive dans les JoinNodes

**Impact** : Haute sévérité - Résultats incorrects dans scénarios réels

**Solution Recommandée** : Implémenter re-évaluation des tokens en mémoire

---

## 🎯 Comportement Observé

[Description détaillée du comportement étrange]

## 🔬 Investigation

### Étapes Réalisées
1. ✅ Création cas de reproduction minimal
2. ✅ Tests de variations (ordre, taille, types)
3. ✅ Formulation hypothèses
4. ✅ Expérimentation systématique
5. ✅ Instrumentation du code
6. ✅ Analyse des traces

### Hypothèses Testées
- ❌ Race condition : Éliminée (comportement déterministe)
- ❌ Buffer overflow : Éliminée (problème dès 3 faits)
- ✅ Ordre de soumission : **CONFIRMÉE**

## 💡 Cause Racine

[Explication détaillée de la cause]

### Diagramme

```
Ordre problématique (A→B→C):

A → AlphaNode[A] → JoinNode
                    ↓ (RightMemory vide)
                    ❌ Pas de match
                    ↓
                   LeftMemory (stockage)

B → AlphaNode[B] → RightMemory du JoinNode
                    ↓
                   ❌ Ne re-vérifie PAS LeftMemory !

Résultat : Tokens en LeftMemory ne sont jamais matchés
```

## 📊 Impact

[Évaluation de l'impact]

## 🛠️ Solutions Possibles

### Solution 1 : Propagation Rétroactive (Recommandée)
**Description** : Quand un token arrive en RightMemory, re-tester tous les tokens en LeftMemory

**Avantages** :
- ✅ Corrige le problème à la source
- ✅ Pas de contrainte sur l'ordre
- ✅ Conforme à l'algorithme RETE standard

**Inconvénients** :
- ⚠️ Complexité d'implémentation moyenne
- ⚠️ Légère perte de performance (acceptable)

**Effort** : 2-3 jours

### Solution 2 : Tri Automatique des Faits
**Description** : Trier les faits par type avant soumission

**Avantages** :
- ✅ Simple à implémenter
- ✅ Workaround rapide

**Inconvénients** :
- ❌ Ne corrige pas la cause racine
- ❌ Peut masquer le problème
- ❌ Ordre "optimal" dépend du réseau

**Effort** : 1 jour

### Solution 3 : Soumission en Deux Passes
**Description** : Soumettre tous les faits, puis déclencher propagation

**Avantages** :
- ✅ Relativement simple

**Inconvénients** :
- ❌ Change le comportement de l'API
- ❌ Perte de propagation incrémentale

**Effort** : 2 jours

## 🎯 Recommandation

**Implémenter Solution 1 (Propagation Rétroactive)**

**Plan d'action** :
1. Créer issue GitHub avec ce rapport
2. Implémenter la solution dans `rete/node_join.go`
3. Ajouter tests de régression
4. Valider avec runner universel (58 tests)
5. Mettre à jour documentation

**Timeline** : Sprint prochain

---

## 📎 Annexes

### Fichiers de Reproduction
- `investigation/minimal.constraint`
- `investigation/minimal.facts`

### Logs Pertinents
- `logs/investigation_2025-11-26.log`

### Code Instrumenté
- `rete/node_join.go` (branche `investigate/token-propagation`)

---

**Prochaines Étapes** :
- [ ] Créer issue GitHub
- [ ] Implémenter solution
- [ ] Tests de régression
- [ ] Code review
- [ ] Merge
```

#### 6.2 Créer une Issue GitHub (si applicable)

```markdown
## 🐛 Bug: Tokens manquants dans jointures selon ordre de soumission

**Type** : Bug  
**Sévérité** : Haute  
**Component** : rete/node_join.go

### Description

Les JoinNodes ne propagent pas correctement les tokens quand les faits
arrivent dans certains ordres.

### Reproduction

```bash
# Ordre A→B→C : ❌ 3 tokens au lieu de 5
make rete-run CONSTRAINT=test/order_abc.constraint FACTS=test/order_abc.facts

# Ordre C→B→A : ✅ 5 tokens (correct)
make rete-run CONSTRAINT=test/order_abc.constraint FACTS=test/order_cba.facts
```

### Cause Racine

Les JoinNodes ne re-évaluent pas les tokens en LeftMemory quand de
nouveaux tokens arrivent en RightMemory.

### Solution Proposée

Implémenter propagation rétroactive : quand un token arrive en RightMemory,
re-tester tous les tokens existants en LeftMemory pour matching.

### Impact

- Résultats incorrects dans scénarios réels
- Affecte toutes les jointures 2+
- Tests actuels ne détectent pas (ordre garanti)

### Rapport Complet

Voir `docs/investigations/token_propagation_2025-11-26.md`

### Checklist

- [x] Cas de reproduction minimal créé
- [x] Cause racine identifiée
- [x] Solution proposée
- [ ] Tests de régression ajoutés
- [ ] Solution implémentée
- [ ] Code review
- [ ] Documentation mise à jour
```

## Critères de Succès

### ✅ Compréhension

- [ ] Comportement étrange clairement documenté
- [ ] Comportement attendu défini
- [ ] Différence comprise et expliquée
- [ ] Cause racine identifiée
- [ ] Mécanisme du problème compris

### ✅ Méthodologie

- [ ] Cas de reproduction minimal créé
- [ ] Variations testées systématiquement
- [ ] Hypothèses formulées et priorisées
- [ ] Expérimentations menées méthodiquement
- [ ] Instrumentation ajoutée si nécessaire

### ✅ Documentation

- [ ] Rapport d'investigation rédigé
- [ ] Résultats des tests documentés
- [ ] Cause racine expliquée clairement
- [ ] Impact évalué précisément
- [ ] Solutions proposées avec avantages/inconvénients

### ✅ Action

- [ ] Issue créée (si bug identifié)
- [ ] Prochaines étapes définies
- [ ] Timeline établie
- [ ] Assignation faite

## Format de Réponse

```markdown
# 🔍 INVESTIGATION : [Titre du Comportement]

## 📋 Résumé

**Problème** : [Description courte]
**Durée Investigation** : [X heures]
**Statut** : [En cours / Cause trouvée / Bloqué]

## 🎯 Comportement Observé

[Description détaillée]

**Attendu** : [Ce qui devrait se passer]
**Observé** : [Ce qui se passe réellement]

## 🔬 Méthodologie

### Cas de Reproduction
[Fichiers .constraint et .facts minimaux]

### Variations Testées
| Variation | Résultat | Comportement |
|-----------|----------|--------------|
| [Test 1] | [Résultat] | ✅/❌ |
| [Test 2] | [Résultat] | ✅/❌ |

### Hypothèses
1. ❌ [Hypothèse 1] : Éliminée car [raison]
2. ❌ [Hypothèse 2] : Éliminée car [raison]
3. ✅ [Hypothèse 3] : **CONFIRMÉE** - [explication]

## 💡 Cause Racine

[Explication détaillée de la cause]

### Mécanisme
[Comment le problème se produit, étape par étape]

### Preuve
[Logs, traces, mesures qui confirment]

## 📊 Impact

**Sévérité** : [Haute/Moyenne/Faible]
**Scope** : [Ce qui est affecté]
**Workarounds** : [Solutions temporaires possibles]

## 🛠️ Solutions Proposées

### Solution 1 : [Nom] (Recommandée)
- **Description** : [...]
- **Avantages** : [...]
- **Inconvénients** : [...]
- **Effort** : [X jours]

### Solution 2 : [Nom]
- **Description** : [...]
- **Avantages** : [...]
- **Inconvénients** : [...]
- **Effort** : [X jours]

## 🎯 Recommandation

[Quelle solution choisir et pourquoi]

**Plan d'action** :
1. [Étape 1]
2. [Étape 2]
3. [Étape 3]

**Timeline** : [Estimation]

## 📎 Fichiers Générés

- Cas minimal : `investigation/[nom].constraint`, `investigation/[nom].facts`
- Logs : `logs/investigation_[date].log`
- Code instrumenté : `rete/[fichier].go` (branche `investigate/[nom]`)
- Rapport complet : `docs/investigations/[nom]_[date].md`

## 🔗 Prochaines Étapes

- [ ] Créer issue GitHub
- [ ] Implémenter solution
- [ ] Tests de régression
- [ ] Documentation
- [ ] Review et merge
```

## Exemple d'Utilisation

```
J'observe un comportement étrange : lors de l'exécution de certaines contraintes
avec le runner universel, seulement 3 tokens arrivent aux TerminalNodes alors
que j'en attends 5. Le plus bizarre, c'est que si je change l'ordre des faits
dans le fichier .facts, j'obtiens bien 5 tokens.

Il n'y a pas d'erreur, les tests passent, mais le comportement est incohérent
selon l'ordre de soumission des faits.

Utilise le prompt "investigate" pour m'aider à comprendre ce qui se passe.
```

## Checklist d'Investigation

### Avant de Commencer

- [ ] Comportement étrange clairement observé (pas une erreur explicite)
- [ ] Comportement attendu défini
- [ ] Contexte documenté (version, commit, environnement)
- [ ] Temps alloué à l'investigation (éviter investigation sans fin)

### Pendant l'Investigation

- [ ] Cas de reproduction minimal créé
- [ ] Variations testées (au moins 3-4)
- [ ] Hypothèses listées et priorisées
- [ ] Au moins 3 hypothèses testées
- [ ] Instrumentation ajoutée si nécessaire
- [ ] Notes prises au fur et à mesure

### Identification de la Cause

- [ ] Cause racine identifiée (pas juste symptôme)
- [ ] Mécanisme compris en détail
- [ ] Preuves collectées (logs, traces, mesures)
- [ ] Impact évalué
- [ ] Solutions possibles identifiées

### Documentation

- [ ] Rapport d'investigation rédigé
- [ ] Fichiers de reproduction sauvegardés
- [ ] Logs archivés
- [ ] Issue créée (si applicable)
- [ ] Code instrumenté committé sur branche

### Clôture

- [ ] Prochaines étapes définies
- [ ] Assignation faite
- [ ] Timeline établie
- [ ] Équipe informée

## Commandes Utiles

```bash
# Reproduction
make rete-run CONSTRAINT=test.constraint FACTS=test.facts

# Logs détaillés
RETE_DEBUG=1 make rete-run CONSTRAINT=test.constraint FACTS=test.facts

# Profiling CPU
go test -cpuprofile=cpu.prof -bench=BenchmarkProblematic ./rete
go tool pprof cpu.prof

# Profiling mémoire
go test -memprofile=mem.prof -bench=BenchmarkProblematic ./rete
go tool pprof mem.prof

# Race detector
go test -race ./rete
go build -race ./cmd/rete-runner

# Traces
go test -trace=trace.out -run TestProblematic ./rete
go tool trace trace.out

# Debugger
dlv test ./rete -- -test.run TestProblematic
dlv debug ./cmd/rete-runner -- test.constraint test.facts

# Exécutions multiples (détecter non-déterminisme)
for i in {1..20}; do
    echo "Run $i:"
    make rete-run CONSTRAINT=test.constraint FACTS=test.facts | grep "tokens:"
done

# Diff entre deux exécutions
make rete-run CONSTRAINT=test.constraint FACTS=order1.facts > /tmp/run1.log
make rete-run CONSTRAINT=test.constraint FACTS=order2.facts > /tmp/run2.log
diff /tmp/run1.log /tmp/run2.log

# Git bisect (si régression)
git bisect start
git bisect bad HEAD
git bisect good v1.0.0
git bisect run make test
```

## Bonnes Pratiques

### Investigation

- **Systématique** : Tester hypothèses une par une, méthodiquement
- **Documentation** : Noter tout au fur et à mesure, pas après
- **Minimal** : Toujours créer cas de reproduction minimal
- **Mesurable** : Quantifier le comportème (temps, tokens, mémoire)
- **Objectif** : S'en tenir aux faits, éviter les suppositions

### Hypothèses

- **Multiples** : Formuler plusieurs hypothèses, pas juste une
- **Priorisées** : Tester les plus probables d'abord
- **Falsifiables** : Chaque hypothèse doit être testable
- **Éliminées** : Documenter pourquoi une hypothèse est éliminée

### Instrumentation

- **Ciblée** : Ajouter logs aux endroits stratégiques
- **Temporaire** : Marquer le code d'investigation (commentaires `// INVESTIGATION`)
- **Verbeux** : Ne pas hésiter à logger beaucoup d'infos
- **Structuré** : Préfixer les logs (ex: `🔍 INVESTIGATION:`)

### Documentation

- **Immédiate** : Documenter pendant, pas après
- **Complète** : Inclure tentatives ratées (c'est instructif)
- **Partageable** : Rapport compréhensible par d'autres
- **Actionnable** : Toujours finir avec prochaines étapes

## Anti-Patterns à Éviter

### ❌ Investigation Sans Reproduction
```
❌ "Je pense savoir ce que c'est" sans reproduire
✅ Toujours créer cas de reproduction avant de conclure
```

### ❌ Suppositions Non Testées
```
❌ "Ça doit être un problème de X" sans tester
✅ Formuler hypothèse ET la tester avec expérience
```

### ❌ Modifications Aléatoires
```
❌ Changer du code au hasard pour "voir si ça marche"
✅ Comprendre d'abord, modifier ensuite avec intention
```

### ❌ Investigation Sans Fin
```
❌ Investiguer indéfiniment sans conclusion
✅ Se fixer une limite de temps, escalader si nécessaire
```

### ❌ Pas de Documentation
```
❌ Investigation dans sa tête, rien d'écrit
✅ Documenter hypothèses, tests, résultats au fur et à mesure
```

### ❌ Ignorer les Patterns
```
❌ Tester variations sans chercher pattern
✅ Analyser les résultats pour identifier tendances
```

## Outils Recommandés

### Profiling et Analyse
- `pprof` - Profiling CPU/mémoire
- `go tool trace` - Traces d'exécution
- `go test -race` - Détection race conditions
- `delve (dlv)` - Debugger Go

### Logging
- `log` package standard
- `logrus` - Logs structurés
- `zap` - Logs haute performance

### Instrumentation
- `expvar` - Variables exportées
- `prometheus` - Métriques
- Custom logging avec préfixes

### Analyse de Données
- `grep`, `awk`, `sed` - Analyse logs
- `jq` - Analyse JSON
- Scripts Python/Go pour analyse personnalisée

## Ressources

- [Makefile](../../Makefile) - Commandes disponibles
- [Debugging Go](https://go.dev/doc/diagnostics) - Guide officiel
- [Delve](https://github.com/go-delve/delve) - Debugger Go
- [pprof](https://go.dev/blog/pprof) - Profiling

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Mainteneur** : Équipe TSD