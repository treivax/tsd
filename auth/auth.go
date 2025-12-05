// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// AuthTypeNone indique qu'aucune authentification n'est requise
	AuthTypeNone = "none"

	// AuthTypeKey indique une authentification par clé API
	AuthTypeKey = "key"

	// AuthTypeJWT indique une authentification par JWT
	AuthTypeJWT = "jwt"

	// DefaultTokenExpiration est la durée de validité par défaut d'un JWT (24h)
	DefaultTokenExpiration = 24 * time.Hour

	// MinKeyLength est la longueur minimale d'une clé API
	MinKeyLength = 32
)

var (
	// ErrInvalidToken indique que le token est invalide
	ErrInvalidToken = errors.New("token invalide")

	// ErrExpiredToken indique que le token a expiré
	ErrExpiredToken = errors.New("token expiré")

	// ErrUnauthorized indique que l'authentification a échoué
	ErrUnauthorized = errors.New("non autorisé")

	// ErrInvalidAuthType indique que le type d'authentification est invalide
	ErrInvalidAuthType = errors.New("type d'authentification invalide")
)

// Config contient la configuration d'authentification
type Config struct {
	// Type est le type d'authentification (none, key, jwt)
	Type string

	// AuthKeys est la liste des clés API valides (pour AuthTypeKey)
	AuthKeys []string

	// JWTSecret est le secret pour signer/vérifier les JWT (pour AuthTypeJWT)
	JWTSecret string

	// JWTExpiration est la durée de validité des JWT (pour AuthTypeJWT)
	JWTExpiration time.Duration

	// JWTIssuer est l'émetteur des JWT (pour AuthTypeJWT)
	JWTIssuer string
}

// Manager gère l'authentification
type Manager struct {
	config *Config
}

// Claims représente les claims d'un JWT
type Claims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// NewManager crée un nouveau gestionnaire d'authentification
func NewManager(config *Config) (*Manager, error) {
	if config == nil {
		return nil, errors.New("config ne peut pas être nil")
	}

	// Valider la configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// Définir les valeurs par défaut
	if config.JWTExpiration == 0 {
		config.JWTExpiration = DefaultTokenExpiration
	}

	if config.JWTIssuer == "" {
		config.JWTIssuer = "tsd-server"
	}

	return &Manager{
		config: config,
	}, nil
}

// validateConfig valide la configuration d'authentification
func validateConfig(config *Config) error {
	switch config.Type {
	case AuthTypeNone:
		// Aucune validation nécessaire
		return nil

	case AuthTypeKey:
		if len(config.AuthKeys) == 0 {
			return errors.New("au moins une clé API doit être configurée")
		}
		for i, key := range config.AuthKeys {
			if len(key) < MinKeyLength {
				return fmt.Errorf("la clé API %d est trop courte (min %d caractères)", i, MinKeyLength)
			}
		}
		return nil

	case AuthTypeJWT:
		if config.JWTSecret == "" {
			return errors.New("le secret JWT doit être configuré")
		}
		if len(config.JWTSecret) < MinKeyLength {
			return fmt.Errorf("le secret JWT est trop court (min %d caractères)", MinKeyLength)
		}
		return nil

	default:
		return fmt.Errorf("%w: %s (valeurs autorisées: none, key, jwt)", ErrInvalidAuthType, config.Type)
	}
}

// IsEnabled retourne vrai si l'authentification est activée
func (m *Manager) IsEnabled() bool {
	return m.config.Type != AuthTypeNone
}

// GetAuthType retourne le type d'authentification configuré
func (m *Manager) GetAuthType() string {
	return m.config.Type
}

// ValidateToken valide un token d'authentification
func (m *Manager) ValidateToken(token string) error {
	if !m.IsEnabled() {
		return nil
	}

	if token == "" {
		return ErrUnauthorized
	}

	switch m.config.Type {
	case AuthTypeKey:
		return m.validateAuthKey(token)

	case AuthTypeJWT:
		_, err := m.validateJWT(token)
		return err

	default:
		return ErrInvalidAuthType
	}
}

// validateAuthKey valide une clé API
func (m *Manager) validateAuthKey(key string) error {
	for _, validKey := range m.config.AuthKeys {
		// Utiliser subtle.ConstantTimeCompare pour éviter les timing attacks
		if subtle.ConstantTimeCompare([]byte(key), []byte(validKey)) == 1 {
			return nil
		}
	}
	return ErrUnauthorized
}

// validateJWT valide un JWT et retourne les claims
func (m *Manager) validateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Vérifier la méthode de signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
		}
		return []byte(m.config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Vérifier l'émetteur si configuré
	if m.config.JWTIssuer != "" && claims.Issuer != m.config.JWTIssuer {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateJWT génère un nouveau JWT
func (m *Manager) GenerateJWT(username string, roles []string) (string, error) {
	if m.config.Type != AuthTypeJWT {
		return "", errors.New("la génération de JWT n'est disponible qu'en mode JWT")
	}

	now := time.Now()
	claims := Claims{
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.config.JWTIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("erreur génération JWT: %w", err)
	}

	return tokenString, nil
}

// GenerateAuthKey génère une nouvelle clé API aléatoire
func GenerateAuthKey() (string, error) {
	// Générer 32 bytes aléatoires (256 bits)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("erreur génération clé: %w", err)
	}

	// Encoder en base64 URL-safe
	key := base64.URLEncoding.EncodeToString(b)
	return key, nil
}

// ExtractTokenFromHeader extrait le token d'un header Authorization
func ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Format attendu: "Bearer <token>" ou juste "<token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 {
		if strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Si pas de préfixe "Bearer", retourner le header complet
	return authHeader
}

// TokenInfo contient les informations d'un token validé
type TokenInfo struct {
	Type     string   // "key" ou "jwt"
	Username string   // Pour JWT uniquement
	Roles    []string // Pour JWT uniquement
	Valid    bool
}

// GetTokenInfo retourne les informations d'un token
func (m *Manager) GetTokenInfo(token string) (*TokenInfo, error) {
	info := &TokenInfo{
		Type:  m.config.Type,
		Valid: false,
	}

	if !m.IsEnabled() {
		info.Valid = true
		return info, nil
	}

	switch m.config.Type {
	case AuthTypeKey:
		if err := m.validateAuthKey(token); err == nil {
			info.Valid = true
		}
		return info, nil

	case AuthTypeJWT:
		claims, err := m.validateJWT(token)
		if err != nil {
			return info, err
		}
		info.Valid = true
		info.Username = claims.Username
		info.Roles = claims.Roles
		return info, nil

	default:
		return info, ErrInvalidAuthType
	}
}
