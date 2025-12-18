# 🔧 Prompt 02 - Package API Pipeline Complet

> **Objectif**: Créer le package `api` centralisant le pipeline complet TSD avec intégration automatique des xuples  
> **Dépendances**: Aucune (peut être exécuté en parallèle du Prompt 01)  
> **Contexte max**: 128k tokens  
> **Durée estimée**: 1 session

---

## ⚠️ CONTRAINTE ARCHITECTURALE STRICTE

**RÈGLE ABSOLUE**: Il est **STRICTEMENT INTERDIT** de créer des xuples directement dans les tests ou en batch (appel direct à `XupleManager.Create()`, `space.Add()`, etc.).

**Les xuples DOIVENT IMPÉRATIVEMENT être générés à partir de faits soumis au réseau RETE via des règles.**

✅ **CORRECT**:
```go
// Soumettre un fait au réseau RETE
network.Assert(ctx, fact)
// Le réseau évalue les règles, exécute l'action Xuple() qui crée le xuple
```

❌ **INTERDIT**:
```go
// NE JAMAIS faire ça dans les tests ou le code métier
xupleManager.Create(ctx, "space", fact)
space.Add(fact)
```

**Justification**:
- Garantit que tous les xuples passent par le réseau RETE
- Assure l'évaluation complète des règles et conditions
- Préserve la traçabilité et l'auditabilité
- Évite les contournements du pipeline qui introduisent des incohérences

Cette règle s'applique à **tous les tests, benchmarks, exemples et code de production**.

---

## 🎯 Objectif

Créer un nouveau package `api` qui servira de point d'entrée unique et simplifié pour l'utilisation de TSD. Ce package:
- Importe et intègre les packages `rete`, `xuples`, et `constraint`
- Fournit une API simple pour les utilisateurs (tests, serveurs, applications)
- Gère automatiquement toute la configuration (xuple-spaces, actions, etc.)
- Élimine le besoin de configuration manuelle
- Évite les cycles d'importation en étant au-dessus de tous les autres packages

**Avant** (utilisation complexe):
```go
import (
    "github.com/treivax/tsd/rete"
    "github.com/treivax/tsd/xuples"
)

// 50+ lignes de configuration manuelle
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)
network.SetXupleSpaceFactory(func(...) { ... })
// ... configuration complexe ...
pipeline := rete.NewConstraintPipeline()
result, metrics, err := pipeline.IngestFile("program.tsd", network, storage)
```

**Après** (utilisation simple):
```go
import "github.com/treivax/tsd/api"

// 3 lignes, tout est automatique
pipeline := api.NewPipeline()
result, err := pipeline.IngestFile("program.tsd")
// result contient réseau RETE + xuples + métriques + tout
```

---

## 📋 Structure du Package

### Arborescence

```
tsd/api/
├── doc.go              # Documentation du package
├── pipeline.go         # Pipeline principal (point d'entrée)
├── result.go           # Résultat d'ingestion
├── config.go           # Configuration optionnelle
├── errors.go           # Erreurs spécifiques à l'API
├── pipeline_test.go    # Tests unitaires pipeline
├── result_test.go      # Tests unitaires result
└── examples_test.go    # Exemples GoDoc
```

---

## 🛠️ Tâches à Réaliser

### Tâche 1: Documentation du Package

**Fichier**: `api/doc.go`

**Contenu**:

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

