# Prompt 06 - Implémentation du module xuples core

## 🎯 Objectif

Implémenter le module `tsd/xuples` avec les structures de données, xuple-spaces, et politiques conçues précédemment.

Ce module doit :
- Gérer les xuples (fait principal + faits déclencheurs)
- Gérer les xuple-spaces avec leurs politiques
- Implémenter les politiques de sélection, consommation et rétention
- Fournir une API claire pour l'action Xuple
- Être totalement découplé du moteur RETE

## 📋 Tâches

### 1. Créer la structure du package xuples

**Objectif** : Mettre en place l'organisation du package.

- [ ] Créer le répertoire `tsd/xuples/`
- [ ] Définir l'organisation des fichiers selon la conception
- [ ] Créer les fichiers de base avec copyright

**Structure attendue** :
```
tsd/xuples/
├── xuples.go              # Types publics, XupleManager
├── xuplespace.go          # Implémentation XupleSpace
├── policies.go            # Types de politiques
├── policy_selection.go    # Implémentations sélection
├── policy_consumption.go  # Implémentations consommation
├── policy_retention.go    # Implémentations rétention
├── lifecycle.go           # Gestion cycle de vie
├── errors.go              # Erreurs spécifiques
├── doc.go                 # Documentation du package
├── xuples_test.go
├── xuplespace_test.go
├── policies_test.go
└── testdata/
    └── examples.tsd
```

**Livrables** :
- [ ] Arborescence créée
- [ ] Fichiers de base avec copyright
- [ ] doc.go avec documentation du package

### 2. Définir les types et structures de données core

**Objectif** : Implémenter les structures fondamentales du module.

