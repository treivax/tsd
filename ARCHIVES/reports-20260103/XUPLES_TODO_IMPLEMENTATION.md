# 📋 TODOs pour Implémentation Complète - Module Xuples

**Date**: 2025-12-17  
**Version actuelle**: 1.0.0  
**Version cible**: 1.1.0  

---

## 🎯 Vue d'ensemble

Le module xuples est **fonctionnel et prêt pour l'intégration**, mais certaines actions définies dans l'API ne sont pas encore implémentées. Ce document détaille les actions nécessaires pour compléter l'implémentation.

---

## ⚠️ Actions Non Implémentées

### 1. Action Update(fact: any)

**Statut**: ⚠️ NON IMPLÉMENTÉ  
**Package concerné**: `rete`  
**Priorité**: HAUTE  

#### Description
L'action `Update` doit permettre de modifier un fait existant dans le réseau RETE et propager les changements à tous les tokens dépendants.

#### Signature TSD
```tsd
action Update(fact: any)
```

#### Spécifications Techniques

**Méthode à implémenter dans `rete.ReteNetwork`**:
```go
// UpdateFact met à jour un fait existant dans le réseau RETE.
//
// Cette méthode doit:
// 1. Localiser le fait existant dans le réseau par son ID
// 2. Mettre à jour ses attributs avec les nouvelles valeurs
// 3. Identifier tous les tokens qui référencent ce fait
// 4. Propager les changements aux nœuds dépendants
// 5. Re-évaluer les conditions affectées par les changements
//
// Paramètres:
//   - fact: fait avec attributs mis à jour (ID doit correspondre au fait existant)
//
// Retourne:
//   - error: si le fait n'existe pas ou si la mise à jour échoue
func (rn *ReteNetwork) UpdateFact(fact *Fact) error {
    // TODO: Implémenter
}
```

#### Algorithme Recommandé

1. **Validation**
   - Vérifier que fact != nil
   - Vérifier que fact.ID existe dans le réseau
   - Valider le type du fait

2. **Localisation**
   - Trouver tous les AlphaNodes qui contiennent ce fait
   - Identifier tous les tokens dans les BetaNodes qui référencent ce fait

3. **Mise à jour**
   - Mettre à jour les attributs du fait
   - Notifier les AlphaNodes de la modification

4. **Propagation**
   - Pour chaque token affecté:
     - Re-évaluer les conditions du BetaNode
     - Si la condition n'est plus satisfaite, retirer le token
     - Si la condition reste satisfaite, propager la modification

5. **Re-évaluation**
   - Propager aux nœuds terminaux
   - Déclencher les règles affectées si nécessaire

#### Tests à Ajouter

```go
func TestUpdateFact(t *testing.T) {
    // Test: mise à jour simple
    // Test: mise à jour avec propagation
    // Test: fait inexistant
    // Test: fait nil
    // Test: type invalide
    // Test: impact sur tokens
    // Test: re-évaluation conditions
}
```

#### Effort Estimé
- Implémentation: 2-3 jours
- Tests: 1 jour
- Documentation: 0.5 jour
- **Total**: ~4 jours

---

### 2. Action Insert(fact: any)

**Statut**: ⚠️ NON IMPLÉMENTÉ  
**Package concerné**: `rete`  
**Priorité**: HAUTE  

#### Description
L'action `Insert` doit permettre de créer et insérer dynamiquement un nouveau fait dans le réseau RETE pendant l'exécution des règles.

#### Signature TSD
```tsd
action Insert(fact: any)
```

#### Spécifications Techniques

**Méthode à implémenter dans `rete.ReteNetwork`**:
```go
// InsertFact insère un nouveau fait dans le réseau RETE.
//
// Cette méthode doit:
// 1. Valider le fait (type, attributs requis)
// 2. Générer un ID unique si non fourni
// 3. Insérer le fait via les nœuds alpha appropriés
// 4. Propager aux nœuds bêta selon les patterns de matching
// 5. Activer les règles qui matchent le nouveau fait
//
// Paramètres:
//   - fact: nouveau fait à insérer (ID généré automatiquement si vide)
//
// Retourne:
//   - error: si la validation échoue ou si l'insertion échoue
func (rn *ReteNetwork) InsertFact(fact *Fact) error {
    // TODO: Implémenter
}
```

#### Algorithme Recommandé

1. **Validation**
   - Vérifier que fact != nil
   - Valider le type du fait (doit correspondre à un type déclaré)
   - Vérifier que les attributs requis sont présents

