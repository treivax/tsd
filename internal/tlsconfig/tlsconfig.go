// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package tlsconfig provides centralized TLS configuration for TSD server and client.
// This package eliminates duplication and ensures consistent security settings.
package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

const (
	// MinTLSVersion est la version TLS minimale acceptée (TLS 1.2)
	MinTLSVersion = tls.VersionTLS12
)

// ServerConfig contient la configuration TLS pour le serveur
type ServerConfig struct {
	// CertFile est le chemin vers le certificat serveur
	CertFile string

	// KeyFile est le chemin vers la clé privée
	KeyFile string

	// MinVersion est la version TLS minimale (défaut : TLS 1.2)
	MinVersion uint16

	// CipherSuites est la liste des cipher suites autorisées
	CipherSuites []uint16

	// PreferServerCipherSuites indique si le serveur préfère ses cipher suites
	PreferServerCipherSuites bool
}

// ClientConfig contient la configuration TLS pour le client
type ClientConfig struct {
	// CAFile est le chemin vers le certificat CA pour vérifier le serveur
	CAFile string

	// InsecureSkipVerify désactive la vérification TLS (⚠️ développement uniquement)
	InsecureSkipVerify bool

	// MinVersion est la version TLS minimale (défaut : TLS 1.2)
	MinVersion uint16
}

// DefaultCipherSuites retourne la liste recommandée de cipher suites
func DefaultCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}
}

// NewServerTLSConfig crée une configuration TLS pour le serveur
func NewServerTLSConfig(config *ServerConfig) (*tls.Config, error) {
	if config == nil {
		return nil, fmt.Errorf("config ne peut pas être nil")
	}

	// Valider que les fichiers existent
	if config.CertFile == "" {
		return nil, fmt.Errorf("certFile requis")
	}
	if config.KeyFile == "" {
		return nil, fmt.Errorf("keyFile requis")
	}

	if _, err := os.Stat(config.CertFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("certificat non trouvé: %s", config.CertFile)
	}
	if _, err := os.Stat(config.KeyFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("clé privée non trouvée: %s", config.KeyFile)
	}

	// Définir les valeurs par défaut
	minVersion := config.MinVersion
	if minVersion == 0 {
		minVersion = MinTLSVersion
	}

	cipherSuites := config.CipherSuites
	if len(cipherSuites) == 0 {
		cipherSuites = DefaultCipherSuites()
	}

	return &tls.Config{
		MinVersion:               minVersion,
		CipherSuites:             cipherSuites,
		PreferServerCipherSuites: config.PreferServerCipherSuites,
	}, nil
}

// NewClientTLSConfig crée une configuration TLS pour le client
func NewClientTLSConfig(config *ClientConfig) (*tls.Config, error) {
	if config == nil {
		return nil, fmt.Errorf("config ne peut pas être nil")
	}

	// Définir version minimale
	minVersion := config.MinVersion
	if minVersion == 0 {
		minVersion = MinTLSVersion
	}

	tlsConfig := &tls.Config{
		MinVersion:         minVersion,
		InsecureSkipVerify: config.InsecureSkipVerify,
	}

	// Si mode insecure, on s'arrête là
	if config.InsecureSkipVerify {
		return tlsConfig, nil
	}

	// Charger le CA si fourni
	if config.CAFile != "" {
		if _, err := os.Stat(config.CAFile); err == nil {
			caCert, err := os.ReadFile(config.CAFile)
			if err != nil {
				return nil, fmt.Errorf("erreur lecture CA: %w", err)
			}

			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				return nil, fmt.Errorf("erreur parsing CA: certificat invalide")
			}

			tlsConfig.RootCAs = caCertPool
		}
	}

	return tlsConfig, nil
}

// DefaultServerConfig retourne une configuration serveur par défaut
func DefaultServerConfig(certFile, keyFile string) *ServerConfig {
	return &ServerConfig{
		CertFile:                 certFile,
		KeyFile:                  keyFile,
		MinVersion:               MinTLSVersion,
		CipherSuites:             DefaultCipherSuites(),
		PreferServerCipherSuites: true,
	}
}

// DefaultClientConfig retourne une configuration client par défaut
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		MinVersion:         MinTLSVersion,
		InsecureSkipVerify: false,
	}
}
