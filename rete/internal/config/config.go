package config

import (
"time"
)

// Config représente la configuration globale du système RETE.
type Config struct {
Storage StorageConfig `json:"storage"`
Network NetworkConfig `json:"network"`
Logger  LoggerConfig  `json:"logger"`
}

// StorageConfig configuration pour le système de persistance.
type StorageConfig struct {
Type     string        `json:"type"`     // "memory" ou "etcd"
Endpoint string        `json:"endpoint"` // pour etcd
Prefix   string        `json:"prefix"`   // préfixe des clés
Timeout  time.Duration `json:"timeout"`
}

// NetworkConfig configuration pour le réseau RETE.
type NetworkConfig struct {
MaxNodes        int           `json:"max_nodes"`
MaxFactsPerNode int           `json:"max_facts_per_node"`
GCInterval      time.Duration `json:"gc_interval"` // garbage collection
}

// LoggerConfig configuration pour le logging.
type LoggerConfig struct {
Level  string `json:"level"`  // "debug", "info", "warn", "error"
Format string `json:"format"` // "json" ou "text"
Output string `json:"output"` // chemin fichier ou "stdout"
}

// DefaultConfig retourne la configuration par défaut.
func DefaultConfig() *Config {
return &Config{
Storage: StorageConfig{
Type:    "memory",
Prefix:  "/rete/nodes",
Timeout: 30 * time.Second,
},
Network: NetworkConfig{
MaxNodes:        1000,
MaxFactsPerNode: 10000,
GCInterval:      5 * time.Minute,
},
Logger: LoggerConfig{
Level:  "info",
Format: "text",
Output: "stdout",
},
}
}

// Validate vérifie la validité de la configuration.
func (c *Config) Validate() error {
// Validation basique
if c.Storage.Type != "memory" && c.Storage.Type != "etcd" {
return &ValidationError{
Field:   "storage.type",
Value:   c.Storage.Type,
Message: "must be 'memory' or 'etcd'",
}
}

validLevels := map[string]bool{
"debug": true, "info": true, "warn": true, "error": true,
}
if !validLevels[c.Logger.Level] {
return &ValidationError{
Field:   "logger.level",
Value:   c.Logger.Level,
Message: "must be one of: debug, info, warn, error",
}
}

return nil
}

// ValidationError représente une erreur de validation de config.
type ValidationError struct {
Field   string
Value   interface{}
Message string
}

func (e *ValidationError) Error() string {
return "config validation error on field '" + e.Field + "': " + e.Message
}