2. **Préparation**
   - Si fact.ID est vide, générer un ID unique (UUID)
   - Enregistrer le fait dans la working memory

3. **Insertion dans AlphaNodes**
   - Trouver les AlphaNodes correspondant au type du fait
   - Évaluer les conditions alpha
   - Créer les tokens alpha pour les conditions satisfaites

4. **Propagation dans BetaNodes**
   - Propager les tokens aux BetaNodes
   - Évaluer les jointures
   - Créer les tokens combinés

5. **Activation**
   - Propager aux nœuds terminaux
   - Ajouter à l'agenda les règles activées
   - Ne PAS exécuter immédiatement (éviter récursion infinie)

#### Tests à Ajouter

```go
func TestInsertFact(t *testing.T) {
    // Test: insertion simple
    // Test: génération automatique ID
    // Test: propagation aux nœuds
    // Test: activation de règles
    // Test: fait nil
    // Test: type invalide
    // Test: attributs manquants
}
```

#### Effort Estimé
- Implémentation: 2-3 jours
- Tests: 1 jour
- Documentation: 0.5 jour
- **Total**: ~4 jours

---

### 3. Action Retract(id: string)

**Statut**: ⚠️ NON IMPLÉMENTÉ  
**Package concerné**: `rete`  
**Priorité**: HAUTE  

#### Description
L'action `Retract` doit permettre de supprimer un fait du réseau RETE et tous les tokens qui en dépendent (truth maintenance).

#### Signature TSD
```tsd
action Retract(id: string)
```

#### Spécifications Techniques

**Méthode à implémenter dans `rete.ReteNetwork`**:
```go
// RetractFact supprime un fait du réseau RETE et tous les tokens dépendants.
//
// Cette méthode implémente le truth maintenance system (TMS):
// 1. Localiser le fait par son ID
// 2. Identifier tous les tokens qui dépendent de ce fait
// 3. Propager la rétraction aux nœuds dépendants
// 4. Supprimer le fait et nettoyer les références
// 5. Désactiver les règles qui ne sont plus satisfaites
//
// Paramètres:
//   - id: identifiant unique du fait à rétracter
//
// Retourne:
//   - error: si le fait n'existe pas ou si la rétraction échoue
func (rn *ReteNetwork) RetractFact(id string) error {
    // TODO: Implémenter
}
```

#### Algorithme Recommandé

1. **Validation**
   - Vérifier que id != ""
   - Vérifier que le fait existe dans la working memory

2. **Identification des Dépendances**
   - Trouver tous les tokens alpha contenant ce fait
   - Trouver tous les tokens beta dépendants (via Parent chain)
   - Construire le graphe de dépendances

3. **Propagation de la Rétraction**
   - Pour chaque token dépendant (ordre inverse):
     - Retirer le token de son nœud
     - Propager la rétraction aux nœuds enfants
     - Retirer des nœuds terminaux si applicable

4. **Nettoyage**
   - Supprimer le fait de la working memory
   - Retirer de tous les AlphaNodes
   - Nettoyer les références

5. **Désactivation**
   - Retirer de l'agenda les activations invalides
   - Marquer les règles comme non satisfaites

#### Tests à Ajouter

```go
func TestRetractFact(t *testing.T) {
    // Test: rétraction simple
    // Test: propagation aux dépendances
    // Test: truth maintenance
    // Test: fait inexistant
    // Test: id vide
    // Test: multiples dépendances
    // Test: désactivation règles
}
```

#### Effort Estimé
- Implémentation: 3-4 jours (plus complexe - TMS)
- Tests: 1.5 jour
- Documentation: 0.5 jour
- **Total**: ~6 jours

---

## 📊 Effort Total Estimé

| Action | Implémentation | Tests | Documentation | Total |
|--------|----------------|-------|---------------|-------|
| Update | 2-3 jours | 1 jour | 0.5 jour | 4 jours |
| Insert | 2-3 jours | 1 jour | 0.5 jour | 4 jours |
| Retract | 3-4 jours | 1.5 jour | 0.5 jour | 6 jours |
| **TOTAL** | **7-10 jours** | **3.5 jours** | **1.5 jour** | **14 jours** |

**Note**: Estimation pour un développeur expérimenté avec bonne connaissance de RETE.

---

## 🔄 Plan d'Implémentation Recommandé

### Phase 1: Insert (Semaine 1)
**Raison**: Plus simple, pose les bases pour Update et Retract

1. Jour 1-3: Implémentation `InsertFact()`
2. Jour 4: Tests complets
3. Jour 5: Documentation et review