/*
Package api fournit une interface simplifiée pour utiliser le moteur de règles TSD.

Ce package est le point d'entrée recommandé pour toutes les applications utilisant TSD.
Il intègre automatiquement les packages rete, xuples, et constraint, et gère toute
la configuration nécessaire.

# Utilisation Basique

La manière la plus simple d'utiliser TSD est via le Pipeline:

	import "github.com/treivax/tsd/api"

	func main() {
		// Créer un pipeline
		pipeline := api.NewPipeline()

		// Ingérer un programme TSD
		result, err := pipeline.IngestFile("program.tsd")
		if err != nil {
			log.Fatal(err)
		}

		// Utiliser les résultats
		fmt.Printf("Types définis: %d\n", result.TypeCount())
		fmt.Printf("Règles actives: %d\n", result.RuleCount())
		fmt.Printf("Faits dans le réseau: %d\n", result.FactCount())
		fmt.Printf("Xuple-spaces créés: %d\n", result.XupleSpaceCount())
	}

# Accès aux Xuples

Les xuples créés par les règles sont accessibles via le résultat:

	result, _ := pipeline.IngestFile("monitoring.tsd")

	// Récupérer tous les xuples d'un xuple-space
	alerts := result.GetXuples("critical_alerts")
	for _, xuple := range alerts {
		fmt.Printf("Alert: %v\n", xuple.Fact.Fields)
	}

	// Consommer un xuple (retrieve)
	xuple, err := result.Retrieve("critical_alerts", "agent1")
	if err == nil {
		fmt.Printf("Consumed: %v\n", xuple.Fact.Fields)
	}

# Configuration Avancée

Pour une configuration personnalisée:

	config := &api.Config{
		LogLevel:          api.LogLevelDebug,
		EnableMetrics:     true,
		MaxFactsInMemory:  100000,
		XupleSpaceDefaults: &api.XupleSpaceDefaults{
			Selection:   api.SelectionFIFO,
			Consumption: api.ConsumptionOnce,
			Retention:   api.RetentionUnlimited,
		},
	}

	pipeline := api.NewPipelineWithConfig(config)
	result, err := pipeline.IngestFile("program.tsd")

# Ingestion Incrémentale

Le pipeline supporte l'ingestion incrémentale de plusieurs fichiers:

	pipeline := api.NewPipeline()

	// Charger les types et règles
	_, err := pipeline.IngestFile("types.tsd")
	if err != nil {
		log.Fatal(err)
	}

	// Ajouter plus de règles
	_, err = pipeline.IngestFile("additional-rules.tsd")
	if err != nil {
		log.Fatal(err)
	}

	// Soumettre des faits
	result, err := pipeline.IngestFile("facts.tsd")

Le réseau RETE est étendu de manière incrémentale et tous les faits précédents
sont automatiquement propagés aux nouvelles règles.

# Thread Safety

Le Pipeline est thread-safe. Plusieurs goroutines peuvent appeler IngestFile
en parallèle, mais notez que l'ordre d'exécution des règles peut varier.

Pour un contrôle strict de l'ordre, utilisez un seul goroutine ou synchronisez
explicitement les appels.

# Métriques

Les métriques d'ingestion sont disponibles dans le résultat:

	result, _ := pipeline.IngestFile("program.tsd")
	metrics := result.Metrics()

	fmt.Printf("Temps de parsing: %v\n", metrics.ParseDuration)
	fmt.Printf("Temps de construction réseau: %v\n", metrics.BuildDuration)
	fmt.Printf("Nombre de propagations: %d\n", metrics.PropagationCount)

# Gestion d'Erreurs

Les erreurs sont détaillées et incluent la position dans le fichier source:

	_, err := pipeline.IngestFile("invalid.tsd")
	if err != nil {
		if parseErr, ok := err.(*api.ParseError); ok {
			fmt.Printf("Erreur ligne %d, colonne %d: %s\n",
				parseErr.Line, parseErr.Column, parseErr.Message)
		}
	}

# Architecture

Le package api est construit au-dessus de:
  - constraint: Parser TSD (PEG)
  - rete: Moteur de règles (algorithme RETE)
  - xuples: Gestion des xuple-spaces et xuples

Il gère automatiquement:
  - Création du réseau RETE
  - Initialisation du XupleManager
  - Création des xuple-spaces à partir des définitions
  - Enregistrement des actions (Xuple, Print, etc.)
  - Configuration des handlers
  - Propagation des faits

L'utilisateur n'a besoin de connaître aucun détail d'implémentation.
*/
package api
```

### Tâche 2: Configuration

**Fichier**: `api/config.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import "time"

// LogLevel représente le niveau de logging
type LogLevel int

const (
	// LogLevelSilent désactive tous les logs
	LogLevelSilent LogLevel = iota
	// LogLevelError affiche uniquement les erreurs
	LogLevelError
	// LogLevelWarn affiche erreurs et avertissements
	LogLevelWarn
	// LogLevelInfo affiche informations, erreurs et avertissements (défaut)
	LogLevelInfo
	// LogLevelDebug affiche tous les logs y compris debug
	LogLevelDebug
)

