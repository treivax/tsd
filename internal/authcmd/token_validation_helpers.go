// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/treivax/tsd/auth"
)

// token_validation_helpers.go contient des fonctions helper pour la validation de tokens.
// Ces fonctions ont été extraites de validateToken() pour réduire la complexité.

// ValidationConfig contient les paramètres de validation de token
type ValidationConfig struct {
	Token       string
	AuthType    string
	Secret      string
	Keys        string
	Format      string
	Interactive bool
}

// parseValidationFlags parse les arguments de ligne de commande pour validation
func parseValidationFlags(args []string, stderr io.Writer) (*ValidationConfig, *flag.FlagSet, error) {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	config := &ValidationConfig{}
	token := fs.String("token", "", "Token à valider (requis)")
	authType := fs.String("type", "", "Type d'auth: key ou jwt (requis)")
	secret := fs.String("secret", "", "Secret JWT (requis si type=jwt)")
	keys := fs.String("keys", "", "Clés API valides séparées par des virgules (requis si type=key)")
	interactive := fs.Bool("i", false, "Mode interactif")
	format := fs.String("format", "text", "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		return nil, fs, err
	}

	config.Token = *token
	config.AuthType = *authType
	config.Secret = *secret
	config.Keys = *keys
	config.Interactive = *interactive
	config.Format = *format

	return config, fs, nil
}

// readInteractiveInput lit les inputs manquants en mode interactif
func readInteractiveInput(config *ValidationConfig, stdin io.Reader, stdout, stderr io.Writer) error {
	reader := bufio.NewReader(stdin)

	// Lire le token si manquant
	if config.Token == "" {
		fmt.Fprint(stdout, "Token: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("erreur lecture token: %w", err)
		}
		config.Token = strings.TrimSpace(input)
	}

	// Lire le type d'auth si manquant
	if config.AuthType == "" {
		fmt.Fprint(stdout, "Type d'authentification (key/jwt): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("erreur lecture type: %w", err)
		}
		config.AuthType = strings.TrimSpace(input)
	}

	// Lire le secret JWT si nécessaire
	if config.AuthType == "jwt" && config.Secret == "" {
		fmt.Fprint(stdout, "Secret JWT: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("erreur lecture secret: %w", err)
		}
		config.Secret = strings.TrimSpace(input)
	}

	// Lire les clés API si nécessaire
	if config.AuthType == "key" && config.Keys == "" {
		fmt.Fprint(stdout, "Clés API (séparées par des virgules): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("erreur lecture clés: %w", err)
		}
		config.Keys = strings.TrimSpace(input)
	}

	return nil
}

// validateConfigParameters valide que tous les paramètres requis sont présents
func validateConfigParameters(config *ValidationConfig) error {
	if config.Token == "" {
		return fmt.Errorf("le token est requis (-token)")
	}

	if config.AuthType == "" {
		return fmt.Errorf("le type d'authentification est requis (-type key|jwt)")
	}

	if config.AuthType == "key" && config.Keys == "" {
		return fmt.Errorf("les clés API sont requises pour type=key (-keys)")
	}

	if config.AuthType == "jwt" && config.Secret == "" {
		return fmt.Errorf("le secret JWT est requis pour type=jwt (-secret)")
	}

	if config.AuthType != "key" && config.AuthType != "jwt" {
		return fmt.Errorf("type invalide '%s' (doit être 'key' ou 'jwt')", config.AuthType)
	}

	return nil
}

// createAuthConfig crée une configuration auth.Config depuis ValidationConfig
func createAuthConfig(config *ValidationConfig) (*auth.Config, error) {
	switch config.AuthType {
	case "key":
		keysList := strings.Split(config.Keys, ",")
		for i, key := range keysList {
			keysList[i] = strings.TrimSpace(key)
		}
		return &auth.Config{
			Type:     auth.AuthTypeKey,
			AuthKeys: keysList,
		}, nil

	case "jwt":
		return &auth.Config{
			Type:      auth.AuthTypeJWT,
			JWTSecret: config.Secret,
		}, nil

	default:
		return nil, fmt.Errorf("type invalide '%s'", config.AuthType)
	}
}

// ValidationResult contient le résultat d'une validation de token
type ValidationResult struct {
	Valid    bool
	Username string
	Roles    []string
	Error    error
}

// validateTokenWithManager valide un token en utilisant un auth.Manager
func validateTokenWithManager(manager *auth.Manager, token string) *ValidationResult {
	info, err := manager.GetTokenInfo(token)

	result := &ValidationResult{
		Valid: err == nil && info != nil && info.Valid,
		Error: err,
	}

	if info != nil {
		result.Username = info.Username
		result.Roles = info.Roles
	}

	return result
}

// formatValidationOutput formate le résultat de validation selon le format demandé
func formatValidationOutput(result *ValidationResult, config *ValidationConfig) string {
	if config.Format == "json" {
		return formatJSONOutput(result, config.AuthType)
	}
	return formatTextOutput(result, config.AuthType)
}

// formatJSONOutput formate la sortie en JSON
func formatJSONOutput(result *ValidationResult, authType string) string {
	output := map[string]interface{}{
		"valid": result.Valid,
		"type":  authType,
	}

	if result.Error != nil {
		output["error"] = result.Error.Error()
	}

	if result.Username != "" {
		output["username"] = result.Username
	}

	if len(result.Roles) > 0 {
		output["roles"] = result.Roles
	}

	// Note: json.MarshalIndent toujours réussit pour ces types simples
	data, _ := marshalJSON(output)
	return string(data)
}

// formatTextOutput formate la sortie en texte
func formatTextOutput(result *ValidationResult, authType string) string {
	var builder strings.Builder

	if result.Valid {
		builder.WriteString("✅ Token valide\n")
		builder.WriteString(fmt.Sprintf("Type: %s\n", authType))

		if authType == "jwt" && result.Username != "" {
			builder.WriteString(fmt.Sprintf("Utilisateur: %s\n", result.Username))
			if len(result.Roles) > 0 {
				builder.WriteString(fmt.Sprintf("Rôles: %s\n", strings.Join(result.Roles, ", ")))
			}
		}
	} else {
		builder.WriteString("❌ Token invalide\n")
		if result.Error != nil {
			builder.WriteString(fmt.Sprintf("Erreur: %v\n", result.Error))
		}
	}

	return builder.String()
}

// marshalJSON est un wrapper pour json.MarshalIndent pour faciliter les tests
var marshalJSON = func(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