**Fichier à créer** : `tsd/xuples/xuples.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package xuples implémente le système de xuple-spaces pour TSD.
//
// Un xuple est un tuple étendu contenant un fait principal et les faits
// déclencheurs qui ont conduit à sa création. Les xuples sont stockés dans
// des xuple-spaces qui appliquent des politiques de sélection, consommation
// et rétention.
//
// Ce package est totalement découplé du moteur RETE et peut être utilisé
// indépendamment.
package xuples

import (
    "time"
    
    "tsd/rete"
)

// Xuple représente un tuple étendu avec fait principal et faits déclencheurs
type Xuple struct {
    ID              string        // Identifiant unique
    Fact            *rete.Fact    // Fait principal
    TriggeringFacts []*rete.Fact  // Faits qui ont déclenché la création
    CreatedAt       time.Time     // Date de création
    Metadata        XupleMetadata // Métadonnées
}

// XupleMetadata contient les métadonnées d'un xuple
type XupleMetadata struct {
    ConsumptionCount int                    // Nombre de consommations
    ConsumedBy       map[string]time.Time   // Agents ayant consommé (agent -> timestamp)
    ExpiresAt        time.Time              // Date d'expiration (0 si illimitée)
    State            XupleState             // État actuel
}

// XupleState représente l'état d'un xuple
type XupleState int

const (
    // XupleStateAvailable indique que le xuple est disponible
    XupleStateAvailable XupleState = iota
    
    // XupleStateConsumed indique que le xuple a été consommé (pour once)
    XupleStateConsumed
    
    // XupleStateExpired indique que le xuple a expiré
    XupleStateExpired
)

// String retourne la représentation textuelle de l'état
func (s XupleState) String() string {
    switch s {
    case XupleStateAvailable:
        return "available"
    case XupleStateConsumed:
        return "consumed"
    case XupleStateExpired:
        return "expired"
    default:
        return "unknown"
    }
}

// IsAvailable retourne true si le xuple est disponible pour consommation
func (x *Xuple) IsAvailable() bool {
    return x.Metadata.State == XupleStateAvailable
}

// IsExpired vérifie si le xuple a expiré
func (x *Xuple) IsExpired() bool {
    if x.Metadata.State == XupleStateExpired {
        return true
    }
    
    if !x.Metadata.ExpiresAt.IsZero() && time.Now().After(x.Metadata.ExpiresAt) {
        x.Metadata.State = XupleStateExpired
        return true
    }
    
    return false
}

// CanBeConsumedBy vérifie si un agent peut consommer ce xuple
func (x *Xuple) CanBeConsumedBy(agentID string, policy ConsumptionPolicy) bool {
    if !x.IsAvailable() || x.IsExpired() {
        return false
    }
    
    return policy.CanConsume(x, agentID)
}

// MarkConsumedBy marque le xuple comme consommé par un agent
func (x *Xuple) MarkConsumedBy(agentID string) {
    if x.Metadata.ConsumedBy == nil {
        x.Metadata.ConsumedBy = make(map[string]time.Time)
    }
    
    x.Metadata.ConsumedBy[agentID] = time.Now()
    x.Metadata.ConsumptionCount++
}

// XupleManager gère les xuple-spaces
type XupleManager interface {
    // CreateXupleSpace crée un nouveau xuple-space avec les politiques données
    CreateXupleSpace(name string, config XupleSpaceConfig) error
    
    // GetXupleSpace retourne un xuple-space par son nom
    GetXupleSpace(name string) (XupleSpace, error)
    
    // CreateXuple crée un xuple dans le xuple-space spécifié
    CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error
    
    // ListXupleSpaces retourne la liste des noms de xuple-spaces
    ListXupleSpaces() []string
    
    // Close ferme le manager et nettoie les ressources
    Close() error
}

// XupleSpace représente un espace de xuples
type XupleSpace interface {
    // Name retourne le nom du xuple-space
    Name() string
    
    // Insert insère un xuple dans le xuple-space
    Insert(xuple *Xuple) error
    
    // Retrieve récupère un xuple pour un agent selon les politiques
    Retrieve(agentID string) (*Xuple, error)
    
    // MarkConsumed marque un xuple comme consommé par un agent
    MarkConsumed(xupleID string, agentID string) error
    
    // Count retourne le nombre de xuples disponibles
    Count() int
    
    // Cleanup nettoie les xuples expirés
    Cleanup() int
    
    // GetConfig retourne la configuration du xuple-space
    GetConfig() XupleSpaceConfig
}

// XupleSpaceConfig configure un xuple-space
type XupleSpaceConfig struct {
    Name              string
    SelectionPolicy   SelectionPolicy
    ConsumptionPolicy ConsumptionPolicy
    RetentionPolicy   RetentionPolicy
}
```

**Livrables** :
- [ ] Fichier xuples.go créé avec copyright
- [ ] Structures Xuple et XupleMetadata complètes
- [ ] Interfaces XupleManager et XupleSpace définies
- [ ] Méthodes de l'état du xuple implémentées
- [ ] Documentation GoDoc complète

### 3. Définir les interfaces de politiques

**Objectif** : Créer les interfaces pour les trois types de politiques.

**Fichier à créer** : `tsd/xuples/policies.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import "time"

// SelectionPolicy définit comment sélectionner un xuple parmi plusieurs
type SelectionPolicy interface {
    // Select sélectionne un xuple parmi une liste de xuples disponibles
    // Retourne nil si aucun xuple n'est disponible
    Select(xuples []*Xuple) *Xuple
    
    // Name retourne le nom de la politique
    Name() string
}

// ConsumptionPolicy définit comment les xuples peuvent être consommés
type ConsumptionPolicy interface {
    // CanConsume vérifie si un agent peut consommer un xuple
    CanConsume(xuple *Xuple, agentID string) bool
    
    // OnConsumed est appelé après qu'un xuple ait été consommé
    // Retourne true si le xuple doit être marqué comme complètement consommé
    OnConsumed(xuple *Xuple, agentID string) bool
    
    // Name retourne le nom de la politique
    Name() string
}

// RetentionPolicy définit la durée de vie des xuples
type RetentionPolicy interface {
    // ComputeExpiration calcule la date d'expiration pour un nouveau xuple
    // Retourne zero time si pas d'expiration
    ComputeExpiration(createdAt time.Time) time.Time
    
    // ShouldRetain vérifie si un xuple doit être conservé
    ShouldRetain(xuple *Xuple) bool
    
    // Name retourne le nom de la politique
    Name() string
}

// PolicyType représente le type de politique
type PolicyType int

const (
    PolicyTypeSelection PolicyType = iota
    PolicyTypeConsumption
    PolicyTypeRetention
)

// String retourne la représentation textuelle du type
func (p PolicyType) String() string {
    switch p {
    case PolicyTypeSelection:
        return "selection"
    case PolicyTypeConsumption:
        return "consumption"
    case PolicyTypeRetention:
        return "retention"
    default:
        return "unknown"
    }
}
```