// SelectionPolicy définit la politique de sélection pour les xuple-spaces
type SelectionPolicy string

const (
	// SelectionFIFO sélectionne le xuple le plus ancien (First In First Out)
	SelectionFIFO SelectionPolicy = "fifo"
	// SelectionLIFO sélectionne le xuple le plus récent (Last In First Out)
	SelectionLIFO SelectionPolicy = "lifo"
	// SelectionRandom sélectionne un xuple aléatoire
	SelectionRandom SelectionPolicy = "random"
)

// ConsumptionPolicy définit la politique de consommation pour les xuple-spaces
type ConsumptionPolicy string

const (
	// ConsumptionOnce permet de consommer chaque xuple une seule fois
	ConsumptionOnce ConsumptionPolicy = "once"
	// ConsumptionPerAgent permet à chaque agent de consommer chaque xuple une fois
	ConsumptionPerAgent ConsumptionPolicy = "per-agent"
)

// RetentionPolicy définit la politique de rétention pour les xuple-spaces
type RetentionPolicy string

const (
	// RetentionUnlimited conserve les xuples indéfiniment
	RetentionUnlimited RetentionPolicy = "unlimited"
	// RetentionDuration conserve les xuples pendant une durée limitée
	RetentionDuration RetentionPolicy = "duration"
)

// XupleSpaceDefaults contient les valeurs par défaut pour les xuple-spaces
type XupleSpaceDefaults struct {
	// Selection définit la politique de sélection par défaut
	Selection SelectionPolicy
	// Consumption définit la politique de consommation par défaut
	Consumption ConsumptionPolicy
	// Retention définit la politique de rétention par défaut
	Retention RetentionPolicy
	// RetentionDuration définit la durée de rétention (si Retention = RetentionDuration)
	RetentionDuration time.Duration
	// MaxSize définit la taille maximale d'un xuple-space (0 = illimité)
	MaxSize int
}

// Config contient la configuration du pipeline
type Config struct {
	// LogLevel définit le niveau de logging (défaut: LogLevelInfo)
	LogLevel LogLevel

	// EnableMetrics active la collecte de métriques détaillées (défaut: true)
	EnableMetrics bool

	// MaxFactsInMemory limite le nombre de faits en mémoire (0 = illimité)
	// Si la limite est atteinte, les faits les plus anciens sont évincés
	MaxFactsInMemory int

	// XupleSpaceDefaults définit les valeurs par défaut pour les xuple-spaces
	// créés sans configuration explicite
	XupleSpaceDefaults *XupleSpaceDefaults

	// EnableTransactions active le système de transactions (défaut: true)
	// Les transactions permettent le rollback en cas d'erreur
	EnableTransactions bool

	// TransactionTimeout définit le timeout pour les transactions
	// (défaut: 30 secondes)
	TransactionTimeout time.Duration
}

// DefaultConfig retourne la configuration par défaut
func DefaultConfig() *Config {
	return &Config{
		LogLevel:           LogLevelInfo,
		EnableMetrics:      true,
		MaxFactsInMemory:   0, // Illimité
		EnableTransactions: true,
		TransactionTimeout: 30 * time.Second,
		XupleSpaceDefaults: &XupleSpaceDefaults{
			Selection:         SelectionFIFO,
			Consumption:       ConsumptionOnce,
			Retention:         RetentionUnlimited,
			RetentionDuration: 0,
			MaxSize:           0, // Illimité
		},
	}
}

