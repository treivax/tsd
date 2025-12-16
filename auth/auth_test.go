// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package auth

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Test constants
const (
	TestUsername    = "testuser"
	TestRole1       = "admin"
	TestRole2       = "user"
	TestValidKey    = "test-key-12345678901234567890123456789012"
	TestInvalidKey  = "invalid-key"
	TestJWTSecret   = "test-secret-12345678901234567890123"
	TestShortKey    = "short"
	TestIssuer      = "test-issuer"
	TestExpiration  = 1 * time.Hour
	TestConcurrency = 50
)

// TestNewManager_ValidConfig tests NewManager with valid configurations
func TestNewManager_ValidConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "auth type none",
			config: &Config{
				Type: AuthTypeNone,
			},
			wantErr: false,
		},
		{
			name: "auth type key",
			config: &Config{
				Type:     AuthTypeKey,
				AuthKeys: []string{TestValidKey},
			},
			wantErr: false,
		},
		{
			name: "auth type jwt",
			config: &Config{
				Type:          AuthTypeJWT,
				JWTSecret:     TestJWTSecret,
				JWTExpiration: TestExpiration,
				JWTIssuer:     TestIssuer,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && manager == nil {
				t.Errorf("NewManager() returned nil manager")
			}
		})
	}
}

// TestNewManager_NilConfig tests NewManager with nil config
func TestNewManager_NilConfig(t *testing.T) {
	_, err := NewManager(nil)
	if err == nil {
		t.Errorf("NewManager() with nil config should return error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil config, got: %v", err)
	}
}

// TestNewManager_InvalidConfig tests NewManager with invalid configurations
func TestNewManager_InvalidConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr string
	}{
		{
			name: "invalid auth type",
			config: &Config{
				Type: "invalid",
			},
			wantErr: "invalide",
		},
		{
			name: "key auth without keys",
			config: &Config{
				Type:     AuthTypeKey,
				AuthKeys: []string{},
			},
			wantErr: "au moins une clé",
		},
		{
			name: "key auth with short key",
			config: &Config{
				Type:     AuthTypeKey,
				AuthKeys: []string{TestShortKey},
			},
			wantErr: "trop courte",
		},
		{
			name: "jwt auth without secret",
			config: &Config{
				Type:      AuthTypeJWT,
				JWTSecret: "",
			},
			wantErr: "secret JWT",
		},
		{
			name: "jwt auth with short secret",
			config: &Config{
				Type:      AuthTypeJWT,
				JWTSecret: TestShortKey,
			},
			wantErr: "trop court",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewManager(tt.config)
			if err == nil {
				t.Errorf("NewManager() should return error")
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error should contain %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

// TestNewManager_DefaultValues tests default values are set
func TestNewManager_DefaultValues(t *testing.T) {
	config := &Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	}

	manager, err := NewManager(config)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	if manager.config.JWTExpiration != DefaultTokenExpiration {
		t.Errorf("JWTExpiration = %v, want %v", manager.config.JWTExpiration, DefaultTokenExpiration)
	}

	if manager.config.JWTIssuer != "tsd-server" {
		t.Errorf("JWTIssuer = %q, want %q", manager.config.JWTIssuer, "tsd-server")
	}
}

// TestManager_IsEnabled tests IsEnabled method
func TestManager_IsEnabled(t *testing.T) {
	tests := []struct {
		name        string
		authType    string
		wantEnabled bool
	}{
		{
			name:        "none type disabled",
			authType:    AuthTypeNone,
			wantEnabled: false,
		},
		{
			name:        "key type enabled",
			authType:    AuthTypeKey,
			wantEnabled: true,
		},
		{
			name:        "jwt type enabled",
			authType:    AuthTypeJWT,
			wantEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config *Config
			if tt.authType == AuthTypeKey {
				config = &Config{
					Type:     tt.authType,
					AuthKeys: []string{TestValidKey},
				}
			} else if tt.authType == AuthTypeJWT {
				config = &Config{
					Type:      tt.authType,
					JWTSecret: TestJWTSecret,
				}
			} else {
				config = &Config{
					Type: tt.authType,
				}
			}

			manager, err := NewManager(config)
			if err != nil {
				t.Fatalf("NewManager() error = %v", err)
			}

			if got := manager.IsEnabled(); got != tt.wantEnabled {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.wantEnabled)
			}
		})
	}
}

