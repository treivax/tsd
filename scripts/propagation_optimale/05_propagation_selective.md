# 🚀 Prompt 05 - Propagation Sélective

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Implémenter le moteur de propagation sélective : `DeltaPropagator` qui orchestre la propagation des changements uniquement vers les nœuds affectés par un `FactDelta`.

Cette propagation sélective est le cœur de l'optimisation RETE-II/TREAT : au lieu de propager vers tous les nœuds, on ne propage que vers ceux qui dépendent des champs modifiés.

**⚠️ IMPORTANT** : Ce prompt génère du code. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompt 01 validé** : Conception disponible
- [x] **Prompt 02 validé** : Modèle de données delta implémenté
- [x] **Prompt 03 validé** : Indexation des dépendances implémentée
- [x] **Prompt 04 validé** : Détection delta implémentée
- [x] **Tests passent** : `go test ./rete/delta/... -v` (100% success)
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - `rete/delta/field_delta.go` - Structures delta
  - `rete/delta/dependency_index.go` - Index de dépendances
  - `rete/delta/delta_detector.go` - Détecteur de changements

---

## 📂 Fichiers à Créer

Ajouter au package `rete/delta` :

```
rete/delta/
├── delta_propagator.go           # Propagateur principal
├── delta_propagator_test.go      # Tests unitaires
├── propagation_config.go         # Configuration propagation
├── propagation_config_test.go    # Tests configuration
├── propagation_strategy.go       # Stratégies de propagation
├── propagation_strategy_test.go  # Tests stratégies
├── propagation_metrics.go        # Métriques de propagation
└── propagation_benchmark_test.go # Benchmarks performance
```

---

## 🔧 Tâche 1 : Configuration de la Propagation

### Fichier : `rete/delta/propagation_config.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"time"
)

// PropagationMode définit le mode de propagation à utiliser.
type PropagationMode int

const (
	// PropagationModeDelta utilise la propagation sélective par delta
	PropagationModeDelta PropagationMode = iota
	
	// PropagationModeClassic utilise Retract+Insert classique (fallback)
	PropagationModeClassic
	
	// PropagationModeAuto choisit automatiquement selon le contexte
	PropagationModeAuto
)

// String retourne la représentation string du PropagationMode
func (pm PropagationMode) String() string {
	switch pm {
	case PropagationModeDelta:
		return "Delta"
	case PropagationModeClassic:
		return "Classic"
	case PropagationModeAuto:
		return "Auto"
	default:
		return "Unknown"
	}
}

// PropagationConfig contient la configuration du DeltaPropagator.
//
// Cette configuration permet de contrôler le comportement de la propagation
// sélective et les critères de fallback vers la propagation classique.
type PropagationConfig struct {
	// Mode de propagation par défaut
	DefaultMode PropagationMode
	
	// EnableDeltaPropagation active/désactive la propagation delta
	// (master switch pour activation/désactivation globale)
	EnableDeltaPropagation bool
	
	// DeltaThreshold est le seuil de ratio de changement au-delà duquel
	// on bascule en mode classique (Retract+Insert).
	// Valeur entre 0.0 et 1.0.
	// Exemple : 0.3 → si > 30% des champs changent, utiliser mode classique
	// Default: 0.5
	DeltaThreshold float64
	
	// MinFieldsForDelta est le nombre minimum de champs dans un fait
	// pour que la propagation delta soit utilisée.
	// Si le fait a moins de champs, utiliser mode classique.
	// Default: 3
	MinFieldsForDelta int
	
	// MaxAffectedNodesForDelta est le nombre maximum de nœuds affectés
	// au-delà duquel on bascule en mode classique.
	// Rationale : si trop de nœuds affectés, overhead delta > bénéfice
	// Default: 100
	MaxAffectedNodesForDelta int
	
	// AllowPrimaryKeyChange indique si les modifications de clé primaire
	// sont autorisées en mode delta.
	// Si false, tout changement de PK force le mode classique.
	// Default: false (car changement PK = changement d'ID interne)
	AllowPrimaryKeyChange bool
	
	// PrimaryKeyFields liste les noms de champs considérés comme clés primaires.
	// Si vide, détection automatique depuis les TypeDefinitions.
	// Default: []
	PrimaryKeyFields []string
	
	// EnableMetrics active la collecte de métriques de propagation
	// Default: true
	EnableMetrics bool
	
	// PropagationTimeout est le timeout maximum pour une propagation delta.
	// Si dépassé, la propagation est annulée (protection deadlock).
	// Default: 30 secondes
	PropagationTimeout time.Duration
	
	// RetryOnError indique si une propagation delta échouée doit être
	// retentée en mode classique (fallback automatique).
	// Default: true
	RetryOnError bool
	
	// MaxConcurrentPropagations est le nombre maximum de propagations
	// delta simultanées autorisées (contrôle charge).
	// Default: 10
	MaxConcurrentPropagations int
	
	// EnableOptimisticPropagation active la propagation optimiste :
	// ne pas attendre la fin de la propagation pour retourner.
	// Default: false (attente synchrone)
	EnableOptimisticPropagation bool
	
	// LogPropagationDetails active le logging détaillé de chaque propagation
	// (utile pour debugging, overhead en production).
	// Default: false
	LogPropagationDetails bool
}

// DefaultPropagationConfig retourne une configuration par défaut.
func DefaultPropagationConfig() PropagationConfig {
	return PropagationConfig{
		DefaultMode:                   PropagationModeAuto,
		EnableDeltaPropagation:        true,
		DeltaThreshold:                0.5,
		MinFieldsForDelta:             3,
		MaxAffectedNodesForDelta:      100,
		AllowPrimaryKeyChange:         false,
		PrimaryKeyFields:              []string{},
		EnableMetrics:                 true,
		PropagationTimeout:            30 * time.Second,
		RetryOnError:                  true,
		MaxConcurrentPropagations:     10,
		EnableOptimisticPropagation:   false,
		LogPropagationDetails:         false,
	}
}

