// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

// Configuration defaults - constantes nommées
const (
	DefaultMaxExpressions = 1000
	DefaultMaxDepth       = 20
	DefaultDebug          = false
	DefaultRecover        = true
	DefaultStrictMode     = true
	DefaultVersion        = "1.0.0"

	DefaultLogLevel  = "info"
	DefaultLogFormat = "json"
	DefaultLogOutput = "stdout"

	// Permissions fichiers
	DefaultDirPermissions  = 0755
	DefaultFilePermissions = 0644
)

// Opérateurs autorisés par défaut
var defaultAllowedOperators = []string{
	"==", "!=", "<", ">", "<=", ">=",
	"AND", "OR", "NOT",
	"+", "-", "*", "/", "%",
}

// Niveaux de log valides
var validLogLevels = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
}

// Formats de log valides
var validLogFormats = map[string]bool{
	"json":  true,
	"text":  true,
	"plain": true,
}

// Variables d'environnement supportées
const (
	EnvPrefix         = "CONSTRAINT_"
	EnvMaxExpressions = EnvPrefix + "MAX_EXPRESSIONS"
	EnvMaxDepth       = EnvPrefix + "MAX_DEPTH"
	EnvDebug          = EnvPrefix + "DEBUG"
	EnvStrictMode     = EnvPrefix + "STRICT_MODE"
	EnvLogLevel       = EnvPrefix + "LOG_LEVEL"
	EnvLogFormat      = EnvPrefix + "LOG_FORMAT"
	EnvLogOutput      = EnvPrefix + "LOG_OUTPUT"
	EnvConfigFile     = EnvPrefix + "CONFIG_FILE"
)

// Config représente la configuration complète du module constraint
type Config struct {
	Parser    domain.ParserConfig    `json:"parser"`
	Validator domain.ValidatorConfig `json:"validator"`
	Logger    domain.LoggerConfig    `json:"logger"`
	Debug     bool                   `json:"debug"`
	Version   string                 `json:"version"`
}

// DefaultConfig retourne une configuration par défaut avec constantes nommées
func DefaultConfig() *Config {
	// Copie du slice pour éviter partage de mémoire
	allowedOps := make([]string, len(defaultAllowedOperators))
	copy(allowedOps, defaultAllowedOperators)

	return &Config{
		Parser: domain.ParserConfig{
			MaxExpressions: DefaultMaxExpressions,
			Debug:          DefaultDebug,
			Recover:        DefaultRecover,
		},
		Validator: domain.ValidatorConfig{
			StrictMode:       DefaultStrictMode,
			AllowedOperators: allowedOps,
			MaxDepth:         DefaultMaxDepth,
		},
		Logger: domain.LoggerConfig{
			Level:  DefaultLogLevel,
			Format: DefaultLogFormat,
			Output: DefaultLogOutput,
		},
		Debug:   DefaultDebug,
		Version: DefaultVersion,
	}
}

// ConfigManager gère la configuration du module
type ConfigManager struct {
	config   *Config
	filePath string
}

// NewConfigManager crée un nouveau gestionnaire de configuration
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		config:   DefaultConfig(),
		filePath: configPath,
	}
}

// LoadFromFile charge la configuration depuis un fichier
func (cm *ConfigManager) LoadFromFile() error {
	if cm.filePath == "" {
		return fmt.Errorf("config file path not set")
	}

	// Vérifier si le fichier existe
	if _, err := os.Stat(cm.filePath); os.IsNotExist(err) {
		// Si le fichier n'existe pas, créer la configuration par défaut
		return cm.SaveToFile()
	}

	data, err := os.ReadFile(cm.filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", cm.filePath, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", cm.filePath, err)
	}

	cm.config = &config
	return nil
}

// SaveToFile sauvegarde la configuration dans un fichier
func (cm *ConfigManager) SaveToFile() error {
	if cm.filePath == "" {
		return fmt.Errorf("config file path not set")
	}

	// Créer le répertoire si nécessaire avec permissions configurables
	dir := filepath.Dir(cm.filePath)
	if err := os.MkdirAll(dir, os.FileMode(DefaultDirPermissions)); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.filePath, data, os.FileMode(DefaultFilePermissions)); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", cm.filePath, err)
	}

	return nil
}

// GetConfig retourne la configuration actuelle
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// SetConfig définit une nouvelle configuration
func (cm *ConfigManager) SetConfig(config *Config) {
	cm.config = config
}

// GetParserConfig retourne la configuration du parser
func (cm *ConfigManager) GetParserConfig() domain.ParserConfig {
	return cm.config.Parser
}

// GetValidatorConfig retourne la configuration du validateur
func (cm *ConfigManager) GetValidatorConfig() domain.ValidatorConfig {
	return cm.config.Validator
}

// GetLoggerConfig retourne la configuration du logger
func (cm *ConfigManager) GetLoggerConfig() domain.LoggerConfig {
	return cm.config.Logger
}

// IsDebugEnabled vérifie si le mode debug est activé
func (cm *ConfigManager) IsDebugEnabled() bool {
	return cm.config.Debug || cm.config.Parser.Debug
}

// UpdateParserConfig met à jour la configuration du parser
func (cm *ConfigManager) UpdateParserConfig(config domain.ParserConfig) {
	cm.config.Parser = config
}

// UpdateValidatorConfig met à jour la configuration du validateur
func (cm *ConfigManager) UpdateValidatorConfig(config domain.ValidatorConfig) {
	cm.config.Validator = config
}