**Livrables** :
- [ ] Interfaces des trois types de politiques
- [ ] Documentation claire des contrats
- [ ] Types et constantes nécessaires

### 4. Implémenter les politiques de sélection

**Objectif** : Créer les implémentations des politiques de sélection.

**Fichier à créer** : `tsd/xuples/policy_selection.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
    "math/rand"
    "time"
)

// RandomSelectionPolicy sélectionne aléatoirement
type RandomSelectionPolicy struct {
    rng *rand.Rand
}

// NewRandomSelectionPolicy crée une nouvelle politique de sélection aléatoire
func NewRandomSelectionPolicy() *RandomSelectionPolicy {
    return &RandomSelectionPolicy{
        rng: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (p *RandomSelectionPolicy) Select(xuples []*Xuple) *Xuple {
    if len(xuples) == 0 {
        return nil
    }
    return xuples[p.rng.Intn(len(xuples))]
}

func (p *RandomSelectionPolicy) Name() string {
    return "random"
}

// FIFOSelectionPolicy sélectionne le premier entré (plus ancien)
type FIFOSelectionPolicy struct{}

// NewFIFOSelectionPolicy crée une nouvelle politique FIFO
func NewFIFOSelectionPolicy() *FIFOSelectionPolicy {
    return &FIFOSelectionPolicy{}
}

func (p *FIFOSelectionPolicy) Select(xuples []*Xuple) *Xuple {
    if len(xuples) == 0 {
        return nil
    }
    
    // Trouver le xuple le plus ancien
    oldest := xuples[0]
    for _, xuple := range xuples[1:] {
        if xuple.CreatedAt.Before(oldest.CreatedAt) {
            oldest = xuple
        }
    }
    
    return oldest
}

func (p *FIFOSelectionPolicy) Name() string {
    return "fifo"
}

// LIFOSelectionPolicy sélectionne le dernier entré (plus récent)
type LIFOSelectionPolicy struct{}

// NewLIFOSelectionPolicy crée une nouvelle politique LIFO
func NewLIFOSelectionPolicy() *LIFOSelectionPolicy {
    return &LIFOSelectionPolicy{}
}

func (p *LIFOSelectionPolicy) Select(xuples []*Xuple) *Xuple {
    if len(xuples) == 0 {
        return nil
    }
    
    // Trouver le xuple le plus récent
    newest := xuples[0]
    for _, xuple := range xuples[1:] {
        if xuple.CreatedAt.After(newest.CreatedAt) {
            newest = xuple
        }
    }
    
    return newest
}

func (p *LIFOSelectionPolicy) Name() string {
    return "lifo"
}
```

**Livrables** :
- [ ] RandomSelectionPolicy implémenté
- [ ] FIFOSelectionPolicy implémenté
- [ ] LIFOSelectionPolicy implémenté
- [ ] Documentation GoDoc pour chaque politique
- [ ] Tests unitaires de chaque politique

### 5. Implémenter les politiques de consommation