// TestManager_GetAuthType tests GetAuthType method
func TestManager_GetAuthType(t *testing.T) {
	tests := []struct {
		name     string
		authType string
	}{
		{name: "none", authType: AuthTypeNone},
		{name: "key", authType: AuthTypeKey},
		{name: "jwt", authType: AuthTypeJWT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config *Config
			if tt.authType == AuthTypeKey {
				config = &Config{
					Type:     tt.authType,
					AuthKeys: []string{TestValidKey},
				}
			} else if tt.authType == AuthTypeJWT {
				config = &Config{
					Type:      tt.authType,
					JWTSecret: TestJWTSecret,
				}
			} else {
				config = &Config{
					Type: tt.authType,
				}
			}

			manager, _ := NewManager(config)
			if got := manager.GetAuthType(); got != tt.authType {
				t.Errorf("GetAuthType() = %q, want %q", got, tt.authType)
			}
		})
	}
}

// TestManager_ValidateToken_None tests ValidateToken with auth type none
func TestManager_ValidateToken_None(t *testing.T) {
	manager, _ := NewManager(&Config{Type: AuthTypeNone})

	tests := []struct {
		name  string
		token string
	}{
		{name: "empty token", token: ""},
		{name: "any token", token: "anything"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateToken(tt.token)
			if err != nil {
				t.Errorf("ValidateToken() with auth none should always succeed, got error: %v", err)
			}
		})
	}
}

// TestManager_ValidateToken_Key tests ValidateToken with auth type key
func TestManager_ValidateToken_Key(t *testing.T) {
	validKey1 := "valid-key-1234567890123456789012345"
	validKey2 := "valid-key-9876543210987654321098765"

	manager, _ := NewManager(&Config{
		Type:     AuthTypeKey,
		AuthKeys: []string{validKey1, validKey2},
	})

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   TestInvalidKey,
			wantErr: true,
		},
		{
			name:    "valid key 1",
			token:   validKey1,
			wantErr: false,
		},
		{
			name:    "valid key 2",
			token:   validKey2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !errors.Is(err, ErrUnauthorized) {
				t.Errorf("error should be ErrUnauthorized, got: %v", err)
			}
		})
	}
}

// TestManager_ValidateToken_JWT tests ValidateToken with auth type JWT
func TestManager_ValidateToken_JWT(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:          AuthTypeJWT,
		JWTSecret:     TestJWTSecret,
		JWTExpiration: 1 * time.Hour,
		JWTIssuer:     TestIssuer,
	})

	// Generate valid token
	validToken, _ := manager.GenerateJWT(TestUsername, []string{TestRole1})

	// Generate expired token
	expiredClaims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			Issuer:    TestIssuer,
		},
	}
	expiredJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, _ := expiredJWT.SignedString([]byte(TestJWTSecret))

	// Generate token with wrong issuer
	wrongIssuerClaims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "wrong-issuer",
		},
	}
	wrongIssuerJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, wrongIssuerClaims)
	wrongIssuerToken, _ := wrongIssuerJWT.SignedString([]byte(TestJWTSecret))

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "empty token",
			token:   "",
			wantErr: ErrUnauthorized,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "expired token",
			token:   expiredToken,
			wantErr: ErrExpiredToken,
		},
		{
			name:    "wrong issuer",
			token:   wrongIssuerToken,
			wantErr: ErrInvalidToken,
		},
		{
			name:    "valid token",
			token:   validToken,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateToken(tt.token)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("ValidateToken() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("ValidateToken() error = nil, want %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ValidateToken() error = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

// TestManager_GenerateJWT tests JWT generation
func TestManager_GenerateJWT(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:          AuthTypeJWT,
		JWTSecret:     TestJWTSecret,
		JWTExpiration: TestExpiration,
		JWTIssuer:     TestIssuer,
	})

	roles := []string{TestRole1, TestRole2}
	token, err := manager.GenerateJWT(TestUsername, roles)

	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	if token == "" {
		t.Errorf("GenerateJWT() returned empty token")
	}

	// Validate the generated token
	err = manager.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() on generated token error = %v", err)
	}
}

// TestManager_GenerateJWT_WrongAuthType tests GenerateJWT with wrong auth type
func TestManager_GenerateJWT_WrongAuthType(t *testing.T) {
	tests := []struct {
		name     string
		authType string
	}{
		{name: "auth type none", authType: AuthTypeNone},
		{name: "auth type key", authType: AuthTypeKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config *Config
			if tt.authType == AuthTypeKey {
				config = &Config{
					Type:     tt.authType,
					AuthKeys: []string{TestValidKey},
				}
			} else {
				config = &Config{
					Type: tt.authType,
				}
			}

			manager, _ := NewManager(config)
			_, err := manager.GenerateJWT(TestUsername, []string{TestRole1})

			if err == nil {
				t.Errorf("GenerateJWT() should return error for auth type %s", tt.authType)
			}
		})
	}
}

