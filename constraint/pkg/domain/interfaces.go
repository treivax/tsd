// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package domain

// Parser interface définit les méthodes de parsing
type Parser interface {
	// Parse analyse une chaîne de caractères et retourne l'AST
	Parse(filename string, input []byte) (interface{}, error)

	// ParseFile analyse un fichier
	ParseFile(filename string) (interface{}, error)

	// ParseReader analyse depuis un Reader
	ParseReader(filename string, reader interface{}) (interface{}, error)
}

// Validator interface définit les méthodes de validation
type Validator interface {
	// ValidateProgram valide un programme complet
	ValidateProgram(program interface{}) error

	// ValidateTypes valide les définitions de types
	ValidateTypes(types []TypeDefinition) error

	// ValidateExpression valide une expression/règle
	ValidateExpression(expr Expression, types []TypeDefinition) error

	// ValidateConstraint valide une contrainte
	ValidateConstraint(constraint interface{}, variables []TypedVariable, types []TypeDefinition) error
}

// TypeChecker interface définit les méthodes de vérification de types
type TypeChecker interface {
	// GetFieldType retourne le type d'un champ
	GetFieldType(fieldAccess interface{}, variables []TypedVariable, types []TypeDefinition) (string, error)

	// GetValueType retourne le type d'une valeur
	GetValueType(value interface{}) string

	// ValidateTypeCompatibility vérifie la compatibilité entre types
	ValidateTypeCompatibility(leftType, rightType, operator string) error
}

// ActionValidator interface définit les méthodes de validation d'actions
type ActionValidator interface {
	// ValidateAction valide une action
	ValidateAction(action *Action) error

	// ValidateJobCall valide un appel de fonction/job
	ValidateJobCall(jobCall JobCall) error
}

// TypeRegistry interface définit les méthodes de gestion des types
type TypeRegistry interface {
	// RegisterType enregistre un nouveau type
	RegisterType(typeDef TypeDefinition) error

	// GetType récupère un type par son nom
	GetType(name string) (*TypeDefinition, error)

	// HasType vérifie si un type existe
	HasType(name string) bool

	// ListTypes retourne tous les types enregistrés
	ListTypes() []TypeDefinition

	// GetTypeFields retourne les champs d'un type
	GetTypeFields(typeName string) (map[string]string, error)
}

// ConstraintEngine interface définit les méthodes d'évaluation des contraintes
type ConstraintEngine interface {
	// EvaluateConstraint évalue une contrainte avec des données
	EvaluateConstraint(constraint interface{}, data map[string]interface{}) (bool, error)

	// EvaluateExpression évalue une expression complète
	EvaluateExpression(expr Expression, data map[string]interface{}) (bool, error)
}

// ActionExecutor interface définit les méthodes d'exécution d'actions
type ActionExecutor interface {
	// ExecuteAction exécute une action
	ExecuteAction(action *Action, context map[string]interface{}) error

	// RegisterJobHandler enregistre un gestionnaire d'action
	RegisterJobHandler(jobName string, handler JobHandler) error
}

// JobHandler définit un gestionnaire d'action personnalisé
type JobHandler func(args []string, context map[string]interface{}) error

// ProgramManager interface définit les méthodes de gestion de programmes
type ProgramManager interface {
	// LoadProgram charge un programme depuis un fichier ou une chaîne
	LoadProgram(source string) (*Program, error)

	// SaveProgram sauvegarde un programme
	SaveProgram(program *Program, destination string) error

	// ValidateAndLoad charge et valide un programme
	ValidateAndLoad(source string) (*Program, error)

	// ExecuteProgram exécute un programme avec des données
	ExecuteProgram(program *Program, data map[string]interface{}) error
}

// Logger interface définit les méthodes de logging
type Logger interface {
	// Debug log un message de debug
	Debug(message string, fields ...interface{})

	// Info log un message d'information
	Info(message string, fields ...interface{})

	// Warn log un avertissement
	Warn(message string, fields ...interface{})

	// Error log une erreur
	Error(message string, err error, fields ...interface{})
}

// MetricsCollector interface définit les méthodes de collecte de métriques
type MetricsCollector interface {
	// IncrementParsed incrémente le nombre de programmes parsés
	IncrementParsed()

	// IncrementValidated incrémente le nombre de programmes validés
	IncrementValidated()

	// IncrementExecuted incrémente le nombre de programmes exécutés
	IncrementExecuted()

	// RecordParseTime enregistre le temps de parsing
	RecordParseTime(duration interface{})

	// RecordValidationTime enregistre le temps de validation
	RecordValidationTime(duration interface{})

	// GetMetrics retourne les métriques actuelles
	GetMetrics() map[string]interface{}
}

// ConfigProvider interface définit les méthodes de configuration
type ConfigProvider interface {
	// GetParserConfig retourne la configuration du parser
	GetParserConfig() ParserConfig

	// GetValidatorConfig retourne la configuration du validateur
	GetValidatorConfig() ValidatorConfig

	// GetLoggerConfig retourne la configuration du logger
	GetLoggerConfig() LoggerConfig

	// IsDebugEnabled vérifie si le mode debug est activé
	IsDebugEnabled() bool
}

// ParserConfig configuration du parser
type ParserConfig struct {
	MaxExpressions int  `json:"max_expressions"`
	Debug          bool `json:"debug"`
	Recover        bool `json:"recover"`
}

// ValidatorConfig configuration du validateur
type ValidatorConfig struct {
	StrictMode       bool     `json:"strict_mode"`
	AllowedOperators []string `json:"allowed_operators"`
	MaxDepth         int      `json:"max_depth"`
}

// LoggerConfig configuration du logger
type LoggerConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}