**Objectif** : Créer les implémentations des politiques de consommation.

**Fichier à créer** : `tsd/xuples/policy_consumption.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

// OnceConsumptionPolicy permet une seule consommation au total
type OnceConsumptionPolicy struct{}

// NewOnceConsumptionPolicy crée une nouvelle politique de consommation unique
func NewOnceConsumptionPolicy() *OnceConsumptionPolicy {
    return &OnceConsumptionPolicy{}
}

func (p *OnceConsumptionPolicy) CanConsume(xuple *Xuple, agentID string) bool {
    return xuple.Metadata.ConsumptionCount == 0
}

func (p *OnceConsumptionPolicy) OnConsumed(xuple *Xuple, agentID string) bool {
    // Marquer comme complètement consommé après la première consommation
    return true
}

func (p *OnceConsumptionPolicy) Name() string {
    return "once"
}

// PerAgentConsumptionPolicy permet une consommation par agent
type PerAgentConsumptionPolicy struct{}

// NewPerAgentConsumptionPolicy crée une nouvelle politique par agent
func NewPerAgentConsumptionPolicy() *PerAgentConsumptionPolicy {
    return &PerAgentConsumptionPolicy{}
}

func (p *PerAgentConsumptionPolicy) CanConsume(xuple *Xuple, agentID string) bool {
    if xuple.Metadata.ConsumedBy == nil {
        return true
    }
    _, alreadyConsumed := xuple.Metadata.ConsumedBy[agentID]
    return !alreadyConsumed
}

func (p *PerAgentConsumptionPolicy) OnConsumed(xuple *Xuple, agentID string) bool {
    // Ne jamais marquer comme complètement consommé (autres agents peuvent consommer)
    return false
}

func (p *PerAgentConsumptionPolicy) Name() string {
    return "per-agent"
}

// LimitedConsumptionPolicy permet un nombre limité de consommations
type LimitedConsumptionPolicy struct {
    MaxConsumptions int
}

// NewLimitedConsumptionPolicy crée une nouvelle politique avec limite
func NewLimitedConsumptionPolicy(maxConsumptions int) *LimitedConsumptionPolicy {
    if maxConsumptions <= 0 {
        maxConsumptions = 1
    }
    return &LimitedConsumptionPolicy{
        MaxConsumptions: maxConsumptions,
    }
}

func (p *LimitedConsumptionPolicy) CanConsume(xuple *Xuple, agentID string) bool {
    return xuple.Metadata.ConsumptionCount < p.MaxConsumptions
}

func (p *LimitedConsumptionPolicy) OnConsumed(xuple *Xuple, agentID string) bool {
    // Marquer comme consommé si la limite est atteinte
    return xuple.Metadata.ConsumptionCount >= p.MaxConsumptions
}

func (p *LimitedConsumptionPolicy) Name() string {
    return "limited"
}
```

**Livrables** :
- [ ] OnceConsumptionPolicy implémenté
- [ ] PerAgentConsumptionPolicy implémenté
- [ ] LimitedConsumptionPolicy implémenté
- [ ] Documentation GoDoc pour chaque politique
- [ ] Tests unitaires de chaque politique

### 6. Implémenter les politiques de rétention

**Objectif** : Créer les implémentations des politiques de rétention.