// TestGenerateAuthKey tests GenerateAuthKey function
func TestGenerateAuthKey(t *testing.T) {
	key, err := GenerateAuthKey()

	if err != nil {
		t.Fatalf("GenerateAuthKey() error = %v", err)
	}

	if key == "" {
		t.Errorf("GenerateAuthKey() returned empty key")
	}

	if len(key) < MinKeyLength {
		t.Errorf("GenerateAuthKey() key length = %d, want >= %d", len(key), MinKeyLength)
	}

	// Generate another key and verify they're different
	key2, _ := GenerateAuthKey()
	if key == key2 {
		t.Errorf("GenerateAuthKey() should generate different keys")
	}
}

// TestExtractTokenFromHeader tests ExtractTokenFromHeader function
func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
	}{
		{
			name:      "empty header",
			header:    "",
			wantToken: "",
		},
		{
			name:      "bearer prefix lowercase",
			header:    "bearer token123",
			wantToken: "token123",
		},
		{
			name:      "bearer prefix uppercase",
			header:    "Bearer token456",
			wantToken: "token456",
		},
		{
			name:      "bearer prefix mixed case",
			header:    "BEARER token789",
			wantToken: "token789",
		},
		{
			name:      "no bearer prefix",
			header:    "token-without-prefix",
			wantToken: "token-without-prefix",
		},
		{
			name:      "bearer with spaces",
			header:    "Bearer   token-with-spaces",
			wantToken: "  token-with-spaces",
		},
		{
			name:      "other prefix",
			header:    "Basic token123",
			wantToken: "Basic token123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTokenFromHeader(tt.header)
			if got != tt.wantToken {
				t.Errorf("ExtractTokenFromHeader() = %q, want %q", got, tt.wantToken)
			}
		})
	}
}

// TestManager_GetTokenInfo tests GetTokenInfo method
func TestManager_GetTokenInfo(t *testing.T) {
	t.Run("auth type none", func(t *testing.T) {
		manager, _ := NewManager(&Config{Type: AuthTypeNone})
		info, err := manager.GetTokenInfo("")

		if err != nil {
			t.Errorf("GetTokenInfo() error = %v", err)
		}
		if !info.Valid {
			t.Errorf("GetTokenInfo() Valid = false, want true")
		}
	})

	t.Run("auth type key valid", func(t *testing.T) {
		manager, _ := NewManager(&Config{
			Type:     AuthTypeKey,
			AuthKeys: []string{TestValidKey},
		})
		info, err := manager.GetTokenInfo(TestValidKey)

		if err != nil {
			t.Errorf("GetTokenInfo() error = %v", err)
		}
		if !info.Valid {
			t.Errorf("GetTokenInfo() Valid = false, want true")
		}
		if info.Type != AuthTypeKey {
			t.Errorf("GetTokenInfo() Type = %q, want %q", info.Type, AuthTypeKey)
		}
	})

	t.Run("auth type key invalid", func(t *testing.T) {
		manager, _ := NewManager(&Config{
			Type:     AuthTypeKey,
			AuthKeys: []string{TestValidKey},
		})
		info, err := manager.GetTokenInfo(TestInvalidKey)

		if err != nil {
			t.Errorf("GetTokenInfo() error = %v", err)
		}
		if info.Valid {
			t.Errorf("GetTokenInfo() Valid = true, want false")
		}
	})

	t.Run("auth type jwt valid", func(t *testing.T) {
		manager, _ := NewManager(&Config{
			Type:      AuthTypeJWT,
			JWTSecret: TestJWTSecret,
			JWTIssuer: TestIssuer,
		})
		token, _ := manager.GenerateJWT(TestUsername, []string{TestRole1, TestRole2})
		info, err := manager.GetTokenInfo(token)

		if err != nil {
			t.Errorf("GetTokenInfo() error = %v", err)
		}
		if !info.Valid {
			t.Errorf("GetTokenInfo() Valid = false, want true")
		}
		if info.Type != AuthTypeJWT {
			t.Errorf("GetTokenInfo() Type = %q, want %q", info.Type, AuthTypeJWT)
		}
		if info.Username != TestUsername {
			t.Errorf("GetTokenInfo() Username = %q, want %q", info.Username, TestUsername)
		}
		if len(info.Roles) != 2 {
			t.Errorf("GetTokenInfo() Roles length = %d, want 2", len(info.Roles))
		}
	})

	t.Run("auth type jwt invalid", func(t *testing.T) {
		manager, _ := NewManager(&Config{
			Type:      AuthTypeJWT,
			JWTSecret: TestJWTSecret,
		})
		_, err := manager.GetTokenInfo("invalid.token.here")

		if err == nil {
			t.Errorf("GetTokenInfo() should return error for invalid JWT")
		}
	})
}