// Validate vérifie que la configuration est valide
func (c *Config) Validate() error {
	if c.TransactionTimeout < 0 {
		return &ConfigError{
			Field:   "TransactionTimeout",
			Message: "ne peut pas être négatif",
		}
	}

	if c.MaxFactsInMemory < 0 {
		return &ConfigError{
			Field:   "MaxFactsInMemory",
			Message: "ne peut pas être négatif",
		}
	}

	if c.XupleSpaceDefaults != nil {
		if err := c.validateXupleSpaceDefaults(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) validateXupleSpaceDefaults() error {
	defaults := c.XupleSpaceDefaults

	// Valider Selection
	switch defaults.Selection {
	case SelectionFIFO, SelectionLIFO, SelectionRandom:
		// OK
	case "":
		defaults.Selection = SelectionFIFO // Défaut
	default:
		return &ConfigError{
			Field:   "XupleSpaceDefaults.Selection",
			Message: "valeur invalide: " + string(defaults.Selection),
		}
	}

	// Valider Consumption
	switch defaults.Consumption {
	case ConsumptionOnce, ConsumptionPerAgent:
		// OK
	case "":
		defaults.Consumption = ConsumptionOnce // Défaut
	default:
		return &ConfigError{
			Field:   "XupleSpaceDefaults.Consumption",
			Message: "valeur invalide: " + string(defaults.Consumption),
		}
	}

	// Valider Retention
	switch defaults.Retention {
	case RetentionUnlimited, RetentionDuration:
		// OK
	case "":
		defaults.Retention = RetentionUnlimited // Défaut
	default:
		return &ConfigError{
			Field:   "XupleSpaceDefaults.Retention",
			Message: "valeur invalide: " + string(defaults.Retention),
		}
	}

	// Valider RetentionDuration si nécessaire
	if defaults.Retention == RetentionDuration && defaults.RetentionDuration <= 0 {
		return &ConfigError{
			Field:   "XupleSpaceDefaults.RetentionDuration",
			Message: "doit être > 0 quand Retention = duration",
		}
	}

	// Valider MaxSize
	if defaults.MaxSize < 0 {
		return &ConfigError{
			Field:   "XupleSpaceDefaults.MaxSize",
			Message: "ne peut pas être négatif",
		}
	}

	return nil
}
```

### Tâche 3: Erreurs

**Fichier**: `api/errors.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import "fmt"

// Error représente une erreur de l'API
type Error struct {
	// Type d'erreur
	Type ErrorType
	// Message d'erreur
	Message string
	// Erreur sous-jacente (optionnel)
	Cause error
}

// ErrorType représente le type d'erreur
type ErrorType string

const (
	// ErrorTypeParse erreur de parsing du fichier TSD
	ErrorTypeParse ErrorType = "parse"
	// ErrorTypeValidation erreur de validation (types, règles, etc.)
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeExecution erreur d'exécution (propagation, actions, etc.)
	ErrorTypeExecution ErrorType = "execution"
	// ErrorTypeConfig erreur de configuration
	ErrorTypeConfig ErrorType = "config"
	// ErrorTypeIO erreur d'entrée/sortie
	ErrorTypeIO ErrorType = "io"
	// ErrorTypeInternal erreur interne (bug)
	ErrorTypeInternal ErrorType = "internal"
)

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// ParseError représente une erreur de parsing avec position
type ParseError struct {
	// Fichier source
	Filename string
	// Ligne (1-based)
	Line int
	// Colonne (1-based)
	Column int
	// Message d'erreur
	Message string
	// Erreur sous-jacente
	Cause error
}

func (e *ParseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s:%d:%d: %s: %v", e.Filename, e.Line, e.Column, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s:%d:%d: %s", e.Filename, e.Line, e.Column, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Cause
}

// ConfigError représente une erreur de configuration
type ConfigError struct {
	// Champ de configuration en erreur
	Field string
	// Message d'erreur
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration invalide pour '%s': %s", e.Field, e.Message)
}

// XupleSpaceError représente une erreur liée aux xuple-spaces
type XupleSpaceError struct {
	// Nom du xuple-space
	SpaceName string
	// Opération tentée
	Operation string
	// Message d'erreur
	Message string
	// Erreur sous-jacente
	Cause error
}

func (e *XupleSpaceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("xuple-space '%s': %s: %s: %v", e.SpaceName, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("xuple-space '%s': %s: %s", e.SpaceName, e.Operation, e.Message)
}

func (e *XupleSpaceError) Unwrap() error {
	return e.Cause
}
```

### Tâche 4: Résultat

**Fichier**: `api/result.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// Result contient le résultat d'une ingestion de programme TSD
type Result struct {
	// Réseau RETE construit
	network *rete.ReteNetwork

	// XupleManager pour accéder aux xuples
	xupleManager xuples.XupleManager

	// Métriques d'ingestion
	metrics *Metrics
}

// Metrics contient les métriques d'ingestion
type Metrics struct {
	// Durée totale de l'ingestion
	TotalDuration time.Duration

	// Durée du parsing
	ParseDuration time.Duration

	// Durée de construction du réseau
	BuildDuration time.Duration

	// Durée de propagation des faits
	PropagationDuration time.Duration

	// Nombre de types définis
	TypeCount int

	// Nombre de règles créées
	RuleCount int

	// Nombre de faits soumis
	FactCount int

	// Nombre de xuple-spaces créés
	XupleSpaceCount int

	// Nombre de propagations effectuées
	PropagationCount int

	// Nombre d'actions exécutées
	ActionCount int
}

// Network retourne le réseau RETE sous-jacent
// Utilisez cette méthode uniquement pour des opérations avancées
func (r *Result) Network() *rete.ReteNetwork {
	return r.network
}

// XupleManager retourne le XupleManager
// Utilisez cette méthode uniquement pour des opérations avancées
func (r *Result) XupleManager() xuples.XupleManager {
	return r.xupleManager
}

// Metrics retourne les métriques d'ingestion
func (r *Result) Metrics() *Metrics {
	return r.metrics
}

// TypeCount retourne le nombre de types définis
func (r *Result) TypeCount() int {
	return r.metrics.TypeCount
}

// RuleCount retourne le nombre de règles actives
func (r *Result) RuleCount() int {
	return r.metrics.RuleCount
}

// FactCount retourne le nombre de faits dans le réseau
func (r *Result) FactCount() int {
	return r.metrics.FactCount
}

// XupleSpaceCount retourne le nombre de xuple-spaces créés
func (r *Result) XupleSpaceCount() int {
	return r.metrics.XupleSpaceCount
}

// GetXuples retourne tous les xuples d'un xuple-space
func (r *Result) GetXuples(spaceName string) ([]*xuples.Xuple, error) {
	if r.xupleManager == nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "GetXuples",
			Message:   "XupleManager non initialisé",
		}
	}

	space, err := r.xupleManager.GetXupleSpace(spaceName)
	if err != nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "GetXuples",
			Message:   "xuple-space non trouvé",
			Cause:     err,
		}
	}

	return space.ListAll(), nil
}