// Validate vérifie que la configuration est valide.
//
// Retourne une erreur si des paramètres sont incohérents.
func (pc *PropagationConfig) Validate() error {
	if pc.DeltaThreshold < 0.0 || pc.DeltaThreshold > 1.0 {
		return fmt.Errorf("DeltaThreshold must be between 0.0 and 1.0, got %v", pc.DeltaThreshold)
	}
	
	if pc.MinFieldsForDelta < 0 {
		return fmt.Errorf("MinFieldsForDelta must be >= 0, got %d", pc.MinFieldsForDelta)
	}
	
	if pc.MaxAffectedNodesForDelta < 1 {
		return fmt.Errorf("MaxAffectedNodesForDelta must be >= 1, got %d", pc.MaxAffectedNodesForDelta)
	}
	
	if pc.PropagationTimeout < 0 {
		return fmt.Errorf("PropagationTimeout must be >= 0, got %v", pc.PropagationTimeout)
	}
	
	if pc.MaxConcurrentPropagations < 1 {
		return fmt.Errorf("MaxConcurrentPropagations must be >= 1, got %d", pc.MaxConcurrentPropagations)
	}
	
	return nil
}

// ShouldUseDelta détermine si la propagation delta doit être utilisée
// pour un FactDelta donné.
//
// Cette méthode applique les heuristiques configurées pour décider
// du mode de propagation optimal.
//
// Paramètres :
//   - delta : le FactDelta à propager
//   - affectedNodesCount : nombre de nœuds qui seraient affectés
//
// Retourne true si la propagation delta doit être utilisée.
func (pc *PropagationConfig) ShouldUseDelta(delta *FactDelta, affectedNodesCount int) bool {
	// 1. Feature flag global
	if !pc.EnableDeltaPropagation {
		return false
	}
	
	// 2. Vérifier mode forcé
	if pc.DefaultMode == PropagationModeClassic {
		return false
	}
	if pc.DefaultMode == PropagationModeDelta {
		return true
	}
	
	// 3. Mode Auto : appliquer heuristiques
	
	// Heuristique 1 : Nombre de champs
	if delta.FieldCount < pc.MinFieldsForDelta {
		return false
	}
	
	// Heuristique 2 : Ratio de changement
	if delta.ChangeRatio() > pc.DeltaThreshold {
		return false
	}
	
	// Heuristique 3 : Nombre de nœuds affectés
	if affectedNodesCount > pc.MaxAffectedNodesForDelta {
		return false
	}
	
	// Heuristique 4 : Changement de clé primaire
	if !pc.AllowPrimaryKeyChange && pc.hasPrimaryKeyChange(delta) {
		return false
	}
	
	// Toutes les conditions passées : utiliser delta
	return true
}

// hasPrimaryKeyChange vérifie si le delta contient un changement de clé primaire.
func (pc *PropagationConfig) hasPrimaryKeyChange(delta *FactDelta) bool {
	if len(pc.PrimaryKeyFields) == 0 {
		// Pas de clés primaires configurées : autoriser
		return false
	}
	
	for _, pkField := range pc.PrimaryKeyFields {
		if _, changed := delta.Fields[pkField]; changed {
			return true
		}
	}
	
	return false
}

