// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
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
	level      LogLevel
	output     io.Writer
	mu         sync.RWMutex
	timestamps bool
	prefix     string
}

var (
	// DefaultLogger est le logger global par défaut
	DefaultLogger *Logger
)

func init() {
	DefaultLogger = NewLogger(LogLevelInfo, os.Stdout)
}

// NewLogger crée un nouveau logger avec le niveau spécifié
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:      level,
		output:     output,
		timestamps: true,
		prefix:     "[RETE]",
	}
}

// SetTimestamps active ou désactive les timestamps
func (l *Logger) SetTimestamps(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timestamps = enabled
}

// SetPrefix change le préfixe du logger
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// SetOutput change la destination des logs
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
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
	l.mu.RLock()
	output := l.output
	timestamps := l.timestamps
	prefix := l.prefix
	l.mu.RUnlock()

	message := fmt.Sprintf(format, args...)

	var logLine string
	if timestamps {
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		logLine = fmt.Sprintf("%s %s [%s] %s\n", timestamp, prefix, level, message)
	} else {
		logLine = fmt.Sprintf("%s [%s] %s\n", prefix, level, message)
	}

	fmt.Fprint(output, logLine)
}

// WithContext retourne un nouveau logger avec un contexte ajouté au préfixe
func (l *Logger) WithContext(context string) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return &Logger{
		level:      l.level,
		output:     l.output,
		timestamps: l.timestamps,
		prefix:     fmt.Sprintf("%s[%s]", l.prefix, context),
	}
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