**Fichier à créer** : `tsd/xuples/policy_retention.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import "time"

// UnlimitedRetentionPolicy conserve les xuples indéfiniment
type UnlimitedRetentionPolicy struct{}

// NewUnlimitedRetentionPolicy crée une nouvelle politique illimitée
func NewUnlimitedRetentionPolicy() *UnlimitedRetentionPolicy {
    return &UnlimitedRetentionPolicy{}
}

func (p *UnlimitedRetentionPolicy) ComputeExpiration(createdAt time.Time) time.Time {
    return time.Time{} // Zero time = pas d'expiration
}

func (p *UnlimitedRetentionPolicy) ShouldRetain(xuple *Xuple) bool {
    return true
}

func (p *UnlimitedRetentionPolicy) Name() string {
    return "unlimited"
}

// DurationRetentionPolicy expire les xuples après une durée
type DurationRetentionPolicy struct {
    Duration time.Duration
}

// NewDurationRetentionPolicy crée une nouvelle politique basée sur la durée
func NewDurationRetentionPolicy(duration time.Duration) *DurationRetentionPolicy {
    if duration <= 0 {
        duration = 1 * time.Hour // Défaut sécurisé
    }
    return &DurationRetentionPolicy{
        Duration: duration,
    }
}

func (p *DurationRetentionPolicy) ComputeExpiration(createdAt time.Time) time.Time {
    return createdAt.Add(p.Duration)
}

func (p *DurationRetentionPolicy) ShouldRetain(xuple *Xuple) bool {
    if xuple.Metadata.ExpiresAt.IsZero() {
        return true
    }
    return time.Now().Before(xuple.Metadata.ExpiresAt)
}

func (p *DurationRetentionPolicy) Name() string {
    return "duration"
}
```

**Livrables** :
- [ ] UnlimitedRetentionPolicy implémenté
- [ ] DurationRetentionPolicy implémenté
- [ ] Documentation GoDoc pour chaque politique
- [ ] Tests unitaires de chaque politique

### 7. Implémenter XupleSpace

**Objectif** : Créer l'implémentation du xuple-space.

**Fichier à créer** : `tsd/xuples/xuplespace.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
    "fmt"
    "sync"
    
    "github.com/google/uuid"
)

// DefaultXupleSpace implémente XupleSpace
type DefaultXupleSpace struct {
    name   string
    config XupleSpaceConfig
    xuples map[string]*Xuple // xupleID -> Xuple
    mu     sync.RWMutex
}

// NewXupleSpace crée un nouveau xuple-space
func NewXupleSpace(config XupleSpaceConfig) *DefaultXupleSpace {
    return &DefaultXupleSpace{
        name:   config.Name,
        config: config,
        xuples: make(map[string]*Xuple),
    }
}

func (xs *DefaultXupleSpace) Name() string {
    return xs.name
}

func (xs *DefaultXupleSpace) Insert(xuple *Xuple) error {
    if xuple == nil {
        return ErrNilXuple
    }
    
    xs.mu.Lock()
    defer xs.mu.Unlock()
    
    // Générer un ID si nécessaire
    if xuple.ID == "" {
        xuple.ID = uuid.New().String()
    }
    
    // Appliquer la politique de rétention
    xuple.Metadata.ExpiresAt = xs.config.RetentionPolicy.ComputeExpiration(xuple.CreatedAt)
    
    xs.xuples[xuple.ID] = xuple
    return nil
}

func (xs *DefaultXupleSpace) Retrieve(agentID string) (*Xuple, error) {
    xs.mu.RLock()
    defer xs.mu.RUnlock()
    
    // Collecter les xuples disponibles pour cet agent
    available := make([]*Xuple, 0)
    for _, xuple := range xs.xuples {
        if xuple.CanBeConsumedBy(agentID, xs.config.ConsumptionPolicy) {
            available = append(available, xuple)
        }
    }
    
    if len(available) == 0 {
        return nil, ErrNoAvailableXuple
    }
    
    // Sélectionner selon la politique
    selected := xs.config.SelectionPolicy.Select(available)
    if selected == nil {
        return nil, ErrNoAvailableXuple
    }
    
    return selected, nil
}

func (xs *DefaultXupleSpace) MarkConsumed(xupleID string, agentID string) error {
    xs.mu.Lock()
    defer xs.mu.Unlock()
    
    xuple, exists := xs.xuples[xupleID]
    if !exists {
        return ErrXupleNotFound
    }
    
    if !xuple.CanBeConsumedBy(agentID, xs.config.ConsumptionPolicy) {
        return ErrXupleNotAvailable
    }
    
    // Marquer comme consommé
    xuple.MarkConsumedBy(agentID)
    
    // Vérifier si le xuple doit être marqué comme complètement consommé
    if xs.config.ConsumptionPolicy.OnConsumed(xuple, agentID) {
        xuple.Metadata.State = XupleStateConsumed
    }
    
    return nil
}

func (xs *DefaultXupleSpace) Count() int {
    xs.mu.RLock()
    defer xs.mu.RUnlock()
    
    count := 0
    for _, xuple := range xs.xuples {
        if xuple.IsAvailable() && !xuple.IsExpired() {
            count++
        }
    }
    
    return count
}

func (xs *DefaultXupleSpace) Cleanup() int {
    xs.mu.Lock()
    defer xs.mu.Unlock()
    
    cleaned := 0
    for id, xuple := range xs.xuples {
        if !xs.config.RetentionPolicy.ShouldRetain(xuple) || xuple.IsExpired() {
            delete(xs.xuples, id)
            cleaned++
        }
    }
    
    return cleaned
}

func (xs *DefaultXupleSpace) GetConfig() XupleSpaceConfig {
    return xs.config
}
```