// Clone crée une copie de la configuration.
func (pc *PropagationConfig) Clone() PropagationConfig {
	pkFields := make([]string, len(pc.PrimaryKeyFields))
	copy(pkFields, pc.PrimaryKeyFields)
	
	return PropagationConfig{
		DefaultMode:                   pc.DefaultMode,
		EnableDeltaPropagation:        pc.EnableDeltaPropagation,
		DeltaThreshold:                pc.DeltaThreshold,
		MinFieldsForDelta:             pc.MinFieldsForDelta,
		MaxAffectedNodesForDelta:      pc.MaxAffectedNodesForDelta,
		AllowPrimaryKeyChange:         pc.AllowPrimaryKeyChange,
		PrimaryKeyFields:              pkFields,
		EnableMetrics:                 pc.EnableMetrics,
		PropagationTimeout:            pc.PropagationTimeout,
		RetryOnError:                  pc.RetryOnError,
		MaxConcurrentPropagations:     pc.MaxConcurrentPropagations,
		EnableOptimisticPropagation:   pc.EnableOptimisticPropagation,
		LogPropagationDetails:         pc.LogPropagationDetails,
	}
}
```

### Tests : `rete/delta/propagation_config_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestPropagationMode_String(t *testing.T) {
	tests := []struct {
		mode PropagationMode
		want string
	}{
		{PropagationModeDelta, "Delta"},
		{PropagationModeClassic, "Classic"},
		{PropagationModeAuto, "Auto"},
		{PropagationMode(999), "Unknown"},
	}
	
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("PropagationMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultPropagationConfig(t *testing.T) {
	config := DefaultPropagationConfig()
	
	if config.DefaultMode != PropagationModeAuto {
		t.Errorf("Expected Auto mode, got %v", config.DefaultMode)
	}
	
	if !config.EnableDeltaPropagation {
		t.Error("Expected EnableDeltaPropagation = true")
	}
	
	if config.DeltaThreshold != 0.5 {
		t.Errorf("Expected DeltaThreshold = 0.5, got %v", config.DeltaThreshold)
	}
	
	if config.MinFieldsForDelta != 3 {
		t.Errorf("Expected MinFieldsForDelta = 3, got %d", config.MinFieldsForDelta)
	}
	
	if config.MaxAffectedNodesForDelta != 100 {
		t.Errorf("Expected MaxAffectedNodesForDelta = 100, got %d", config.MaxAffectedNodesForDelta)
	}
	
	if config.PropagationTimeout != 30*time.Second {
		t.Errorf("Expected timeout = 30s, got %v", config.PropagationTimeout)
	}
}

func TestPropagationConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    PropagationConfig
		wantError bool
	}{
		{
			name:      "valid default",
			config:    DefaultPropagationConfig(),
			wantError: false,
		},
		{
			name: "invalid delta threshold (negative)",
			config: PropagationConfig{
				DeltaThreshold:            -0.1,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid delta threshold (> 1)",
			config: PropagationConfig{
				DeltaThreshold:            1.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid min fields (negative)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         -1,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid max nodes (zero)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  0,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid timeout (negative)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        -time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid max concurrent (zero)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 0,
			},
			wantError: true,
		},
		{
			name: "edge case: zero threshold",
			config: PropagationConfig{
				DeltaThreshold:            0.0,
				MinFieldsForDelta:         0,
				MaxAffectedNodesForDelta:  1,
				PropagationTimeout:        0,
				MaxConcurrentPropagations: 1,
			},
			wantError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPropagationConfig_ShouldUseDelta_FeatureFlagDisabled(t *testing.T) {
	config := DefaultPropagationConfig()
	config.EnableDeltaPropagation = false
	
	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("field1", "old", "new")
	
	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when feature flag disabled")
	}
}

func TestPropagationConfig_ShouldUseDelta_ForcedMode(t *testing.T) {
	t.Run("forced classic", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.DefaultMode = PropagationModeClassic
		
		delta := NewFactDelta("Test~1", "Test")
		delta.FieldCount = 10
		delta.AddFieldChange("field1", "old", "new")
		
		if config.ShouldUseDelta(delta, 5) {
			t.Error("Expected false when mode forced to Classic")
		}
	})
	
	t.Run("forced delta", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.DefaultMode = PropagationModeDelta
		
		delta := NewFactDelta("Test~1", "Test")
		delta.FieldCount = 1 // Normalement trop peu
		delta.AddFieldChange("field1", "old", "new")
		
		if !config.ShouldUseDelta(delta, 5) {
			t.Error("Expected true when mode forced to Delta")
		}
	})
}

func TestPropagationConfig_ShouldUseDelta_MinFields(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.MinFieldsForDelta = 5
	
	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 3 // < 5
	delta.AddFieldChange("field1", "old", "new")
	
	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when field count below threshold")
	}
}

func TestPropagationConfig_ShouldUseDelta_ChangeRatio(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.DeltaThreshold = 0.3 // 30%
	config.MinFieldsForDelta = 0
	
	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	
	// Modifier 4 champs sur 10 = 40% > 30%
	for i := 0; i < 4; i++ {
		delta.AddFieldChange("field"+string(rune('0'+i)), "old", "new")
	}
	
	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when change ratio exceeds threshold")
	}
}

func TestPropagationConfig_ShouldUseDelta_AffectedNodes(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.MaxAffectedNodesForDelta = 10
	config.MinFieldsForDelta = 0
	
	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("field1", "old", "new")
	
	if config.ShouldUseDelta(delta, 15) {
		t.Error("Expected false when affected nodes exceed limit")
	}
}

func TestPropagationConfig_ShouldUseDelta_PrimaryKeyChange(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.AllowPrimaryKeyChange = false
	config.PrimaryKeyFields = []string{"id", "pk"}
	config.MinFieldsForDelta = 0
	
	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("id", "123", "456") // PK change
	
	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when PK changed and not allowed")
	}
}

func TestPropagationConfig_ShouldUseDelta_AllConditionsPass(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.MinFieldsForDelta = 5
	config.DeltaThreshold = 0.5
	config.MaxAffectedNodesForDelta = 50
	
	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("field1", "old", "new") // 10% change
	
	if !config.ShouldUseDelta(delta, 20) {
		t.Error("Expected true when all conditions pass")
	}
}

func TestPropagationConfig_hasPrimaryKeyChange(t *testing.T) {
	tests := []struct {
		name      string
		pkFields  []string
		delta     *FactDelta
		wantChange bool
	}{
		{
			name:     "no PK fields configured",
			pkFields: []string{},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("id", "1", "2")
				return d
			}(),
			wantChange: false,
		},
		{
			name:     "PK field changed",
			pkFields: []string{"id"},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("id", "1", "2")
				return d
			}(),
			wantChange: true,
		},
		{
			name:     "PK field not changed",
			pkFields: []string{"id"},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("name", "old", "new")
				return d
			}(),
			wantChange: false,
		},
		{
			name:     "multiple PK fields, one changed",
			pkFields: []string{"id", "tenant_id"},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("tenant_id", "A", "B")
				return d
			}(),
			wantChange: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultPropagationConfig()
			config.PrimaryKeyFields = tt.pkFields
			
			got := config.hasPrimaryKeyChange(tt.delta)
			if got != tt.wantChange {
				t.Errorf("hasPrimaryKeyChange() = %v, want %v", got, tt.wantChange)
			}
		})
	}
}

