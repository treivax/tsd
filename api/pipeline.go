// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// Pipeline est le point d'entrée principal pour utiliser TSD
type Pipeline struct {
	config       *Config
	network      *rete.ReteNetwork
	storage      rete.Storage
	xupleManager xuples.XupleManager
	retePipeline *rete.ConstraintPipeline
	mu           sync.RWMutex
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

	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("configuration invalide: %v", err))
	}

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	xupleManager := xuples.NewXupleManager()

	// Configurer le handler pour l'action Xuple
	network.SetXupleManager(xupleManager)
	network.SetXupleHandler(func(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
		return xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
	})

	retePipeline := rete.NewConstraintPipeline()
	logger := createLogger(config.LogLevel)
	retePipeline.SetLogger(logger)

	// Créer le pipeline
	p := &Pipeline{
		config:       config,
		network:      network,
		storage:      storage,
		xupleManager: xupleManager,
		retePipeline: retePipeline,
	}

	// Configurer le callback pour créer les xuple-spaces dès qu'ils sont détectés
	retePipeline.SetOnXupleSpacesDetected(p.createXupleSpacesFromDefinitionsCallback)

	return p
}

// IngestFile ingère un fichier TSD et retourne le résultat
func (p *Pipeline) IngestFile(filename string) (*Result, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startTime := time.Now()

	if _, err := os.Stat(filename); err != nil {
		return nil, &Error{
			Type:    ErrorTypeIO,
			Message: "fichier inaccessible",
			Cause:   err,
		}
	}

	network, reteMetrics, err := p.retePipeline.IngestFile(filename, p.network, p.storage)
	if err != nil {
		return nil, p.wrapError(err, filename)
	}
	p.network = network

	metrics := &Metrics{
		TotalDuration:       time.Since(startTime),
		ParseDuration:       reteMetrics.ParsingDuration,
		BuildDuration:       reteMetrics.RuleCreationDuration + reteMetrics.TypeCreationDuration,
		PropagationDuration: reteMetrics.PropagationDuration,
		TypeCount:           reteMetrics.TypesAdded,
		RuleCount:           reteMetrics.RulesAdded,
		FactCount:           reteMetrics.FactsSubmitted,
		XupleSpaceCount:     len(p.xupleManager.ListXupleSpaces()),
		PropagationCount:    reteMetrics.FactsPropagated,
		ActionCount:         reteMetrics.PropagationTargets,
	}

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

	if _, err := tmpFile.WriteString(program); err != nil {
		return nil, &Error{
			Type:    ErrorTypeIO,
			Message: "impossible d'écrire dans fichier temporaire",
			Cause:   err,
		}
	}
	tmpFile.Close()

	p.mu.Unlock()
	result, err := p.IngestFile(tmpFile.Name())
	p.mu.Lock()

	if err != nil {
		return nil, err
	}

	result.metrics.TotalDuration = time.Since(startTime)
	return result, nil
}

// Reset réinitialise complètement le pipeline
func (p *Pipeline) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.storage = rete.NewMemoryStorage()
	p.network = rete.NewReteNetwork(p.storage)
	p.xupleManager = xuples.NewXupleManager()

	// Configurer le handler pour l'action Xuple
	p.network.SetXupleManager(p.xupleManager)
	p.network.SetXupleHandler(func(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
		return p.xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
	})
}

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

// createXupleSpacesFromDefinitionsCallback est le callback appelé par le pipeline RETE
// dès que les définitions de xuple-spaces sont détectées, AVANT la soumission des faits inline.
// Cela garantit que les xuple-spaces existent quand les actions Xuple sont exécutées.
func (p *Pipeline) createXupleSpacesFromDefinitionsCallback(network *rete.ReteNetwork, definitions []interface{}) error {

	for _, xsDef := range definitions {
		if err := p.createXupleSpaceFromDefinition(xsDef); err != nil {
			return err
		}
	}
	return nil
}

// createXupleSpaceFromDefinition crée un xuple-space à partir de sa définition parsée.
func (p *Pipeline) createXupleSpaceFromDefinition(xsDef interface{}) error {
	xsMap, ok := xsDef.(map[string]interface{})
	if !ok {
		return fmt.Errorf("format de xuple-space invalide: %T", xsDef)
	}

	name, _ := xsMap["name"].(string)
	if name == "" {
		return fmt.Errorf("nom de xuple-space manquant")
	}

	selPolicy := p.parseSelectionPolicy(xsMap)
	consPolicy := p.parseConsumptionPolicy(xsMap)
	retPolicy := p.parseRetentionPolicy(xsMap)

	// Parser max-size avec fallback vers config
	maxSize := p.config.XupleSpaceDefaults.MaxSize
	if ms, ok := xsMap["maxSize"]; ok {
		switch v := ms.(type) {
		case int:
			maxSize = v
		case float64:
			maxSize = int(v)
		}
	}

	xsConfig := xuples.XupleSpaceConfig{
		Name:              name,
		SelectionPolicy:   selPolicy,
		ConsumptionPolicy: consPolicy,
		RetentionPolicy:   retPolicy,
		MaxSize:           maxSize,
	}

	return p.xupleManager.CreateXupleSpace(name, xsConfig)
}

func (p *Pipeline) parseSelectionPolicy(xsMap map[string]interface{}) xuples.SelectionPolicy {
	selectionStr, _ := xsMap["selectionPolicy"].(string)
	switch selectionStr {
	case "fifo":
		return xuples.NewFIFOSelectionPolicy()
	case "lifo":
		return xuples.NewLIFOSelectionPolicy()
	case "random":
		return xuples.NewRandomSelectionPolicy()
	default:
		switch p.config.XupleSpaceDefaults.Selection {
		case SelectionLIFO:
			return xuples.NewLIFOSelectionPolicy()
		case SelectionRandom:
			return xuples.NewRandomSelectionPolicy()
		default:
			return xuples.NewFIFOSelectionPolicy()
		}
	}
}

func (p *Pipeline) parseConsumptionPolicy(xsMap map[string]interface{}) xuples.ConsumptionPolicy {
	consumptionMap, _ := xsMap["consumptionPolicy"].(map[string]interface{})
	consType, _ := consumptionMap["type"].(string)
	switch consType {
	case "once":
		return xuples.NewOnceConsumptionPolicy()
	case "per-agent":
		return xuples.NewPerAgentConsumptionPolicy()
	default:
		if p.config.XupleSpaceDefaults.Consumption == ConsumptionPerAgent {
			return xuples.NewPerAgentConsumptionPolicy()
		}
		return xuples.NewOnceConsumptionPolicy()
	}
}

func (p *Pipeline) parseRetentionPolicy(xsMap map[string]interface{}) xuples.RetentionPolicy {
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
		if p.config.XupleSpaceDefaults.Retention == RetentionDuration {
			return xuples.NewDurationRetentionPolicy(p.config.XupleSpaceDefaults.RetentionDuration)
		}
		return xuples.NewUnlimitedRetentionPolicy()
	default:
		if p.config.XupleSpaceDefaults.Retention == RetentionDuration {
			return xuples.NewDurationRetentionPolicy(p.config.XupleSpaceDefaults.RetentionDuration)
		}
		return xuples.NewUnlimitedRetentionPolicy()
	}
}

func (p *Pipeline) wrapError(err error, filename string) error {
	return &Error{
		Type:    ErrorTypeExecution,
		Message: "erreur d'ingestion",
		Cause:   err,
	}
}