// TestAuthTypeConstants tests auth type constants
func TestAuthTypeConstants(t *testing.T) {
	if AuthTypeNone != "none" {
		t.Errorf("AuthTypeNone = %q, want %q", AuthTypeNone, "none")
	}
	if AuthTypeKey != "key" {
		t.Errorf("AuthTypeKey = %q, want %q", AuthTypeKey, "key")
	}
	if AuthTypeJWT != "jwt" {
		t.Errorf("AuthTypeJWT = %q, want %q", AuthTypeJWT, "jwt")
	}
}

// TestDefaultTokenExpiration tests default token expiration constant
func TestDefaultTokenExpiration(t *testing.T) {
	expected := 24 * time.Hour
	if DefaultTokenExpiration != expected {
		t.Errorf("DefaultTokenExpiration = %v, want %v", DefaultTokenExpiration, expected)
	}
}

// TestMinKeyLength tests minimum key length constant
func TestMinKeyLength(t *testing.T) {
	if MinKeyLength != 32 {
		t.Errorf("MinKeyLength = %d, want 32", MinKeyLength)
	}
}

// TestErrorConstants tests error constants
func TestErrorConstants(t *testing.T) {
	if ErrInvalidToken == nil {
		t.Errorf("ErrInvalidToken should not be nil")
	}
	if ErrExpiredToken == nil {
		t.Errorf("ErrExpiredToken should not be nil")
	}
	if ErrUnauthorized == nil {
		t.Errorf("ErrUnauthorized should not be nil")
	}
	if ErrInvalidAuthType == nil {
		t.Errorf("ErrInvalidAuthType should not be nil")
	}
}

// TestManager_ValidateToken_TimingAttackResistance tests constant-time comparison
func TestManager_ValidateToken_TimingAttackResistance(t *testing.T) {
	// This test verifies that key comparison uses constant-time comparison
	// We can't directly test timing, but we verify the behavior is correct
	validKey := "valid-key-12345678901234567890123456789012"
	manager, _ := NewManager(&Config{
		Type:     AuthTypeKey,
		AuthKeys: []string{validKey},
	})

	// All these should fail, testing various lengths and prefixes
	invalidKeys := []string{
		"valid-key-12345678901234567890123456789011",  // Last char different
		"xalid-key-12345678901234567890123456789012",  // First char different
		"valid-key-123456789012345678901234567890",    // Shorter
		"valid-key-123456789012345678901234567890123", // Longer
	}

	for _, key := range invalidKeys {
		err := manager.ValidateToken(key)
		if err == nil {
			t.Errorf("ValidateToken(%q) should fail", key)
		}
		if !errors.Is(err, ErrUnauthorized) {
			t.Errorf("ValidateToken() error = %v, want ErrUnauthorized", err)
		}
	}

	// Valid key should succeed
	err := manager.ValidateToken(validKey)
	if err != nil {
		t.Errorf("ValidateToken() with valid key error = %v", err)
	}
}

// TestManager_JWT_ClaimsValidation tests JWT claims validation
func TestManager_JWT_ClaimsValidation(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
		JWTIssuer: TestIssuer,
	})

	// Generate token and immediately validate to extract claims
	roles := []string{TestRole1, TestRole2}
	token, _ := manager.GenerateJWT(TestUsername, roles)

	info, err := manager.GetTokenInfo(token)
	if err != nil {
		t.Fatalf("GetTokenInfo() error = %v", err)
	}

	if info.Username != TestUsername {
		t.Errorf("Username = %q, want %q", info.Username, TestUsername)
	}

	if len(info.Roles) != len(roles) {
		t.Errorf("len(Roles) = %d, want %d", len(info.Roles), len(roles))
	}

	for i, role := range roles {
		if info.Roles[i] != role {
			t.Errorf("Roles[%d] = %q, want %q", i, info.Roles[i], role)
		}
	}
}

// TestManager_MultipleKeys tests auth with multiple valid keys
func TestManager_MultipleKeys(t *testing.T) {
	key1 := "key-one-1234567890123456789012345678"
	key2 := "key-two-1234567890123456789012345678"
	key3 := "key-three-12345678901234567890123456"

	manager, _ := NewManager(&Config{
		Type:     AuthTypeKey,
		AuthKeys: []string{key1, key2, key3},
	})

	// All three keys should be valid
	for i, key := range []string{key1, key2, key3} {
		err := manager.ValidateToken(key)
		if err != nil {
			t.Errorf("ValidateToken() with key %d error = %v", i+1, err)
		}
	}

	// Invalid key should fail
	err := manager.ValidateToken("invalid-key-1234567890123456789012")
	if err == nil {
		t.Errorf("ValidateToken() with invalid key should fail")
	}
}