func TestPropagationConfig_Clone(t *testing.T) {
	original := PropagationConfig{
		DefaultMode:                   PropagationModeDelta,
		EnableDeltaPropagation:        true,
		DeltaThreshold:                0.3,
		MinFieldsForDelta:             5,
		MaxAffectedNodesForDelta:      50,
		AllowPrimaryKeyChange:         true,
		PrimaryKeyFields:              []string{"id", "pk"},
		EnableMetrics:                 true,
		PropagationTimeout:            10 * time.Second,
		RetryOnError:                  true,
		MaxConcurrentPropagations:     20,
		EnableOptimisticPropagation:   true,
		LogPropagationDetails:         true,
	}
	
	cloned := original.Clone()
	
	// Vérifier égalité valeurs
	if cloned.DefaultMode != original.DefaultMode {
		t.Error("DefaultMode not cloned")
	}
	if cloned.DeltaThreshold != original.DeltaThreshold {
		t.Error("DeltaThreshold not cloned")
	}
	
	// Vérifier indépendance slices
	if len(cloned.PrimaryKeyFields) != len(original.PrimaryKeyFields) {
		t.Error("PrimaryKeyFields length mismatch")
	}
	
	cloned.PrimaryKeyFields[0] = "modified"
	if original.PrimaryKeyFields[0] == "modified" {
		t.Error("Clone not independent (slice mutation affected original)")
	}
}
```

---

## 🔧 Tâche 2 : Stratégies de Propagation

### Fichier : `rete/delta/propagation_strategy.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// PropagationStrategy définit une stratégie de propagation.
//
// Cette interface permet d'implémenter différentes stratégies
// de propagation pour s'adapter à différents scénarios.
type PropagationStrategy interface {
	// GetName retourne le nom de la stratégie
	GetName() string
	
	// ShouldPropagate détermine si la propagation doit avoir lieu
	ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool
	
	// GetPropagationOrder retourne l'ordre de propagation des nœuds
	GetPropagationOrder(nodes []NodeReference) []NodeReference
}

// SequentialStrategy propage vers les nœuds dans l'ordre séquentiel.
//
// Cette stratégie est simple et prévisible : alpha → beta → terminal.
type SequentialStrategy struct{}

// GetName retourne "Sequential"
func (s *SequentialStrategy) GetName() string {
	return "Sequential"
}

// ShouldPropagate retourne toujours true (propage toujours)
func (s *SequentialStrategy) ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool {
	return len(affectedNodes) > 0
}

// GetPropagationOrder trie les nœuds par type : alpha, puis beta, puis terminal
func (s *SequentialStrategy) GetPropagationOrder(nodes []NodeReference) []NodeReference {
	// Séparer par type
	var alphaNodes, betaNodes, terminalNodes []NodeReference
	
	for _, node := range nodes {
		switch node.NodeType {
		case "alpha":
			alphaNodes = append(alphaNodes, node)
		case "beta":
			betaNodes = append(betaNodes, node)
		case "terminal":
			terminalNodes = append(terminalNodes, node)
		}
	}
	
	// Concaténer dans l'ordre alpha → beta → terminal
	ordered := make([]NodeReference, 0, len(nodes))
	ordered = append(ordered, alphaNodes...)
	ordered = append(ordered, betaNodes...)
	ordered = append(ordered, terminalNodes...)
	
	return ordered
}

// TopologicalStrategy propage en respectant les dépendances topologiques.
//
// Cette stratégie garantit qu'un nœud parent est toujours traité avant
// ses nœuds enfants (ordre topologique du graphe RETE).
type TopologicalStrategy struct {
	// nodeDepths stocke la profondeur de chaque nœud dans le graphe
	// (calculé lors de la construction du réseau)
	nodeDepths map[string]int
}

// NewTopologicalStrategy crée une nouvelle stratégie topologique
func NewTopologicalStrategy() *TopologicalStrategy {
	return &TopologicalStrategy{
		nodeDepths: make(map[string]int),
	}
}

// GetName retourne "Topological"
func (ts *TopologicalStrategy) GetName() string {
	return "Topological"
}

// ShouldPropagate retourne true si au moins un nœud est affecté
func (ts *TopologicalStrategy) ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool {
	return len(affectedNodes) > 0
}

// GetPropagationOrder trie les nœuds par profondeur topologique
func (ts *TopologicalStrategy) GetPropagationOrder(nodes []NodeReference) []NodeReference {
	// Si pas de profondeurs calculées, fallback sur ordre séquentiel
	if len(ts.nodeDepths) == 0 {
		sequential := &SequentialStrategy{}
		return sequential.GetPropagationOrder(nodes)
	}
	
	// Trier par profondeur croissante
	ordered := make([]NodeReference, len(nodes))
	copy(ordered, nodes)
	
	// Tri par insertion simple (taille typique petite)
	for i := 1; i < len(ordered); i++ {
		key := ordered[i]
		keyDepth := ts.getDepth(key.NodeID)
		j := i - 1
		
		for j >= 0 && ts.getDepth(ordered[j].NodeID) > keyDepth {
			ordered[j+1] = ordered[j]
			j--
		}
		ordered[j+1] = key
	}
	
	return ordered
}

// SetNodeDepth enregistre la profondeur d'un nœud
func (ts *TopologicalStrategy) SetNodeDepth(nodeID string, depth int) {
	ts.nodeDepths[nodeID] = depth
}

// getDepth retourne la profondeur d'un nœud (0 si inconnu)
func (ts *TopologicalStrategy) getDepth(nodeID string) int {
	if depth, exists := ts.nodeDepths[nodeID]; exists {
		return depth
	}
	return 0
}

// OptimizedStrategy est une stratégie hybride qui optimise selon le contexte.
//
// Elle combine plusieurs heuristiques :
// - Trier par type (alpha → beta → terminal)
// - Grouper par factType (meilleure localité cache)
// - Prioriser les nœuds avec moins de dépendances
type OptimizedStrategy struct{}

// GetName retourne "Optimized"
func (os *OptimizedStrategy) GetName() string {
	return "Optimized"
}

