package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"SERVER_PORT":          "8080",
				"OPENAI_SK":            "sk-123",
				"DEEPSEEK_SK":          "sk-456",
				"GEMINI_SK":            "sk-789",
				"MONGODB_URI":          "mongodb://localhost:27017",
				"MONGODB_DATABASE":     "ai_router",
				"GOOGLE_CLIENT_ID":     "client-123",
				"GOOGLE_CLIENT_SECRET": "secret-456",
				"JWT_SECRET":           "jwt-secret-789",
				"CLIENT_URL":           "http://localhost:3000",
				"AUTH_REDIRECT_URL":    "http://localhost:8080/callback",
			},
			expectError: false,
		},
		{
			name: "missing SERVER_PORT",
			envVars: map[string]string{
				"OPENAI_SK":            "sk-123",
				"DEEPSEEK_SK":          "sk-456",
				"GEMINI_SK":            "sk-789",
				"MONGODB_URI":          "mongodb://localhost:27017",
				"MONGODB_DATABASE":     "ai_router",
				"GOOGLE_CLIENT_ID":     "client-123",
				"GOOGLE_CLIENT_SECRET": "secret-456",
				"JWT_SECRET":           "jwt-secret-789",
				"CLIENT_URL":           "http://localhost:3000",
				"AUTH_REDIRECT_URL":    "http://localhost:8080/callback",
			},
			expectError: true,
		},
		{
			name: "missing MongoDB configuration",
			envVars: map[string]string{
				"SERVER_PORT":          "8080",
				"OPENAI_SK":            "sk-123",
				"DEEPSEEK_SK":          "sk-456",
				"GEMINI_SK":            "sk-789",
				"GOOGLE_CLIENT_ID":     "client-123",
				"GOOGLE_CLIENT_SECRET": "secret-456",
				"JWT_SECRET":           "jwt-secret-789",
				"CLIENT_URL":           "http://localhost:3000",
				"AUTH_REDIRECT_URL":    "http://localhost:8080/callback",
			},
			expectError: true,
		},
		{
			name: "missing Google OAuth configuration",
			envVars: map[string]string{
				"SERVER_PORT":       "8080",
				"OPENAI_SK":         "sk-123",
				"DEEPSEEK_SK":       "sk-456",
				"GEMINI_SK":         "sk-789",
				"MONGODB_URI":       "mongodb://localhost:27017",
				"MONGODB_DATABASE":  "ai_router",
				"JWT_SECRET":        "jwt-secret-789",
				"CLIENT_URL":        "http://localhost:3000",
				"AUTH_REDIRECT_URL": "http://localhost:8080/callback",
			},
			expectError: true,
		},
		{
			name: "missing JWT configuration",
			envVars: map[string]string{
				"SERVER_PORT":          "8080",
				"OPENAI_SK":            "sk-123",
				"DEEPSEEK_SK":          "sk-456",
				"GEMINI_SK":            "sk-789",
				"MONGODB_URI":          "mongodb://localhost:27017",
				"MONGODB_DATABASE":     "ai_router",
				"GOOGLE_CLIENT_ID":     "client-123",
				"GOOGLE_CLIENT_SECRET": "secret-456",
				"CLIENT_URL":           "http://localhost:3000",
				"AUTH_REDIRECT_URL":    "http://localhost:8080/callback",
			},
			expectError: true,
		},
		{
			name: "missing URL configuration",
			envVars: map[string]string{
				"SERVER_PORT":          "8080",
				"OPENAI_SK":            "sk-123",
				"DEEPSEEK_SK":          "sk-456",
				"GEMINI_SK":            "sk-789",
				"MONGODB_URI":          "mongodb://localhost:27017",
				"MONGODB_DATABASE":     "ai_router",
				"GOOGLE_CLIENT_ID":     "client-123",
				"GOOGLE_CLIENT_SECRET": "secret-456",
				"JWT_SECRET":           "jwt-secret-789",
			},
			expectError: true,
		},
		{
			name: "invalid SERVER_PORT",
			envVars: map[string]string{
				"SERVER_PORT":          "not-a-number",
				"OPENAI_SK":            "sk-123",
				"DEEPSEEK_SK":          "sk-456",
				"GEMINI_SK":            "sk-789",
				"MONGODB_URI":          "mongodb://localhost:27017",
				"MONGODB_DATABASE":     "ai_router",
				"GOOGLE_CLIENT_ID":     "client-123",
				"GOOGLE_CLIENT_SECRET": "secret-456",
				"JWT_SECRET":           "jwt-secret-789",
				"CLIENT_URL":           "http://localhost:3000",
				"AUTH_REDIRECT_URL":    "http://localhost:8080/callback",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg, err := LoadConfig(true)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.envVars["SERVER_PORT"], strconv.Itoa(cfg.ServerPort))
			assert.Equal(t, tt.envVars["OPENAI_SK"], cfg.OPENAI_SK)
			assert.Equal(t, tt.envVars["DEEPSEEK_SK"], cfg.DEEPSEEK_SK)
			assert.Equal(t, tt.envVars["GEMINI_SK"], cfg.GEMINI_SK)
			assert.Equal(t, tt.envVars["MONGODB_URI"], cfg.MongoDBURI)
			assert.Equal(t, tt.envVars["MONGODB_DATABASE"], cfg.MongoDBDatabase)
			assert.Equal(t, tt.envVars["GOOGLE_CLIENT_ID"], cfg.GoogleClientID)
			assert.Equal(t, tt.envVars["GOOGLE_CLIENT_SECRET"], cfg.GoogleClientSecret)
			assert.Equal(t, tt.envVars["JWT_SECRET"], cfg.JWTSecret)
			assert.Equal(t, tt.envVars["CLIENT_URL"], cfg.ClientURL)
			assert.Equal(t, tt.envVars["AUTH_REDIRECT_URL"], cfg.AuthRedirectURL)
		})
	}
}