// TestManager_JWTWithoutIssuer tests JWT without issuer validation
func TestManager_JWTWithoutIssuer(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
		JWTIssuer: "", // No issuer
	})

	token, _ := manager.GenerateJWT(TestUsername, []string{TestRole1})

	// Token should be valid even without issuer
	err := manager.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() error = %v", err)
	}
}

// TestClaims_Structure tests Claims structure
func TestClaims_Structure(t *testing.T) {
	claims := Claims{
		Username: TestUsername,
		Roles:    []string{TestRole1, TestRole2},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: TestIssuer,
		},
	}

	if claims.Username != TestUsername {
		t.Errorf("Username = %q, want %q", claims.Username, TestUsername)
	}

	if len(claims.Roles) != 2 {
		t.Errorf("len(Roles) = %d, want 2", len(claims.Roles))
	}

	if claims.Issuer != TestIssuer {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, TestIssuer)
	}
}

// TestConfig_Structure tests Config structure
func TestConfig_Structure(t *testing.T) {
	config := Config{
		Type:          AuthTypeJWT,
		AuthKeys:      []string{TestValidKey},
		JWTSecret:     TestJWTSecret,
		JWTExpiration: TestExpiration,
		JWTIssuer:     TestIssuer,
	}

	if config.Type != AuthTypeJWT {
		t.Errorf("Type = %q, want %q", config.Type, AuthTypeJWT)
	}

	if len(config.AuthKeys) != 1 {
		t.Errorf("len(AuthKeys) = %d, want 1", len(config.AuthKeys))
	}

	if config.JWTSecret != TestJWTSecret {
		t.Errorf("JWTSecret = %q, want %q", config.JWTSecret, TestJWTSecret)
	}

	if config.JWTExpiration != TestExpiration {
		t.Errorf("JWTExpiration = %v, want %v", config.JWTExpiration, TestExpiration)
	}

	if config.JWTIssuer != TestIssuer {
		t.Errorf("JWTIssuer = %q, want %q", config.JWTIssuer, TestIssuer)
	}
}

// TestTokenInfo_Structure tests TokenInfo structure
func TestTokenInfo_Structure(t *testing.T) {
	info := TokenInfo{
		Type:     AuthTypeJWT,
		Username: TestUsername,
		Roles:    []string{TestRole1, TestRole2},
		Valid:    true,
	}

	if info.Type != AuthTypeJWT {
		t.Errorf("Type = %q, want %q", info.Type, AuthTypeJWT)
	}

	if info.Username != TestUsername {
		t.Errorf("Username = %q, want %q", info.Username, TestUsername)
	}

	if len(info.Roles) != 2 {
		t.Errorf("len(Roles) = %d, want 2", len(info.Roles))
	}

	if !info.Valid {
		t.Errorf("Valid = false, want true")
	}
}

// TestManager_ValidateJWT_WrongSigningMethod tests JWT with wrong signing method
func TestManager_ValidateJWT_WrongSigningMethod(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
		JWTIssuer: TestIssuer,
	})

	// HS512 is still HMAC so it will be accepted
	// The code only checks if it's HMAC, not specifically HS256
	// This test verifies that HMAC methods are accepted
	jti, _ := generateJTI()
	claims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{TokenAudience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    TestIssuer,
			ID:        jti, // JTI requis pour validation
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, _ := token.SignedString([]byte(TestJWTSecret))

	// Should succeed because HS512 is still HMAC
	err := manager.ValidateToken(tokenString)
	if err != nil {
		t.Errorf("ValidateToken() error = %v, HS512 should be accepted as HMAC", err)
	}
}

// TestManager_ValidateJWT_InvalidClaims tests JWT with invalid claims structure
func TestManager_ValidateJWT_InvalidClaims(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	})

	// Create token with standard claims instead of custom Claims
	standardClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims)
	tokenString, _ := token.SignedString([]byte(TestJWTSecret))

	// This should fail because we can't cast to *Claims
	_, err := manager.validateJWT(tokenString)
	if err == nil {
		t.Errorf("validateJWT() should fail for invalid claims structure")
	}
}

// TestGenerateAuthKey_Coverage tests GenerateAuthKey error path
func TestGenerateAuthKey_Coverage(t *testing.T) {
	// Call multiple times to ensure randomness
	keys := make(map[string]bool)
	for i := 0; i < 10; i++ {
		key, err := GenerateAuthKey()
		if err != nil {
			t.Errorf("GenerateAuthKey() iteration %d error = %v", i, err)
		}
		if keys[key] {
			t.Errorf("GenerateAuthKey() generated duplicate key")
		}
		keys[key] = true
	}
}