// ShouldPropagate retourne true si propagation justifiée
func (os *OptimizedStrategy) ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool {
	// Ne pas propager si aucun nœud affecté
	if len(affectedNodes) == 0 {
		return false
	}
	
	// Ne pas propager si delta vide (safety check)
	if delta.IsEmpty() {
		return false
	}
	
	return true
}

// GetPropagationOrder optimise l'ordre de propagation
func (os *OptimizedStrategy) GetPropagationOrder(nodes []NodeReference) []NodeReference {
	if len(nodes) == 0 {
		return nodes
	}
	
	// Étape 1 : Grouper par (type, factType) pour localité
	groups := make(map[string][]NodeReference)
	
	for _, node := range nodes {
		key := node.NodeType + ":" + node.FactType
		groups[key] = append(groups[key], node)
	}
	
	// Étape 2 : Ordre de traitement des groupes
	// Alpha d'abord, puis beta, puis terminal
	ordered := make([]NodeReference, 0, len(nodes))
	
	// Traiter alpha nodes par factType
	for key, group := range groups {
		if len(key) > 5 && key[:5] == "alpha" {
			ordered = append(ordered, group...)
		}
	}
	
	// Traiter beta nodes
	for key, group := range groups {
		if len(key) > 4 && key[:4] == "beta" {
			ordered = append(ordered, group...)
		}
	}
	
	// Traiter terminal nodes
	for key, group := range groups {
		if len(key) > 8 && key[:8] == "terminal" {
			ordered = append(ordered, group...)
		}
	}
	
	return ordered
}
```

### Tests : `rete/delta/propagation_strategy_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

func TestSequentialStrategy_GetName(t *testing.T) {
	strategy := &SequentialStrategy{}
	if strategy.GetName() != "Sequential" {
		t.Errorf("Expected 'Sequential', got '%s'", strategy.GetName())
	}
}

func TestSequentialStrategy_ShouldPropagate(t *testing.T) {
	strategy := &SequentialStrategy{}
	delta := NewFactDelta("Test~1", "Test")
	
	t.Run("with nodes", func(t *testing.T) {
		nodes := []NodeReference{{NodeID: "n1", NodeType: "alpha"}}
		if !strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected true when nodes present")
		}
	})
	
	t.Run("without nodes", func(t *testing.T) {
		nodes := []NodeReference{}
		if strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected false when no nodes")
		}
	})
}

func TestSequentialStrategy_GetPropagationOrder(t *testing.T) {
	strategy := &SequentialStrategy{}
	
	nodes := []NodeReference{
		{NodeID: "term1", NodeType: "terminal", FactType: "Test"},
		{NodeID: "alpha1", NodeType: "alpha", FactType: "Test"},
		{NodeID: "beta1", NodeType: "beta", FactType: "Test"},
		{NodeID: "alpha2", NodeType: "alpha", FactType: "Test"},
		{NodeID: "term2", NodeType: "terminal", FactType: "Test"},
	}
	
	ordered := strategy.GetPropagationOrder(nodes)
	
	if len(ordered) != len(nodes) {
		t.Fatalf("Expected %d nodes, got %d", len(nodes), len(ordered))
	}
	
	// Vérifier ordre : alphas d'abord
	alphaCount := 0
	for i, node := range ordered {
		if node.NodeType == "alpha" {
			alphaCount++
		} else if node.NodeType == "beta" {
			// Beta doit venir après tous les alphas
			if alphaCount != 2 {
				t.Errorf("Beta at position %d before all alphas", i)
			}
		}
	}
	
	if alphaCount != 2 {
		t.Errorf("Expected 2 alpha nodes, got %d", alphaCount)
	}
}

func TestTopologicalStrategy_GetName(t *testing.T) {
	strategy := NewTopologicalStrategy()
	if strategy.GetName() != "Topological" {
		t.Errorf("Expected 'Topological', got '%s'", strategy.GetName())
	}
}

func TestTopologicalStrategy_SetNodeDepth(t *testing.T) {
	strategy := NewTopologicalStrategy()
	
	strategy.SetNodeDepth("node1", 1)
	strategy.SetNodeDepth("node2", 2)
	
	if strategy.getDepth("node1") != 1 {
		t.Errorf("Expected depth 1, got %d", strategy.getDepth("node1"))
	}
	
	if strategy.getDepth("node2") != 2 {
		t.Errorf("Expected depth 2, got %d", strategy.getDepth("node2"))
	}
	
	if strategy.getDepth("unknown") != 0 {
		t.Errorf("Expected depth 0 for unknown node, got %d", strategy.getDepth("unknown"))
	}
}

func TestTopologicalStrategy_GetPropagationOrder(t *testing.T) {
	strategy := NewTopologicalStrategy()
	
	// Setup depths
	strategy.SetNodeDepth("node1", 3)
	strategy.SetNodeDepth("node2", 1)
	strategy.SetNodeDepth("node3", 2)
	
	nodes := []NodeReference{
		{NodeID: "node1", NodeType: "alpha"},
		{NodeID: "node2", NodeType: "alpha"},
		{NodeID: "node3", NodeType: "beta"},
	}
	
	ordered := strategy.GetPropagationOrder(nodes)
	
	// Vérifier ordre par profondeur : node2 (1) → node3 (2) → node1 (3)
	if ordered[0].NodeID != "node2" {
		t.Errorf("Expected node2 first, got %s", ordered[0].NodeID)
	}
	if ordered[1].NodeID != "node3" {
		t.Errorf("Expected node3 second, got %s", ordered[1].NodeID)
	}
	if ordered[2].NodeID != "node1" {
		t.Errorf("Expected node1 third, got %s", ordered[2].NodeID)
	}
}