**Livrables** :
- [ ] DefaultXupleSpace implémenté avec copyright
- [ ] Thread-safe (sync.RWMutex)
- [ ] Toutes les méthodes de l'interface implémentées
- [ ] Génération d'ID pour les xuples
- [ ] Application des politiques
- [ ] Gestion d'erreurs robuste
- [ ] Documentation GoDoc complète

### 8. Implémenter XupleManager

**Objectif** : Créer l'implémentation du manager de xuple-spaces.

**Ajout dans** : `tsd/xuples/xuples.go`

**Code attendu** :
```go
// DefaultXupleManager implémente XupleManager
type DefaultXupleManager struct {
    spaces map[string]XupleSpace
    mu     sync.RWMutex
}

// NewXupleManager crée un nouveau manager de xuple-spaces
func NewXupleManager() *DefaultXupleManager {
    return &DefaultXupleManager{
        spaces: make(map[string]XupleSpace),
    }
}

func (m *DefaultXupleManager) CreateXupleSpace(name string, config XupleSpaceConfig) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if _, exists := m.spaces[name]; exists {
        return ErrXupleSpaceExists
    }
    
    config.Name = name
    xs := NewXupleSpace(config)
    m.spaces[name] = xs
    
    return nil
}

func (m *DefaultXupleManager) GetXupleSpace(name string) (XupleSpace, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    xs, exists := m.spaces[name]
    if !exists {
        return nil, ErrXupleSpaceNotFound
    }
    
    return xs, nil
}

func (m *DefaultXupleManager) CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
    xs, err := m.GetXupleSpace(xuplespace)
    if err != nil {
        return err
    }
    
    xuple := &Xuple{
        Fact:            fact,
        TriggeringFacts: triggeringFacts,
        CreatedAt:       time.Now(),
        Metadata: XupleMetadata{
            State:      XupleStateAvailable,
            ConsumedBy: make(map[string]time.Time),
        },
    }
    
    return xs.Insert(xuple)
}

func (m *DefaultXupleManager) ListXupleSpaces() []string {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    names := make([]string, 0, len(m.spaces))
    for name := range m.spaces {
        names = append(names, name)
    }
    
    return names
}

func (m *DefaultXupleManager) Close() error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    // Nettoyer tous les xuple-spaces
    for _, xs := range m.spaces {
        xs.Cleanup()
    }
    
    m.spaces = make(map[string]XupleSpace)
    return nil
}
```

**Livrables** :
- [ ] DefaultXupleManager implémenté
- [ ] Thread-safe
- [ ] Toutes les méthodes de l'interface implémentées
- [ ] Gestion d'erreurs
- [ ] Documentation GoDoc