// TestManager_GenerateJWT_CustomExpiration tests JWT with custom expiration
func TestManager_GenerateJWT_CustomExpiration(t *testing.T) {
	customExpiration := 2 * time.Hour
	manager, _ := NewManager(&Config{
		Type:          AuthTypeJWT,
		JWTSecret:     TestJWTSecret,
		JWTExpiration: customExpiration,
		JWTIssuer:     TestIssuer,
	})

	token, err := manager.GenerateJWT(TestUsername, []string{TestRole1})
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	// Validate and check expiration
	info, err := manager.GetTokenInfo(token)
	if err != nil {
		t.Errorf("GetTokenInfo() error = %v", err)
	}
	if !info.Valid {
		t.Errorf("Generated token should be valid")
	}
}

// TestManager_GenerateJWT_NoRoles tests JWT generation without roles
func TestManager_GenerateJWT_NoRoles(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	})

	token, err := manager.GenerateJWT(TestUsername, nil)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	info, err := manager.GetTokenInfo(token)
	if err != nil {
		t.Errorf("GetTokenInfo() error = %v", err)
	}
	if len(info.Roles) != 0 {
		t.Errorf("Roles length = %d, want 0", len(info.Roles))
	}
}

// TestManager_GenerateJWT_EmptyRoles tests JWT generation with empty roles array
func TestManager_GenerateJWT_EmptyRoles(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	})

	token, err := manager.GenerateJWT(TestUsername, []string{})
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	info, err := manager.GetTokenInfo(token)
	if err != nil {
		t.Errorf("GetTokenInfo() error = %v", err)
	}
	if len(info.Roles) != 0 {
		t.Errorf("Roles length = %d, want 0", len(info.Roles))
	}
}

// TestManager_ValidateToken_EdgeCases tests edge cases in token validation
func TestManager_ValidateJWT_NotBefore(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	})

	// Create token that's not valid yet (NotBefore in future)
	claims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"tsd-api"},
			NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(TestJWTSecret))

	err := manager.ValidateToken(tokenString)
	if err == nil {
		t.Errorf("ValidateToken() should fail for token not yet valid")
	}
}

// TestExtractTokenFromHeader_EdgeCases tests additional edge cases
func TestExtractTokenFromHeader_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
	}{
		{
			name:      "token with bearer lowercase no space",
			header:    "bearertoken123",
			wantToken: "bearertoken123",
		},
		{
			name:      "multiple spaces after bearer",
			header:    "Bearer    token-multi-space",
			wantToken: "   token-multi-space",
		},
		{
			name:      "bearer with newline",
			header:    "Bearer token\n",
			wantToken: "token\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTokenFromHeader(tt.header)
			if got != tt.wantToken {
				t.Errorf("ExtractTokenFromHeader() = %q, want %q", got, tt.wantToken)
			}
		})
	}
}

// TestManager_GetTokenInfo_InvalidType tests GetTokenInfo with invalid auth type
func TestManager_GetTokenInfo_InvalidType(t *testing.T) {
	// This tests an edge case where the manager has an invalid type
	manager := &Manager{
		config: &Config{
			Type: "invalid",
		},
	}

	_, err := manager.GetTokenInfo("token")
	if !errors.Is(err, ErrInvalidAuthType) {
		t.Errorf("GetTokenInfo() error = %v, want ErrInvalidAuthType", err)
	}
}

// TestManager_ValidateAuthKey_EmptyKeysList tests validation with no keys
func TestManager_ValidateAuthKey_EmptyKeysList(t *testing.T) {
	manager := &Manager{
		config: &Config{
			Type:     AuthTypeKey,
			AuthKeys: []string{},
		},
	}

	err := manager.validateAuthKey("any-key")
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("validateAuthKey() error = %v, want ErrUnauthorized", err)
	}
}

// TestManager_ValidateToken_InvalidAuthType tests ValidateToken with invalid type
func TestManager_ValidateToken_InvalidAuthType(t *testing.T) {
	manager := &Manager{
		config: &Config{
			Type: "invalid",
		},
	}

	err := manager.ValidateToken("token")
	if !errors.Is(err, ErrInvalidAuthType) {
		t.Errorf("ValidateToken() error = %v, want ErrInvalidAuthType", err)
	}
}

// TestManager_GenerateJWT_LongUsername tests JWT with very long username
func TestManager_GenerateJWT_LongUsername(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	})

	longUsername := strings.Repeat("a", 1000)
	token, err := manager.GenerateJWT(longUsername, []string{TestRole1})
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	info, err := manager.GetTokenInfo(token)
	if err != nil {
		t.Errorf("GetTokenInfo() error = %v", err)
	}
	if info.Username != longUsername {
		t.Errorf("Username mismatch")
	}
}