func TestTopologicalStrategy_GetPropagationOrder_NoDepths(t *testing.T) {
	strategy := NewTopologicalStrategy()
	
	// Pas de profondeurs définies → fallback sequential
	nodes := []NodeReference{
		{NodeID: "term1", NodeType: "terminal"},
		{NodeID: "alpha1", NodeType: "alpha"},
	}
	
	ordered := strategy.GetPropagationOrder(nodes)
	
	// Devrait être ordonné type-first (alpha avant terminal)
	if len(ordered) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(ordered))
	}
	
	if ordered[0].NodeType != "alpha" {
		t.Error("Expected alpha first in fallback mode")
	}
}

func TestOptimizedStrategy_GetName(t *testing.T) {
	strategy := &OptimizedStrategy{}
	if strategy.GetName() != "Optimized" {
		t.Errorf("Expected 'Optimized', got '%s'", strategy.GetName())
	}
}

func TestOptimizedStrategy_ShouldPropagate(t *testing.T) {
	strategy := &OptimizedStrategy{}
	
	t.Run("empty delta", func(t *testing.T) {
		delta := NewFactDelta("Test~1", "Test")
		nodes := []NodeReference{{NodeID: "n1"}}
		
		if strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected false for empty delta")
		}
	})
	
	t.Run("no nodes", func(t *testing.T) {
		delta := NewFactDelta("Test~1", "Test")
		delta.AddFieldChange("field", "old", "new")
		nodes := []NodeReference{}
		
		if strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected false when no nodes")
		}
	})
	
	t.Run("valid propagation", func(t *testing.T) {
		delta := NewFactDelta("Test~1", "Test")
		delta.AddFieldChange("field", "old", "new")
		nodes := []NodeReference{{NodeID: "n1"}}
		
		if !strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected true for valid propagation")
		}
	})
}

func TestOptimizedStrategy_GetPropagationOrder(t *testing.T) {
	strategy := &OptimizedStrategy{}
	
	nodes := []NodeReference{
		{NodeID: "term1", NodeType: "terminal", FactType: "Product"},
		{NodeID: "alpha1", NodeType: "alpha", FactType: "Product"},
		{NodeID: "beta1", NodeType: "beta", FactType: "Order"},
		{NodeID: "alpha2", NodeType: "alpha", FactType: "Order"},
		{NodeID: "term2", NodeType: "terminal", FactType: "Product"},
	}
	
	ordered := strategy.GetPropagationOrder(nodes)
	
	if len(ordered) != len(nodes) {
		t.Fatalf("Expected %d nodes, got %d", len(nodes), len(ordered))
	}
	
	// Vérifier que alphas viennent avant betas et terminals
	seenBeta := false
	seenTerminal := false
	
	for _, node := range ordered {
		if node.NodeType == "beta" {
			seenBeta = true
		}
		if node.NodeType == "terminal" {
			seenTerminal = true
		}
		
		// Alpha ne devrait pas venir après beta ou terminal
		if node.NodeType == "alpha" && (seenBeta || seenTerminal) {
			t.Error("Alpha node found after beta or terminal")
		}
	}
}

func TestOptimizedStrategy_GetPropagationOrder_Empty(t *testing.T) {
	strategy := &OptimizedStrategy{}
	
	ordered := strategy.GetPropagationOrder([]NodeReference{})
	
	if len(ordered) != 0 {
		t.Errorf("Expected empty result, got %d nodes", len(ordered))
	}
}
```

---

## 🔧 Tâche 3 : Métriques de Propagation

### Fichier : `rete/delta/propagation_metrics.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"sync"
	"time"
)

// PropagationMetrics collecte des statistiques sur les propagations delta.
//
// Ces métriques permettent de monitorer la performance et l'efficacité
// du système de propagation sélective.
type PropagationMetrics struct {
	// Compteurs de propagations
	TotalPropagations     int64
	DeltaPropagations     int64
	ClassicPropagations   int64
	FailedPropagations    int64
	
	// Performance
	TotalPropagationTime  time.Duration
	AvgPropagationTime    time.Duration
	MinPropagationTime    time.Duration
	MaxPropagationTime    time.Duration
	
	// Efficacité delta
	TotalNodesEvaluated   int64
	NodesSkippedByDelta   int64
	AvgNodesPerPropagation float64
	
	// Champs modifiés
	TotalFieldsChanged    int64
	AvgFieldsPerPropagation float64
	
	// Fallbacks
	FallbacksDueToRatio   int64
	FallbacksDueToNodes   int64
	FallbacksDueToPK      int64
	FallbacksDueToError   int64
	
	// Timestamps
	FirstPropagation      time.Time
	LastPropagation       time.Time
	
	// Protection concurrence
	mutex                 sync.RWMutex
}

// NewPropagationMetrics crée une nouvelle instance de métriques.
func NewPropagationMetrics() *PropagationMetrics {
	return &PropagationMetrics{
		MinPropagationTime: time.Duration(1<<63 - 1), // Max int64
	}
}

// RecordDeltaPropagation enregistre une propagation delta.
//
// Paramètres :
//   - duration : temps pris par la propagation
//   - nodesAffected : nombre de nœuds affectés
//   - fieldsChanged : nombre de champs modifiés
func (pm *PropagationMetrics) RecordDeltaPropagation(
	duration time.Duration,
	nodesAffected int,
	fieldsChanged int,
) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	pm.TotalPropagations++
	pm.DeltaPropagations++
	pm.TotalNodesEvaluated += int64(nodesAffected)
	pm.TotalFieldsChanged += int64(fieldsChanged)
	
	pm.updateTiming(duration)
	pm.updateTimestamps()
	pm.recalculateAverages()
}