// UpdateLoggerConfig met à jour la configuration du logger
func (cm *ConfigManager) UpdateLoggerConfig(config domain.LoggerConfig) {
	cm.config.Logger = config
}

// SetDebug active/désactive le mode debug
func (cm *ConfigManager) SetDebug(enabled bool) {
	cm.config.Debug = enabled
}

// Validate valide la configuration
func (cm *ConfigManager) Validate() error {
	config := cm.config

	// Validation du parser
	if config.Parser.MaxExpressions <= 0 {
		return fmt.Errorf("parser.max_expressions must be positive, got %d", config.Parser.MaxExpressions)
	}

	// Validation du validateur
	if config.Validator.MaxDepth <= 0 {
		return fmt.Errorf("validator.max_depth must be positive, got %d", config.Validator.MaxDepth)
	}

	if len(config.Validator.AllowedOperators) == 0 {
		return fmt.Errorf("validator.allowed_operators cannot be empty")
	}

	// Validation du logger
	if !validLogLevels[config.Logger.Level] {
		return fmt.Errorf("invalid logger level: %s", config.Logger.Level)
	}

	if !validLogFormats[config.Logger.Format] {
		return fmt.Errorf("invalid logger format: %s", config.Logger.Format)
	}

	return nil
}

// String retourne une représentation JSON de la configuration
func (cm *ConfigManager) String() string {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Sprintf("ConfigManager{error: %v}", err)
	}
	return string(data)
}

// Clone crée une copie profonde de la configuration
func (cm *ConfigManager) Clone() *ConfigManager {
	// Deep copy de la config
	configCopy := *cm.config

	// Deep copy du slice AllowedOperators
	configCopy.Validator.AllowedOperators = make([]string, len(cm.config.Validator.AllowedOperators))
	copy(configCopy.Validator.AllowedOperators, cm.config.Validator.AllowedOperators)

	return &ConfigManager{
		config:   &configCopy,
		filePath: cm.filePath,
	}
}

// Reset remet la configuration aux valeurs par défaut
func (cm *ConfigManager) Reset() {
	cm.config = DefaultConfig()
}

// LoadFromEnv charge la configuration depuis les variables d'environnement
// Les variables d'environnement surchargent les valeurs du fichier de configuration
func (cm *ConfigManager) LoadFromEnv() error {
	// Parser MaxExpressions
	if val := os.Getenv(EnvMaxExpressions); val != "" {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid %s: %w", EnvMaxExpressions, err)
		}
		cm.config.Parser.MaxExpressions = parsed
	}

	// Parser MaxDepth
	if val := os.Getenv(EnvMaxDepth); val != "" {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid %s: %w", EnvMaxDepth, err)
		}
		cm.config.Validator.MaxDepth = parsed
	}

	// Parser Debug
	if val := os.Getenv(EnvDebug); val != "" {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid %s: %w", EnvDebug, err)
		}
		cm.config.Debug = parsed
	}

	// Parser StrictMode
	if val := os.Getenv(EnvStrictMode); val != "" {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid %s: %w", EnvStrictMode, err)
		}
		cm.config.Validator.StrictMode = parsed
	}

	// Logger Level
	if val := os.Getenv(EnvLogLevel); val != "" {
		cm.config.Logger.Level = strings.ToLower(val)
	}

	// Logger Format
	if val := os.Getenv(EnvLogFormat); val != "" {
		cm.config.Logger.Format = strings.ToLower(val)
	}

	// Logger Output
	if val := os.Getenv(EnvLogOutput); val != "" {
		cm.config.Logger.Output = val
	}

	// Valider après chargement
	return cm.Validate()
}

// MergeConfig fusionne une configuration avec la configuration actuelle
// Les valeurs non-nulles de la config source écrasent les valeurs actuelles
func (cm *ConfigManager) MergeConfig(source *Config) {
	if source == nil {
		return
	}

	// Merge Parser config
	if source.Parser.MaxExpressions > 0 {
		cm.config.Parser.MaxExpressions = source.Parser.MaxExpressions
	}
	// Note: Debug et Recover sont des bool, on ne peut pas distinguer "non défini" de "false"
	// On les merge systématiquement
	cm.config.Parser.Debug = source.Parser.Debug
	cm.config.Parser.Recover = source.Parser.Recover

	// Merge Validator config
	if len(source.Validator.AllowedOperators) > 0 {
		cm.config.Validator.AllowedOperators = make([]string, len(source.Validator.AllowedOperators))
		copy(cm.config.Validator.AllowedOperators, source.Validator.AllowedOperators)
	}
	if source.Validator.MaxDepth > 0 {
		cm.config.Validator.MaxDepth = source.Validator.MaxDepth
	}
	cm.config.Validator.StrictMode = source.Validator.StrictMode

	// Merge Logger config
	if source.Logger.Level != "" {
		cm.config.Logger.Level = source.Logger.Level
	}
	if source.Logger.Format != "" {
		cm.config.Logger.Format = source.Logger.Format
	}
	if source.Logger.Output != "" {
		cm.config.Logger.Output = source.Logger.Output
	}

	// Merge autres champs
	if source.Version != "" {
		cm.config.Version = source.Version
	}
	cm.config.Debug = source.Debug
}

// GetConfigFilePath retourne le chemin du fichier de configuration depuis ENV ou valeur par défaut
func GetConfigFilePath(defaultPath string) string {
	if path := os.Getenv(EnvConfigFile); path != "" {
		return path
	}
	return defaultPath
}
