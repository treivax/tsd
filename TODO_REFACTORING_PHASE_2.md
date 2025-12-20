# TODO - Suite du Refactoring API TSD

**Date de création** : 2025-12-19  
**Contexte** : Post-refactoring tsdio (cachage de `_id_`)  
**Statut actuel** : ✅ Phase 1 terminée

---

## ✅ Terminé - Phase 1

- [x] Cacher `_id_` de `tsdio.Fact` dans l'API JSON
- [x] Ajouter méthodes `GetInternalID()` et `SetInternalID()`
- [x] Mettre à jour tous les usages dans le code
- [x] Tests à 100% de couverture pour tsdio
- [x] Documentation complète de l'API
- [x] Validation des changements

---

## 🔜 À Implémenter - Phase 2 (Prompt 07)

### 1. Support des Affectations de Variables

**Priorité** : Haute  
**Complexité** : Moyenne

#### Structures à Créer

```go
// tsdio/program.go (nouveau fichier)

type FactAssignment struct {
    Variable string `json:"variable"`
    Fact     *Fact  `json:"fact"`
}

type Program struct {
    Types           []TypeDefinition   `json:"types,omitempty"`
    Actions         []ActionDefinition `json:"actions,omitempty"`
    FactAssignments []FactAssignment   `json:"factAssignments,omitempty"`
    Facts           []*Fact            `json:"facts,omitempty"`
    Rules           []Rule             `json:"rules,omitempty"`
}
```

#### Méthodes à Implémenter

- [ ] `NewFactAssignment(variable string, fact *Fact) (*FactAssignment, error)`
- [ ] `(fa *FactAssignment) Validate() error`
- [ ] `NewProgram() *Program`
- [ ] `(p *Program) AddFactAssignment(assignment *FactAssignment) error`
- [ ] `(p *Program) AddFact(fact *Fact) error`
- [ ] `(p *Program) Validate() error`

#### Tests à Créer

- [ ] `TestFactAssignment_Creation`
- [ ] `TestFactAssignment_Validate`
- [ ] `TestProgram_AddFactAssignment`
- [ ] `TestProgram_Validate`
- [ ] Tests d'intégration avec le moteur RETE

#### Impact sur le Code

**Fichiers à modifier** :
- `constraint/parser.go` - Support syntaxe `$var = fact`
- `constraint/program_state.go` - Gestion des variables
- `rete/network_manager.go` - Association variable → ID
- `internal/servercmd/servercmd.go` - Traitement des affectations

**Exemple de syntaxe TSD** :
```tsd
type User : <name: string, age: number>

// Affectation de variable
$alice = User <name: "Alice", age: 30>

// Utilisation de la variable
fact login : Login <
    user: $alice,
    email: "alice@example.com"
>
```

---

### 2. API Validator (Si Nécessaire)

**Priorité** : Basse  
**Complexité** : Faible

**Note** : La validation de `_id_` est déjà présente dans `constraint/parser.go`. Un validateur API dédié n'est nécessaire que si on veut valider d'autres aspects côté API.

#### Potentielles Validations API

- [ ] Validation format JSON
- [ ] Validation taille des requêtes
- [ ] Validation types de données
- [ ] Sanitization des entrées

Si implémenté, créer :
```go
// api/validator.go

type APIValidator struct {}

func NewAPIValidator() *APIValidator
func (v *APIValidator) ValidateFact(fact interface{}) error
func (v *APIValidator) ValidateFactAssignment(assignment interface{}) error
func (v *APIValidator) SanitizeFact(fact map[string]interface{}) map[string]interface{}
```

---

### 3. Modifications API Result (package api/)

**Priorité** : Moyenne  
**Complexité** : Faible

**Note** : Le package `api/` utilise actuellement `xuples.Xuple` pour les résultats, pas `tsdio.Fact`. 

#### Vérifications à Faire

- [ ] Vérifier si `xuples.Xuple` expose `_id_` en JSON
- [ ] Si oui, appliquer le même pattern de cachage
- [ ] Tests de sérialisation JSON pour les Xuples

#### Si Modification Nécessaire

```go
// xuples/xuple.go

type Xuple struct {
    internalID string  // Caché
    Type       string                 `json:"type"`
    Fields     map[string]interface{} `json:"fields"`
}

func (x *Xuple) GetInternalID() string
func (x *Xuple) SetInternalID(id string)
```

---

## 🐛 Bugs/Issues à Investiguer

### Tests Échouant dans constraint/