// RecordClassicPropagation enregistre une propagation classique (Retract+Insert).
//
// Paramètres :
//   - duration : temps pris
//   - totalNodes : nombre total de nœuds dans le réseau
func (pm *PropagationMetrics) RecordClassicPropagation(
	duration time.Duration,
	totalNodes int,
) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	pm.TotalPropagations++
	pm.ClassicPropagations++
	pm.TotalNodesEvaluated += int64(totalNodes)
	
	pm.updateTiming(duration)
	pm.updateTimestamps()
	pm.recalculateAverages()
}

// RecordFailedPropagation enregistre une propagation échouée.
func (pm *PropagationMetrics) RecordFailedPropagation() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	pm.TotalPropagations++
	pm.FailedPropagations++
	pm.updateTimestamps()
}

// RecordFallback enregistre un fallback vers mode classique.
//
// Paramètres :
//   - reason : raison du fallback ("ratio", "nodes", "pk", "error")
func (pm *PropagationMetrics) RecordFallback(reason string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	switch reason {
	case "ratio":
		pm.FallbacksDueToRatio++
	case "nodes":
		pm.FallbacksDueToNodes++
	case "pk":
		pm.FallbacksDueToPK++
	case "error":
		pm.FallbacksDueToError++
	}
}

// RecordNodesSkipped enregistre des nœuds évités grâce au delta.
//
// Paramètres :
//   - count : nombre de nœuds évités
func (pm *PropagationMetrics) RecordNodesSkipped(count int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	pm.NodesSkippedByDelta += int64(count)
}

// GetSnapshot retourne un instantané des métriques actuelles.
//
// Retourne une copie thread-safe des métriques.
func (pm *PropagationMetrics) GetSnapshot() PropagationMetrics {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	return PropagationMetrics{
		TotalPropagations:       pm.TotalPropagations,
		DeltaPropagations:       pm.DeltaPropagations,
		ClassicPropagations:     pm.ClassicPropagations,
		FailedPropagations:      pm.FailedPropagations,
		TotalPropagationTime:    pm.TotalPropagationTime,
		AvgPropagationTime:      pm.AvgPropagationTime,
		MinPropagationTime:      pm.MinPropagationTime,
		MaxPropagationTime:      pm.MaxPropagationTime,
		TotalNodesEvaluated:     pm.TotalNodesEvaluated,
		NodesSkippedByDelta:     pm.NodesSkippedByDelta,
		AvgNodesPerPropagation:  pm.AvgNodesPerPropagation,
		TotalFieldsChanged:      pm.TotalFieldsChanged,
		AvgFieldsPerPropagation: pm.AvgFieldsPerPropagation,
		FallbacksDueToRatio:     pm.FallbacksDueToRatio,
		FallbacksDueToNodes:     pm.FallbacksDueToNodes,
		FallbacksDueToPK:        pm.FallbacksDueToPK,
		FallbacksDueToError:     pm.FallbacksDueToError,
		FirstPropagation:        pm.FirstPropagation,
		LastPropagation:         pm.LastPropagation,
	}
}

// GetEfficiencyRatio retourne le ratio d'efficacité de la propagation delta.
//
// Ratio = NodesSkipped / TotalNodesEvaluated
// Plus ce ratio est élevé, plus le système delta est efficace.
//
// Retourne une valeur entre 0.0 et 1.0, ou 0.0 si aucune propagation.
func (pm *PropagationMetrics) GetEfficiencyRatio() float64 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	if pm.TotalNodesEvaluated == 0 {
		return 0.0
	}
	
	// Estimer les nœuds totaux qui auraient été évalués en mode classique
	totalClassicNodes := pm.TotalNodesEvaluated + pm.NodesSkippedByDelta
	
	if totalClassicNodes == 0 {
		return 0.0
	}
	
	return float64(pm.NodesSkippedByDelta) / float64(totalClassicNodes)
}

// GetDeltaUsageRatio retourne le ratio d'utilisation du mode delta.
//
// Ratio = DeltaPropagations / TotalPropagations
//
// Retourne une valeur entre 0.0 et 1.0.
func (pm *PropagationMetrics) GetDeltaUsageRatio() float64 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	if pm.TotalPropagations == 0 {
		return 0.0
	}
	
	return float64(pm.DeltaPropagations) / float64(pm.TotalPropagations)
}

// Reset réinitialise toutes les métriques à zéro.
func (pm *PropagationMetrics) Reset() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	pm.TotalPropagations = 0
	pm.DeltaPropagations = 0
	pm.ClassicPropagations = 0
	pm.FailedPropagations = 0
	pm.TotalPropagationTime = 0
	pm.AvgPropagationTime = 0
	pm.MinPropagationTime = time.Duration(1<<63 - 1)
	pm.MaxPropagationTime = 0
	pm.TotalNodesEvaluated = 0
	pm.NodesSkippedByDelta = 0
	pm.AvgNodesPerPropagation = 0
	pm.TotalFieldsChanged = 0
	pm.AvgFieldsPerPropagation = 0
	pm.FallbacksDueToRatio = 0
	pm.FallbacksDueToNodes = 0
	pm.FallbacksDueToPK = 0
	pm.FallbacksDueToError = 0
	pm.FirstPropagation = time.Time{}
	pm.LastPropagation = time.Time{}
}

// updateTiming met à jour les statistiques de timing.
// ATTENTION : doit être appelé avec mutex déjà acquis.
func (pm *PropagationMetrics) updateTiming(duration time.Duration) {
	pm.TotalPropagationTime += duration
	
	if duration < pm.MinPropagationTime {
		pm.MinPropagationTime = duration
	}
	
	if duration > pm.MaxPropagationTime {
		pm.MaxPropagationTime = duration
	}
}