### 9. Définir les erreurs spécifiques

**Objectif** : Créer des erreurs typed pour le module.

**Fichier à créer** : `tsd/xuples/errors.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import "errors"

// Erreurs du module xuples
var (
    // ErrNilXuple est retourné quand un xuple nil est fourni
    ErrNilXuple = errors.New("xuple cannot be nil")
    
    // ErrXupleNotFound est retourné quand un xuple n'existe pas
    ErrXupleNotFound = errors.New("xuple not found")
    
    // ErrXupleNotAvailable est retourné quand un xuple n'est pas disponible
    ErrXupleNotAvailable = errors.New("xuple not available for consumption")
    
    // ErrNoAvailableXuple est retourné quand aucun xuple n'est disponible
    ErrNoAvailableXuple = errors.New("no available xuple")
    
    // ErrXupleSpaceNotFound est retourné quand un xuple-space n'existe pas
    ErrXupleSpaceNotFound = errors.New("xuple-space not found")
    
    // ErrXupleSpaceExists est retourné lors d'une tentative de création d'un xuple-space existant
    ErrXupleSpaceExists = errors.New("xuple-space already exists")
    
    // ErrInvalidPolicy est retourné quand une politique est invalide
    ErrInvalidPolicy = errors.New("invalid policy")
    
    // ErrInvalidConfiguration est retourné quand une configuration est invalide
    ErrInvalidConfiguration = errors.New("invalid xuple-space configuration")
)
```

**Livrables** :
- [ ] Erreurs typed définies
- [ ] Documentation de chaque erreur
- [ ] Conventions Go respectées (errors.New)

### 10. Créer les tests complets du module

**Objectif** : Tester exhaustivement toutes les fonctionnalités.

**Fichiers à créer** :
- `tsd/xuples/xuples_test.go`
- `tsd/xuples/xuplespace_test.go`
- `tsd/xuples/policies_test.go`

**Tests attendus** (exemples) :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
    "testing"
    "time"
    
    "tsd/rete"
)

func TestXuple_Lifecycle(t *testing.T) {
    t.Log("🧪 TEST CYCLE DE VIE XUPLE")
    
    fact := &rete.Fact{ID: "f1", Type: "Person"}
    triggering := []*rete.Fact{
        {ID: "t1", Type: "Event"},
        {ID: "t2", Type: "Condition"},
    }
    
    xuple := &Xuple{
        ID:              "x1",
        Fact:            fact,
        TriggeringFacts: triggering,
        CreatedAt:       time.Now(),
        Metadata: XupleMetadata{
            State:      XupleStateAvailable,
            ConsumedBy: make(map[string]time.Time),
        },
    }
    
    // Vérifier l'état initial
    if !xuple.IsAvailable() {
        t.Error("❌ Xuple devrait être disponible")
    }
    
    if xuple.IsExpired() {
        t.Error("❌ Xuple ne devrait pas être expiré")
    }
    
    t.Log("✅ Cycle de vie de base correct")
}

func TestXupleSpace_InsertAndRetrieve(t *testing.T) {
    t.Log("🧪 TEST INSERT/RETRIEVE XUPLE-SPACE")
    
    config := XupleSpaceConfig{
        Name:              "test",
        SelectionPolicy:   NewFIFOSelectionPolicy(),
        ConsumptionPolicy: NewOnceConsumptionPolicy(),
        RetentionPolicy:   NewUnlimitedRetentionPolicy(),
    }
    
    xs := NewXupleSpace(config)
    
    // Insérer un xuple
    xuple := &Xuple{
        Fact:      &rete.Fact{ID: "f1"},
        CreatedAt: time.Now(),
        Metadata: XupleMetadata{
            State:      XupleStateAvailable,
            ConsumedBy: make(map[string]time.Time),
        },
    }
    
    err := xs.Insert(xuple)
    if err != nil {
        t.Fatalf("❌ Erreur insertion: %v", err)
    }
    
    // Récupérer le xuple
    retrieved, err := xs.Retrieve("agent1")
    if err != nil {
        t.Fatalf("❌ Erreur récupération: %v", err)
    }
    
    if retrieved == nil {
        t.Fatal("❌ Xuple récupéré est nil")
    }
    
    if retrieved.Fact.ID != "f1" {
        t.Errorf("❌ Mauvais xuple récupéré")
    }
    
    t.Log("✅ Insert/Retrieve fonctionne")
}