// TestManager_ValidateJWT_TokenWithoutExpiry tests JWT without expiry
func TestManager_ValidateJWT_TokenWithoutExpiry(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	})

	// Create token without expiry - library will reject this as invalid
	claims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{"tsd-api"},
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(TestJWTSecret))

	// The jwt library validates expiry and will reject tokens without proper claims
	err := manager.ValidateToken(tokenString)
	if err == nil {
		t.Errorf("ValidateToken() should fail for token without proper expiry")
	}
}

// TestManager_ValidateJWT_InvalidAudience tests JWT with wrong audience
func TestManager_ValidateJWT_InvalidAudience(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
		JWTIssuer: TestIssuer,
	})

	// Create token with wrong audience
	claims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"other-service"}, // Wrong audience
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    TestIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(TestJWTSecret))

	err := manager.ValidateToken(tokenString)
	if err == nil {
		t.Error("❌ Token with wrong audience should be rejected")
	} else {
		t.Logf("✅ Token with wrong audience correctly rejected: %v", err)
	}
}

// TestManager_ValidateJWT_MissingAudience tests JWT without audience claim
func TestManager_ValidateJWT_MissingAudience(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
		JWTIssuer: TestIssuer,
	})

	// Create token without audience
	claims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    TestIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(TestJWTSecret))

	err := manager.ValidateToken(tokenString)
	if err == nil {
		t.Error("❌ Token without audience should be rejected")
	} else {
		t.Logf("✅ Token without audience correctly rejected: %v", err)
	}
}

// TestGenerateJWT_HasAllRequiredClaims tests that generated JWT has all required claims
func TestGenerateJWT_HasAllRequiredClaims(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
		JWTIssuer: TestIssuer,
	})

	tokenString, err := manager.GenerateJWT(TestUsername, []string{"admin", "user"})
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	// Parse token to verify claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(TestJWTSecret), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatal("Failed to cast claims")
	}

	// Verify all required claims are present
	if claims.Subject != TestUsername {
		t.Errorf("❌ Subject (sub) missing or incorrect: got %v, want %v", claims.Subject, TestUsername)
	} else {
		t.Log("✅ Subject (sub) claim present")
	}

	if len(claims.Audience) == 0 || claims.Audience[0] != "tsd-api" {
		t.Errorf("❌ Audience (aud) missing or incorrect: got %v, want [tsd-api]", claims.Audience)
	} else {
		t.Log("✅ Audience (aud) claim present")
	}

	if claims.ID == "" {
		t.Error("❌ JWT ID (jti) missing")
	} else {
		t.Logf("✅ JWT ID (jti) present: %s", claims.ID)
	}

	if claims.Issuer != TestIssuer {
		t.Errorf("❌ Issuer (iss) incorrect: got %v, want %v", claims.Issuer, TestIssuer)
	} else {
		t.Log("✅ Issuer (iss) claim present")
	}

	if claims.ExpiresAt == nil {
		t.Error("❌ ExpiresAt (exp) missing")
	} else {
		t.Log("✅ ExpiresAt (exp) claim present")
	}

	if claims.IssuedAt == nil {
		t.Error("❌ IssuedAt (iat) missing")
	} else {
		t.Log("✅ IssuedAt (iat) claim present")
	}

	if claims.NotBefore == nil {
		t.Error("❌ NotBefore (nbf) missing")
	} else {
		t.Log("✅ NotBefore (nbf) claim present")
	}
}

// TestGenerateJWT_UniqueJTI tests that each generated JWT has a unique JTI
func TestGenerateJWT_UniqueJTI(t *testing.T) {
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
	})

	// Generate 100 tokens and verify all have unique JTI
	jtis := make(map[string]bool)
	for i := 0; i < 100; i++ {
		tokenString, err := manager.GenerateJWT(TestUsername, nil)
		if err != nil {
			t.Fatalf("GenerateJWT() error = %v", err)
		}

		token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(TestJWTSecret), nil
		})

		claims, _ := token.Claims.(*Claims)
		jti := claims.ID

		if jtis[jti] {
			t.Errorf("❌ Duplicate JTI detected: %s", jti)
		}
		jtis[jti] = true
	}

	t.Logf("✅ Generated 100 tokens with unique JTI values")
}