// Retrieve récupère et consomme un xuple d'un xuple-space selon sa politique
func (r *Result) Retrieve(spaceName string, agentID string) (*xuples.Xuple, error) {
	if r.xupleManager == nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "Retrieve",
			Message:   "XupleManager non initialisé",
		}
	}

	space, err := r.xupleManager.GetXupleSpace(spaceName)
	if err != nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "Retrieve",
			Message:   "xuple-space non trouvé",
			Cause:     err,
		}
	}

	xuple, err := space.Retrieve(agentID)
	if err != nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "Retrieve",
			Message:   "échec de récupération",
			Cause:     err,
		}
	}

	return xuple, nil
}

// XupleSpaceNames retourne les noms de tous les xuple-spaces
func (r *Result) XupleSpaceNames() []string {
	if r.xupleManager == nil {
		return []string{}
	}
	return r.xupleManager.ListXupleSpaces()
}

// XupleCount retourne le nombre de xuples dans un xuple-space
func (r *Result) XupleCount(spaceName string) (int, error) {
	xuples, err := r.GetXuples(spaceName)
	if err != nil {
		return 0, err
	}
	return len(xuples), nil
}

// Summary retourne un résumé texte du résultat
func (r *Result) Summary() string {
	summary := fmt.Sprintf("=== Résultat d'Ingestion TSD ===\n")
	summary += fmt.Sprintf("Types définis:        %d\n", r.TypeCount())
	summary += fmt.Sprintf("Règles actives:       %d\n", r.RuleCount())
	summary += fmt.Sprintf("Faits dans réseau:    %d\n", r.FactCount())
	summary += fmt.Sprintf("Xuple-spaces créés:   %d\n", r.XupleSpaceCount())

	if r.xupleManager != nil {
		for _, spaceName := range r.XupleSpaceNames() {
			count, _ := r.XupleCount(spaceName)
			summary += fmt.Sprintf("  - %s: %d xuples\n", spaceName, count)
		}
	}

	summary += fmt.Sprintf("\nMétriques de Performance:\n")
	summary += fmt.Sprintf("Durée totale:         %v\n", r.metrics.TotalDuration)
	summary += fmt.Sprintf("  - Parsing:          %v\n", r.metrics.ParseDuration)
	summary += fmt.Sprintf("  - Construction:     %v\n", r.metrics.BuildDuration)
	summary += fmt.Sprintf("  - Propagation:      %v\n", r.metrics.PropagationDuration)
	summary += fmt.Sprintf("Propagations:         %d\n", r.metrics.PropagationCount)
	summary += fmt.Sprintf("Actions exécutées:    %d\n", r.metrics.ActionCount)

	return summary
}
```

### Tâche 5: Pipeline Principal

**Fichier**: `api/pipeline.go`

**Contenu** (voir partie suivante pour le code complet - fichier long):

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// Pipeline est le point d'entrée principal pour utiliser TSD
type Pipeline struct {
	// Configuration
	config *Config

	// Réseau RETE (persistant entre ingestions)
	network *rete.ReteNetwork

	// Storage pour les faits
	storage rete.Storage

	// XupleManager (persistant)
	xupleManager xuples.XupleManager

	// Pipeline RETE sous-jacent
	retePipeline *rete.ConstraintPipeline

	// Mutex pour thread-safety
	mu sync.RWMutex
}

// NewPipeline crée un nouveau pipeline avec la configuration par défaut
func NewPipeline() *Pipeline {
	return NewPipelineWithConfig(DefaultConfig())
}

// NewPipelineWithConfig crée un nouveau pipeline avec une configuration personnalisée
func NewPipelineWithConfig(config *Config) *Pipeline {
	if config == nil {
		config = DefaultConfig()
	}

	// Valider la configuration
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("configuration invalide: %v", err))
	}

	// Créer le storage
	storage := rete.NewMemoryStorage()

	// Créer le réseau RETE
	network := rete.NewReteNetwork(storage)

	// Créer le XupleManager
	xupleManager := xuples.NewXupleManager()

	// Configurer le réseau avec le XupleManager
	network.SetXupleManager(xupleManager)

	// Configurer le handler Xuple
	network.SetXupleHandler(func(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
		return xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
	})

	// Créer le pipeline RETE
	retePipeline := rete.NewConstraintPipeline()

	// Configurer le logger selon le niveau
	logger := createLogger(config.LogLevel)
	retePipeline.SetLogger(logger)

	// Configurer la factory pour créer automatiquement les xuple-spaces
	network.SetXupleSpaceFactory(createXupleSpaceFactory(xupleManager, config))

	return &Pipeline{
		config:       config,
		network:      network,
		storage:      storage,
		xupleManager: xupleManager,
		retePipeline: retePipeline,
	}
}

// IngestFile ingère un fichier TSD et retourne le résultat
func (p *Pipeline) IngestFile(filename string) (*Result, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startTime := time.Now()

	// Vérifier que le fichier existe
	if _, err := os.Stat(filename); err != nil {
		return nil, &Error{
			Type:    ErrorTypeIO,
			Message: "fichier inaccessible",
			Cause:   err,
		}
	}

	// Ingérer via le pipeline RETE
	parseStart := time.Now()
	network, reteMetrics, err := p.retePipeline.IngestFile(filename, p.network, p.storage)
	if err != nil {
		return nil, p.wrapError(err, filename)
	}
	p.network = network

	// Construire les métriques
	metrics := &Metrics{
		TotalDuration:       time.Since(startTime),
		ParseDuration:       time.Since(parseStart),
		BuildDuration:       reteMetrics.BuildDuration,
		PropagationDuration: reteMetrics.PropagationDuration,
		TypeCount:           reteMetrics.TypeCount,
		RuleCount:           reteMetrics.RuleCount,
		FactCount:           reteMetrics.FactCount,
		XupleSpaceCount:     len(p.xupleManager.ListXupleSpaces()),
		PropagationCount:    reteMetrics.PropagationCount,
		ActionCount:         reteMetrics.ActionCount,
	}

	// Créer le résultat
	result := &Result{
		network:      p.network,
		xupleManager: p.xupleManager,
		metrics:      metrics,
	}

	return result, nil
}

// IngestString ingère un programme TSD depuis une chaîne
func (p *Pipeline) IngestString(program string) (*Result, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startTime := time.Now()

	// Créer un fichier temporaire
	tmpFile, err := os.CreateTemp("", "tsd-*.tsd")
	if err != nil {
		return nil, &Error{
			Type:    ErrorTypeIO,
			Message: "impossible de créer fichier temporaire",
			Cause:   err,
		}
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Écrire le programme
	if _, err := tmpFile.WriteString(program); err != nil {
		return nil, &Error{
			Type:    ErrorTypeIO,
			Message: "impossible d'écrire dans fichier temporaire",
			Cause:   err,
		}
	}
	tmpFile.Close()

	// Déléguer à IngestFile (sans lock, on est déjà locké)
	p.mu.Unlock()
	result, err := p.IngestFile(tmpFile.Name())
	p.mu.Lock()

	if err != nil {
		return nil, err
	}

	// Ajuster le temps total
	result.metrics.TotalDuration = time.Since(startTime)

	return result, nil
}

// Reset réinitialise complètement le pipeline (efface tout)
func (p *Pipeline) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Créer un nouveau storage
	p.storage = rete.NewMemoryStorage()

	// Créer un nouveau réseau
	p.network = rete.NewReteNetwork(p.storage)

	// Créer un nouveau XupleManager
	p.xupleManager = xuples.NewXupleManager()

	// Reconfigurer
	p.network.SetXupleManager(p.xupleManager)
	p.network.SetXupleHandler(func(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
		return p.xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
	})
	p.network.SetXupleSpaceFactory(createXupleSpaceFactory(p.xupleManager, p.config))
}

// Fonctions utilitaires privées

func createLogger(level LogLevel) *rete.Logger {
	var reteLevel rete.LogLevel
	switch level {
	case LogLevelSilent:
		reteLevel = rete.LogLevelSilent
	case LogLevelError:
		reteLevel = rete.LogLevelError
	case LogLevelWarn:
		reteLevel = rete.LogLevelWarn
	case LogLevelInfo:
		reteLevel = rete.LogLevelInfo
	case LogLevelDebug:
		reteLevel = rete.LogLevelDebug
	default:
		reteLevel = rete.LogLevelInfo
	}
	return rete.NewLogger(reteLevel, os.Stdout)
}

func createXupleSpaceFactory(xupleManager xuples.XupleManager, config *Config) rete.XupleSpaceFactoryFunc {
	return func(network *rete.ReteNetwork, definitions []interface{}) error {
		for _, xsDef := range definitions {
			if err := createXupleSpaceFromDefinition(xupleManager, xsDef, config); err != nil {
				return err
			}
		}
		return nil
	}
}

func createXupleSpaceFromDefinition(xupleManager xuples.XupleManager, xsDef interface{}, config *Config) error {
	xsMap, ok := xsDef.(map[string]interface{})
	if !ok {
		return fmt.Errorf("format de xuple-space invalide: %T", xsDef)
	}

	name, _ := xsMap["name"].(string)
	if name == "" {
		return fmt.Errorf("nom de xuple-space manquant")
	}

	// Parser les politiques
	selPolicy := parseSelectionPolicy(xsMap, config)
	consPolicy := parseConsumptionPolicy(xsMap, config)
	retPolicy := parseRetentionPolicy(xsMap, config)

	// Créer la configuration
	xsConfig := xuples.XupleSpaceConfig{
		Name:              name,
		SelectionPolicy:   selPolicy,
		ConsumptionPolicy: consPolicy,
		RetentionPolicy:   retPolicy,
		MaxSize:           config.XupleSpaceDefaults.MaxSize,
	}

	return xupleManager.CreateXupleSpace(name, xsConfig)
}

func parseSelectionPolicy(xsMap map[string]interface{}, config *Config) xuples.SelectionPolicy {
	selectionStr, _ := xsMap["selectionPolicy"].(string)
	switch selectionStr {
	case "fifo":
		return xuples.NewFIFOSelectionPolicy()
	case "lifo":
		return xuples.NewLIFOSelectionPolicy()
	case "random":
		return xuples.NewRandomSelectionPolicy()
	default:
		// Utiliser défaut de la config
		switch config.XupleSpaceDefaults.Selection {
		case SelectionLIFO:
			return xuples.NewLIFOSelectionPolicy()
		case SelectionRandom:
			return xuples.NewRandomSelectionPolicy()
		default:
			return xuples.NewFIFOSelectionPolicy()
		}
	}
}

func parseConsumptionPolicy(xsMap map[string]interface{}, config *Config) xuples.ConsumptionPolicy {
	consumptionMap, _ := xsMap["consumptionPolicy"].(map[string]interface{})
	consType, _ := consumptionMap["type"].(string)
	switch consType {
	case "once":
		return xuples.NewOnceConsumptionPolicy()
	case "per-agent":
		return xuples.NewPerAgentConsumptionPolicy()
	default:
		// Utiliser défaut de la config
		if config.XupleSpaceDefaults.Consumption == ConsumptionPerAgent {
			return xuples.NewPerAgentConsumptionPolicy()
		}
		return xuples.NewOnceConsumptionPolicy()
	}
}

func parseRetentionPolicy(xsMap map[string]interface{}, config *Config) xuples.RetentionPolicy {
	retentionMap, _ := xsMap["retentionPolicy"].(map[string]interface{})
	retType, _ := retentionMap["type"].(string)
	switch retType {
	case "unlimited":
		return xuples.NewUnlimitedRetentionPolicy()
	case "duration":
		duration := 0
		if d, ok := retentionMap["duration"].(float64); ok {
			duration = int(d)
		}
		if duration > 0 {
			return xuples.NewDurationRetentionPolicy(time.Duration(duration) * time.Second)
		}
		// Si duration invalide, utiliser défaut config
		if config.XupleSpaceDefaults.Retention == RetentionDuration {
			return xuples.NewDurationRetentionPolicy(config.XupleSpaceDefaults.RetentionDuration)
		}
		return xuples.NewUnlimitedRetentionPolicy()
	default:
		// Utiliser défaut config
		if config.XupleSpaceDefaults.Retention == RetentionDuration {
			return xuples.NewDurationRetentionPolicy(config.XupleSpaceDefaults.RetentionDuration)
		}
		return xuples.NewUnlimitedRetentionPolicy()
	}
}

func (p *Pipeline) wrapError(err error, filename string) error {
	// Tenter de wrapper en ParseError si possible
	// (À adapter selon le format d'erreur du parser constraint)
	if parseErr, ok := err.(*constraint.ParseError); ok {
		return &ParseError{
			Filename: filename,
			Line:     parseErr.Line,
			Column:   parseErr.Column,
			Message:  parseErr.Message,
			Cause:    parseErr,
		}
	}

	// Sinon, erreur générique
	return &Error{
		Type:    ErrorTypeExecution,
		Message: "erreur d'ingestion",
		Cause:   err,
	}
}
```

---

## 🧪 Tests (suite dans fichiers suivants)

Les tests unitaires et exemples sont dans les fichiers:
- `api/pipeline_test.go`
- `api/result_test.go`
- `api/config_test.go`
- `api/examples_test.go`

---

## ✅ Checklist de Validation

- [ ] Package `api` créé avec tous les fichiers
- [ ] `NewPipeline()` fonctionne avec config par défaut
- [ ] `IngestFile()` ingère correctement un fichier TSD
- [ ] `IngestString()` ingère correctement une chaîne
- [ ] Xuple-spaces créés automatiquement
- [ ] XupleManager configuré et accessible
- [ ] Actions Xuple enregistrées automatiquement
- [ ] Métriques collectées et accessibles
- [ ] Erreurs bien typées et informatives
- [ ] Thread-safe (tests de concurrence)
- [ ] Documentation GoDoc complète
- [ ] Tous les tests passent
- [ ] Couverture > 80%
- [ ] `make validate` passe

---

**Prochaine étape**: Après validation, passer au **Prompt 03 - Création Automatique Xuple-Spaces**