func TestXupleManager_CreateAndManage(t *testing.T) {
    t.Log("🧪 TEST XUPLE MANAGER")
    
    manager := NewXupleManager()
    
    // Créer un xuple-space
    config := XupleSpaceConfig{
        SelectionPolicy:   NewFIFOSelectionPolicy(),
        ConsumptionPolicy: NewOnceConsumptionPolicy(),
        RetentionPolicy:   NewUnlimitedRetentionPolicy(),
    }
    
    err := manager.CreateXupleSpace("myspace", config)
    if err != nil {
        t.Fatalf("❌ Erreur création xuple-space: %v", err)
    }
    
    // Vérifier que le xuple-space existe
    xs, err := manager.GetXupleSpace("myspace")
    if err != nil {
        t.Fatalf("❌ Erreur récupération xuple-space: %v", err)
    }
    
    if xs.Name() != "myspace" {
        t.Errorf("❌ Mauvais nom: %s", xs.Name())
    }
    
    // Créer un xuple
    fact := &rete.Fact{ID: "f1"}
    triggering := []*rete.Fact{{ID: "t1"}}
    
    err = manager.CreateXuple("myspace", fact, triggering)
    if err != nil {
        t.Fatalf("❌ Erreur création xuple: %v", err)
    }
    
    // Vérifier le compte
    if xs.Count() != 1 {
        t.Errorf("❌ Attendu 1 xuple, reçu %d", xs.Count())
    }
    
    t.Log("✅ XupleManager fonctionne")
}

// Plus de tests pour chaque politique...
```

**Livrables** :
- [ ] Tests de Xuple et métadonnées
- [ ] Tests de XupleSpace (toutes méthodes)
- [ ] Tests de XupleManager
- [ ] Tests de chaque politique
- [ ] Tests de concurrence
- [ ] Tests d'erreurs
- [ ] Couverture > 80%
- [ ] Tous les tests passent

## 📁 Structure finale

```
tsd/xuples/
├── doc.go
├── xuples.go              # Types core, interfaces, manager
├── xuplespace.go          # Implémentation xuple-space
├── policies.go            # Interfaces politiques
├── policy_selection.go    # Politiques sélection
├── policy_consumption.go  # Politiques consommation
├── policy_retention.go    # Politiques rétention
├── errors.go              # Erreurs
├── xuples_test.go
├── xuplespace_test.go
├── policies_test.go
└── testdata/
```

## ✅ Critères de succès

- [ ] Package xuples complet et fonctionnel
- [ ] Toutes les structures implémentées avec copyright
- [ ] Toutes les interfaces implémentées
- [ ] Toutes les politiques implémentées
- [ ] Thread-safe (sync.RWMutex)
- [ ] Aucun hardcoding
- [ ] Gestion d'erreurs robuste
- [ ] Documentation GoDoc complète
- [ ] Tests complets avec couverture > 80%
- [ ] Tous les tests passent
- [ ] `make test-unit` passe
- [ ] Totalement découplé de RETE

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/docs/xuples/design/` - Conception détaillée
- Effective Go - https://go.dev/doc/effective_go
- Go Concurrency Patterns

## 🎯 Prochaine étape

Une fois le module xuples core implémenté, passer au prompt **07-integrate-xuple-action.md** pour intégrer l'action Xuple avec le module xuples et permettre la création de xuples depuis les règles.