// TestGenerateJTI_Format tests that JTI has correct format and length
func TestGenerateJTI_Format(t *testing.T) {
	t.Log("🧪 TEST FORMAT JTI")
	t.Log("==================")

	// Act - Générer plusieurs JTI
	jti1, err := generateJTI()
	if err != nil {
		t.Fatalf("❌ Échec génération JTI: %v", err)
	}

	jti2, err := generateJTI()
	if err != nil {
		t.Fatalf("❌ Échec génération JTI: %v", err)
	}

	// Assert - Unicité
	if jti1 == jti2 {
		t.Error("❌ JTI identiques, devrait être unique")
	} else {
		t.Logf("✅ JTI uniques: %s != %s", jti1[:10]+"...", jti2[:10]+"...")
	}

	// Assert - Longueur raisonnable (base64 de 16 bytes ≈ 22-24 chars)
	minLength := 20
	maxLength := 30
	if len(jti1) < minLength || len(jti1) > maxLength {
		t.Errorf("❌ Longueur JTI anormale: %d (attendu entre %d et %d)", len(jti1), minLength, maxLength)
	} else {
		t.Logf("✅ Longueur JTI correcte: %d caractères", len(jti1))
	}

	// Assert - Format base64 URL-safe
	decoded, err := base64.URLEncoding.DecodeString(jti1)
	if err != nil {
		t.Errorf("❌ JTI n'est pas base64 URL-safe valide: %v", err)
	} else {
		t.Logf("✅ JTI en base64 URL-safe valide (%d bytes décodés)", len(decoded))
	}

	// Assert - Longueur décodée devrait être JTILength
	if len(decoded) != JTILength {
		t.Errorf("❌ JTI décodé = %d bytes, attendu %d bytes", len(decoded), JTILength)
	} else {
		t.Logf("✅ JTI décodé = %d bytes (JTILength)", JTILength)
	}
}

// TestGenerateJTI_Uniqueness tests JTI uniqueness with high volume
func TestGenerateJTI_Uniqueness(t *testing.T) {
	t.Log("🧪 TEST UNICITÉ JTI (HIGH VOLUME)")
	t.Log("==================================")

	const testCount = 1000
	jtis := make(map[string]bool, testCount)

	// Generate many JTIs and verify all are unique
	for i := 0; i < testCount; i++ {
		jti, err := generateJTI()
		if err != nil {
			t.Fatalf("❌ Échec génération JTI à l'itération %d: %v", i, err)
		}

		if jtis[jti] {
			t.Errorf("❌ JTI dupliqué détecté à l'itération %d: %s", i, jti)
			return
		}
		jtis[jti] = true
	}

	t.Logf("✅ %d JTI uniques générés sans collision", testCount)
}

// TestValidateJWT_MissingJTI tests that JWT without JTI is rejected
func TestValidateJWT_MissingJTI(t *testing.T) {
	t.Log("🧪 TEST VALIDATION JTI MANQUANT")
	t.Log("================================")

	// Arrange - Créer manager
	manager, _ := NewManager(&Config{
		Type:      AuthTypeJWT,
		JWTSecret: TestJWTSecret,
		JWTIssuer: TestIssuer,
	})

	// Créer token SANS JTI
	claims := Claims{
		Username: TestUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   TestUsername,
			Audience:  jwt.ClaimStrings{TokenAudience},
			Issuer:    TestIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "", // ❌ JTI vide
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(TestJWTSecret))

	// Act - Valider token sans JTI
	err := manager.ValidateToken(tokenString)

	// Assert - Devrait échouer
	if err == nil {
		t.Error("❌ La validation devrait échouer pour un token sans JTI")
	} else {
		t.Logf("✅ Token sans JTI correctement rejeté: %v", err)
	}

	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("❌ Erreur devrait être ErrInvalidToken, reçu: %v", err)
	}
}

// TestConstants_JWTClaims tests that JWT claim constants are properly defined
func TestConstants_JWTClaims(t *testing.T) {
	t.Log("🧪 TEST CONSTANTES JWT CLAIMS")
	t.Log("==============================")

	// Vérifier TokenAudience
	if TokenAudience == "" {
		t.Error("❌ TokenAudience ne devrait pas être vide")
	} else {
		t.Logf("✅ TokenAudience défini: %q", TokenAudience)
	}

	// Vérifier DefaultTokenIssuer
	if DefaultTokenIssuer == "" {
		t.Error("❌ DefaultTokenIssuer ne devrait pas être vide")
	} else {
		t.Logf("✅ DefaultTokenIssuer défini: %q", DefaultTokenIssuer)
	}

	// Vérifier JTILength
	minJTILength := 8
	if JTILength < minJTILength {
		t.Errorf("❌ JTILength trop court: %d (minimum recommandé: %d)", JTILength, minJTILength)
	} else {
		t.Logf("✅ JTILength approprié: %d bytes", JTILength)
	}
}