// updateTimestamps met à jour les timestamps.
// ATTENTION : doit être appelé avec mutex déjà acquis.
func (pm *PropagationMetrics) updateTimestamps() {
	now := time.Now()
	
	if pm.FirstPropagation.IsZero() {
		pm.FirstPropagation = now
	}
	
	pm.LastPropagation = now
}

// recalculateAverages recalcule les moyennes.
// ATTENTION : doit être appelé avec mutex déjà acquis.
func (pm *PropagationMetrics) recalculateAverages() {
	if pm.TotalPropagations > 0 {
		pm.AvgPropagationTime = time.Duration(
			int64(pm.TotalPropagationTime) / pm.TotalPropagations,
		)
		
		pm.AvgNodesPerPropagation = float64(pm.TotalNodesEvaluated) /
			float64(pm.TotalPropagations)
		
		pm.AvgFieldsPerPropagation = float64(pm.TotalFieldsChanged) /
			float64(pm.TotalPropagations)
	}
}
```

---

## 🔧 Tâche 4 : DeltaPropagator Principal

### Fichier : `rete/delta/delta_propagator.go`

**Contenu** (partie 1/2) :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DeltaPropagator orchestre la propagation sélective des changements.
//
// Il coordonne la détection de delta, la recherche de nœuds affectés,
// et la propagation vers ces nœuds selon la stratégie configurée.
//
// Thread-safety : DeltaPropagator est safe pour utilisation concurrent.
type DeltaPropagator struct {
	// Dépendances
	detector  *DeltaDetector
	index     *DependencyIndex
	strategy  PropagationStrategy
	
	// Configuration
	config    PropagationConfig
	
	// Métriques
	metrics   *PropagationMetrics
	
	// Contrôle de concurrence
	semaphore chan struct{}
	mutex     sync.RWMutex
	
	// Callbacks pour interaction avec le réseau RETE
	// (définis comme interfaces pour éviter dépendances circulaires)
	onPropagate func(nodeID string, delta *FactDelta) error
}

// DeltaPropagatorBuilder construit un DeltaPropagator avec pattern builder.
type DeltaPropagatorBuilder struct {
	detector  *DeltaDetector
	index     *DependencyIndex
	strategy  PropagationStrategy
	config    PropagationConfig
	onPropagate func(string, *FactDelta) error
}

// NewDeltaPropagatorBuilder crée un nouveau builder.
func NewDeltaPropagatorBuilder() *DeltaPropagatorBuilder {
	return &DeltaPropagatorBuilder{
		config: DefaultPropagationConfig(),
	}
}

// WithDetector configure le détecteur de delta.
func (b *DeltaPropagatorBuilder) WithDetector(detector *DeltaDetector) *DeltaPropagatorBuilder {
	b.detector = detector
	return b
}

// WithIndex configure l'index de dépendances.
func (b *DeltaPropagatorBuilder) WithIndex(index *DependencyIndex) *DeltaPropagatorBuilder {
	b.index = index
	return b
}

// WithStrategy configure la stratégie de propagation.
func (b *DeltaPropagatorBuilder) WithStrategy(strategy PropagationStrategy) *DeltaPropagatorBuilder {
	b.strategy = strategy
	return b
}

// WithConfig configure la propagation.
func (b *DeltaPropagatorBuilder) WithConfig(config PropagationConfig) *DeltaPropagatorBuilder {
	b.config = config
	return b
}

// WithPropagateCallback configure le callback de propagation vers le réseau RETE.
func (b *DeltaPropagatorBuilder) WithPropagateCallback(
	callback func(string, *FactDelta) error,
) *DeltaPropagatorBuilder {
	b.onPropagate = callback
	return b
}

// Build construit le DeltaPropagator.
func (b *DeltaPropagatorBuilder) Build() (*DeltaPropagator, error) {
	// Valider configuration
	if err := b.config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	
	// Dépendances obligatoires
	if b.index == nil {
		return nil, fmt.Errorf("dependency index is required")
	}
	
	// Dépendances avec fallback
	if b.detector == nil {
		b.detector = NewDeltaDetector()
	}
	
	if b.strategy == nil {
		b.strategy = &SequentialStrategy{}
	}
	
	// Créer sémaphore pour contrôle de concurrence
	semaphore := make(chan struct{}, b.config.MaxConcurrentPropagations)
	
	return &DeltaPropagator{
		detector:    b.detector,
		index:       b.index,
		strategy:    b.strategy,
		config:      b.config,
		metrics:     NewPropagationMetrics(),
		semaphore:   semaphore,
		onPropagate: b.onPropagate,
	}, nil
}

// PropagateUpdate propage une mise à jour de fait.
//
// Cette méthode est le point d'entrée principal pour propager un Update.
//
// Paramètres :
//   - oldFact : fait avant modification
//   - newFact : fait après modification
//   - factID : identifiant interne du fait
//   - factType : type du fait
//
// Retourne une erreur si la propagation échoue.
func (dp *DeltaPropagator) PropagateUpdate(
	oldFact, newFact map[string]interface{},
	factID, factType string,
) error {
	return dp.PropagateUpdateWithContext(
		context.Background(),
		oldFact, newFact,
		factID, factType,
	)
}

// PropagateUpdateWithContext propage avec un contexte (timeout, cancellation).
//
// Paramètres :
//   - ctx : contexte de propagation (timeout, cancel)
//   - oldFact, newFact : faits avant/après
//   - factID, factType : identifiant et type
//
// Retourne une erreur si propagation échouée ou timeout.
func (dp *DeltaPropagator) PropagateUpdateWithContext(
	ctx context.