**Livrables**:
- ✅ `rete.ReteNetwork.InsertFact()` implémenté
- ✅ Tests passent (couverture > 80%)
- ✅ Documentation GoDoc complète
- ✅ Action `Insert` fonctionnelle dans builtin.go

### Phase 2: Update (Semaine 2)
**Raison**: Nécessite Insert pour les tests

1. Jour 1-3: Implémentation `UpdateFact()`
2. Jour 4: Tests complets
3. Jour 5: Documentation et review

**Livrables**:
- ✅ `rete.ReteNetwork.UpdateFact()` implémenté
- ✅ Tests passent (couverture > 80%)
- ✅ Documentation GoDoc complète
- ✅ Action `Update` fonctionnelle dans builtin.go

### Phase 3: Retract (Semaine 3)
**Raison**: Plus complexe, nécessite TMS

1. Jour 1-4: Implémentation `RetractFact()` + TMS
2. Jour 5-6: Tests complets
3. Jour 7: Documentation et review

**Livrables**:
- ✅ `rete.ReteNetwork.RetractFact()` implémenté
- ✅ Truth Maintenance System fonctionnel
- ✅ Tests passent (couverture > 80%)
- ✅ Documentation GoDoc complète
- ✅ Action `Retract` fonctionnelle dans builtin.go

### Phase 4: Validation Finale (Jour 15)
1. Tests d'intégration complets
2. Tests de régression
3. Documentation utilisateur
4. Release v1.1.0

---

## 📝 Checklist d'Implémentation

Pour chaque action, suivre cette checklist:

### Implémentation
- [ ] Méthode dans `rete.ReteNetwork` créée
- [ ] Validation des paramètres
- [ ] Algorithme implémenté
- [ ] Gestion d'erreurs robuste
- [ ] Thread-safety garantie

### Tests
- [ ] Tests unitaires (cas nominaux)
- [ ] Tests cas d'erreur
- [ ] Tests cas limites
- [ ] Tests de concurrence (race detector)
- [ ] Couverture > 80%

### Documentation
- [ ] GoDoc complet
- [ ] Exemples dans les commentaires
- [ ] README mis à jour
- [ ] CHANGELOG.md mis à jour

### Intégration
- [ ] Modification de `builtin.go` pour appeler la méthode
- [ ] Retrait du `return fmt.Errorf("not yet implemented")`
- [ ] Tests de `builtin.go` mis à jour
- [ ] Validation avec `make validate`

### Validation
- [ ] Tous les tests passent
- [ ] `go vet` OK
- [ ] `staticcheck` OK
- [ ] `errcheck` OK
- [ ] Aucune régression

---

## 🎯 Critères de Succès

Une implémentation est considérée comme réussie si:

1. ✅ La méthode dans `rete.ReteNetwork` fonctionne correctement
2. ✅ L'action dans `builtin.go` appelle la méthode RETE
3. ✅ Tous les tests passent (couverture > 80%)
4. ✅ Aucune régression dans les tests existants
5. ✅ Documentation complète (GoDoc + README)
6. ✅ Toutes les vérifications qualité passent
7. ✅ Thread-safety garantie (race detector OK)

---

## 📚 Ressources

### Code Existant
- `rete/network.go` - Structure du réseau RETE
- `rete/node.go` - Définition des nœuds
- `rete/token.go` - Gestion des tokens
- `rete/actions/builtin.go` - Actions natives

### Documentation
- [RETE Algorithm](https://en.wikipedia.org/wiki/Rete_algorithm)
- [Truth Maintenance Systems](https://en.wikipedia.org/wiki/Truth_maintenance_system)
- [Effective Go](https://go.dev/doc/effective_go)

### Standards Projet
- `.github/prompts/common.md` - Standards généraux
- `.github/prompts/review.md` - Standards de revue

---

## 🚀 Prochaines Étapes

1. **Prioriser**: Décider de l'ordre d'implémentation (recommandation: Insert → Update → Retract)
2. **Planifier**: Allouer les ressources et le temps nécessaires
3. **Implémenter**: Suivre le plan et la checklist
4. **Valider**: Tester exhaustivement chaque action
5. **Documenter**: Mettre à jour toute la documentation
6. **Release**: Publier v1.1.0 avec les actions complètes

---

## ✍️ Signature

**Créé par**: resinsec (GitHub Copilot CLI)  
**Date**: 2025-12-17  
**Version**: 1.0  

**Statut**: 📋 TODO - Actions à implémenter pour v1.1.0

---

*Ce document sera mis à jour au fur et à mesure de l'implémentation des actions.*
