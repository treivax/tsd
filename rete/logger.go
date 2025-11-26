// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// LogLevel représente le niveau de logging
type LogLevel int

const (
	// LogLevelSilent désactive tous les logs
	LogLevelSilent LogLevel = iota
	// LogLevelError affiche uniquement les erreurs
	LogLevelError
	// LogLevelWarn affiche warnings et erreurs
	LogLevelWarn
	// LogLevelInfo affiche info, warnings et erreurs
	LogLevelInfo
	// LogLevelDebug affiche tous les logs incluant debug
	LogLevelDebug
)

// Logger est un logger configurable pour RETE
type Logger struct {
	level  LogLevel
	output io.Writer
	mu     sync.RWMutex
}

var (
	// DefaultLogger est le logger global par défaut
	DefaultLogger *Logger
	once          sync.Once
)

func init() {
	DefaultLogger = NewLogger(LogLevelInfo, os.Stdout)
}

// NewLogger crée un nouveau logger avec le niveau spécifié
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
	}
}

// SetLevel change le niveau de logging
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel retourne le niveau de logging actuel
func (l *Logger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// Debug log un message de debug
func (l *Logger) Debug(format string, args ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level >= LogLevelDebug {
		l.log("DEBUG", format, args...)
	}
}

// Info log un message informatif
func (l *Logger) Info(format string, args ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level >= LogLevelInfo {
		l.log("INFO", format, args...)
	}
}

// Warn log un avertissement
func (l *Logger) Warn(format string, args ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level >= LogLevelWarn {
		l.log("WARN", format, args...)
	}
}

// Error log une erreur
func (l *Logger) Error(format string, args ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level >= LogLevelError {
		l.log("ERROR", format, args...)
	}
}

// log est la méthode interne de logging
func (l *Logger) log(level string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.SetOutput(l.output)
	log.SetFlags(0) // Pas de timestamp par défaut pour compatibilité avec les tests
	log.Printf("%s", message)
}

// SetGlobalLogLevel change le niveau du logger global
func SetGlobalLogLevel(level LogLevel) {
	DefaultLogger.SetLevel(level)
}

// GetGlobalLogLevel retourne le niveau du logger global
func GetGlobalLogLevel() LogLevel {
	return DefaultLogger.GetLevel()
}

// Helper functions pour utiliser le logger global

// Debug log un message de debug avec le logger global
func Debug(format string, args ...interface{}) {
	DefaultLogger.Debug(format, args...)
}

// Info log un message informatif avec le logger global
func Info(format string, args ...interface{}) {
	DefaultLogger.Info(format, args...)
}

// Warn log un avertissement avec le logger global
func Warn(format string, args ...interface{}) {
	DefaultLogger.Warn(format, args...)
}

// Error log une erreur avec le logger global
func Error(format string, args ...interface{}) {
	DefaultLogger.Error(format, args...)
}