**Fichiers** :
- `constraint/aggregation_calculation_test.go`
- `constraint/arithmetic_expressions_test.go`

**Problèmes** :
- Tests d'agrégation (AVG, SUM, COUNT, MIN) échouent
- Nombre d'activations attendu != nombre réel
- Logs montrent que les actions s'exécutent mais activations = 0

**Actions** :
- [ ] Investiguer pourquoi les activations ne sont pas comptées
- [ ] Vérifier si c'est lié au système de tokens/bindings
- [ ] Corriger ou mettre à jour les tests

**Note** : Ces échecs sont **préexistants** et **non causés** par le refactoring tsdio.

---

## 📚 Documentation à Compléter

### Documentation Technique

- [ ] Guide d'utilisation des affectations de variables
- [ ] Exemples d'utilisation avancés
- [ ] Diagrammes d'architecture mis à jour

### Documentation API

- [ ] Swagger/OpenAPI spec pour les endpoints
- [ ] Guide de migration v1.1 → v1.2
- [ ] Changelog détaillé

---

## 🧪 Tests à Ajouter

### Tests d'Intégration

- [ ] Test flow complet : Affectation → Référence → Exécution
- [ ] Test avec multiples variables
- [ ] Test de résolution de variables en cascade
- [ ] Test d'erreurs (variable non définie, etc.)

### Tests E2E

- [ ] Test via API HTTP avec affectations
- [ ] Test CLI avec fichier TSD contenant affectations
- [ ] Test de performance avec nombreuses affectations

---

## 🔧 Améliorations Futures

### Performance

- [ ] Profiling des conversions RETE → tsdio
- [ ] Optimisation de la sérialisation JSON
- [ ] Cache pour les ID internes fréquemment accédés

### Fonctionnalités

- [ ] Support des variables dans les règles
- [ ] Support des variables dans les contraintes
- [ ] Scoping des variables (local/global)

### Sécurité

- [ ] Rate limiting pour l'API
- [ ] Validation stricte des types
- [ ] Audit logging des opérations API

---

## 📅 Planning Suggéré

### Sprint 1 (1-2 jours)
- [ ] Implémenter structures `FactAssignment` et `Program`
- [ ] Tests unitaires de base
- [ ] Documentation API

### Sprint 2 (2-3 jours)
- [ ] Modifications parser pour syntaxe `$var = fact`
- [ ] Intégration avec le moteur RETE
- [ ] Tests d'intégration

### Sprint 3 (1 jour)
- [ ] Corriger les tests échouants dans constraint/
- [ ] Tests E2E complets
- [ ] Documentation finale

### Sprint 4 (1 jour)
- [ ] Revue de code
- [ ] Validation complète
- [ ] Release v1.3.0

---

## 🚨 Risques Identifiés

### Risque 1 : Compatibilité Parser

**Impact** : Élevé  
**Probabilité** : Moyenne

La syntaxe `$var = fact` pourrait entrer en conflit avec la syntaxe existante.

**Mitigation** :
- Analyse complète de la grammaire
- Tests exhaustifs des cas limites
- Backward compatibility tests

### Risque 2 : Performance

**Impact** : Moyen  
**Probabilité** : Faible

L'association variable → ID pourrait impacter les performances.

**Mitigation** :
- Benchmarks avant/après
- Profiling des opérations critiques
- Optimisation si nécessaire

### Risque 3 : Complexité

**Impact** : Moyen  
**Probabilité** : Moyenne

L'ajout de variables augmente la complexité du système.

**Mitigation** :
- Documentation claire
- Tests exhaustifs
- Exemples simples

---

## 📞 Contacts & Ressources

### Documentation de Référence

- `06-prompt-api-tsdio.md` - Spécifications originales
- `RAPPORT_REFACTORING_TSDIO_API.md` - Rapport phase 1
- `tsdio/API_DOCUMENTATION.md` - Documentation API actuelle

### Points de Contact

- **Code Review** : À définir
- **Architecture** : À définir
- **Tests** : À définir

---

## ✅ Checklist avant de Commencer Phase 2

- [ ] Phase 1 validée et mergée
- [ ] Tests constraint/ corrigés ou documentés
- [ ] Design document pour FactAssignment approuvé
- [ ] Parser grammar analysée
- [ ] Équipe informée des changements

---

**Dernière mise à jour** : 2025-12-19  
**Auteur** : Assistant AI (resinsec)  
**Status** : 📋 TODO - En attente de début Phase 2
