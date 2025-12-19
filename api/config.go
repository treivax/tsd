// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import "time"

// LogLevel représente le niveau de logging
type LogLevel int

const (
	LogLevelSilent LogLevel = iota
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// SelectionPolicy définit la politique de sélection pour les xuple-spaces
type SelectionPolicy string

const (
	SelectionFIFO   SelectionPolicy = "fifo"
	SelectionLIFO   SelectionPolicy = "lifo"
	SelectionRandom SelectionPolicy = "random"
)

// ConsumptionPolicy définit la politique de consommation pour les xuple-spaces
type ConsumptionPolicy string

const (
	ConsumptionOnce     ConsumptionPolicy = "once"
	ConsumptionPerAgent ConsumptionPolicy = "per-agent"
)

// RetentionPolicy définit la politique de rétention pour les xuple-spaces
type RetentionPolicy string

const (
	RetentionUnlimited RetentionPolicy = "unlimited"
	RetentionDuration  RetentionPolicy = "duration"
)

// XupleSpaceDefaults contient les valeurs par défaut pour les xuple-spaces
type XupleSpaceDefaults struct {
	Selection         SelectionPolicy
	Consumption       ConsumptionPolicy
	Retention         RetentionPolicy
	RetentionDuration time.Duration
	MaxSize           int
}

// Config contient la configuration du pipeline
type Config struct {
	LogLevel           LogLevel
	EnableMetrics      bool
	MaxFactsInMemory   int
	XupleSpaceDefaults *XupleSpaceDefaults
	EnableTransactions bool
	TransactionTimeout time.Duration
}

// DefaultConfig retourne la configuration par défaut
func DefaultConfig() *Config {
	return &Config{
		LogLevel:           LogLevelInfo,
		EnableMetrics:      true,
		MaxFactsInMemory:   0,
		EnableTransactions: true,
		TransactionTimeout: 30 * time.Second,
		XupleSpaceDefaults: &XupleSpaceDefaults{
			Selection:         SelectionFIFO,
			Consumption:       ConsumptionOnce,
			Retention:         RetentionUnlimited,
			RetentionDuration: 0,
			MaxSize:           0,
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

	switch defaults.Selection {
	case SelectionFIFO, SelectionLIFO, SelectionRandom:
	case "":
		defaults.Selection = SelectionFIFO
	default:
		return &ConfigError{
			Field:   "XupleSpaceDefaults.Selection",
			Message: "valeur invalide: " + string(defaults.Selection),
		}
	}

	switch defaults.Consumption {
	case ConsumptionOnce, ConsumptionPerAgent:
	case "":
		defaults.Consumption = ConsumptionOnce
	default:
		return &ConfigError{
			Field:   "XupleSpaceDefaults.Consumption",
			Message: "valeur invalide: " + string(defaults.Consumption),
		}
	}

	switch defaults.Retention {
	case RetentionUnlimited, RetentionDuration:
	case "":
		defaults.Retention = RetentionUnlimited
	default:
		return &ConfigError{
			Field:   "XupleSpaceDefaults.Retention",
			Message: "valeur invalide: " + string(defaults.Retention),
		}
	}

	if defaults.Retention == RetentionDuration && defaults.RetentionDuration <= 0 {
		return &ConfigError{
			Field:   "XupleSpaceDefaults.RetentionDuration",
			Message: "doit être > 0 quand Retention = duration",
		}
	}

	if defaults.MaxSize < 0 {
		return &ConfigError{
			Field:   "XupleSpaceDefaults.MaxSize",
			Message: "ne peut pas être négatif",
		}
	}

	return nil
}
