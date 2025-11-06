package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

// Config représente la configuration complète du module constraint
type Config struct {
	Parser    domain.ParserConfig    `json:"parser"`
	Validator domain.ValidatorConfig `json:"validator"`
	Logger    domain.LoggerConfig    `json:"logger"`
	Debug     bool                   `json:"debug"`
	Version   string                 `json:"version"`
}

// DefaultConfig retourne une configuration par défaut
func DefaultConfig() *Config {
	return &Config{
		Parser: domain.ParserConfig{
			MaxExpressions: 1000,
			Debug:          false,
			Recover:        true,
		},
		Validator: domain.ValidatorConfig{
			StrictMode: true,
			AllowedOperators: []string{
				"==", "!=", "<", ">", "<=", ">=",
				"AND", "OR", "NOT",
				"+", "-", "*", "/", "%",
			},
			MaxDepth: 20,
		},
		Logger: domain.LoggerConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
		Debug:   false,
		Version: "1.0.0",
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

	// Créer le répertoire si nécessaire
	dir := filepath.Dir(cm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.filePath, data, 0644); err != nil {
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
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[config.Logger.Level] {
		return fmt.Errorf("invalid logger level: %s", config.Logger.Level)
	}

	validFormats := map[string]bool{
		"json": true, "text": true, "plain": true,
	}
	if !validFormats[config.Logger.Format] {
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

// Clone crée une copie de la configuration
func (cm *ConfigManager) Clone() *ConfigManager {
	configCopy := *cm.config
	return &ConfigManager{
		config:   &configCopy,
		filePath: cm.filePath,
	}
}

// Reset remet la configuration aux valeurs par défaut
func (cm *ConfigManager) Reset() {
	cm.config = DefaultConfig()
}