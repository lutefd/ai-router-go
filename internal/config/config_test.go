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
				"SERVER_PORT": "8080",
				"OPENAI_SK":   "sk-123",
				"DEEPSEEK_SK": "sk-456",
				"GEMINI_SK":   "sk-789",
			},
			expectError: false,
		},
		{
			name: "missing SERVER_PORT",
			envVars: map[string]string{
				"OPENAI_SK":   "sk-123",
				"DEEPSEEK_SK": "sk-456",
				"GEMINI_SK":   "sk-789",
			},
			expectError: true,
		},
		{
			name: "invalid SERVER_PORT",
			envVars: map[string]string{
				"SERVER_PORT": "not-a-number",
				"OPENAI_SK":   "sk-123",
				"DEEPSEEK_SK": "sk-456",
				"GEMINI_SK":   "sk-789",
			},
			expectError: true,
		},
		{
			name: "missing OPENAI_SK",
			envVars: map[string]string{
				"SERVER_PORT": "8080",
				"DEEPSEEK_SK": "sk-456",
				"GEMINI_SK":   "sk-789",
			},
			expectError: true,
		},
		{
			name: "missing DEEPSEEK_SK",
			envVars: map[string]string{
				"SERVER_PORT": "8080",
				"OPENAI_SK":   "sk-123",
				"GEMINI_SK":   "sk-789",
			},
			expectError: true,
		},
		{
			name: "missing GEMINI_SK",
			envVars: map[string]string{
				"SERVER_PORT": "8080",
				"OPENAI_SK":   "sk-123",
				"DEEPSEEK_SK": "sk-456",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			cfg, err := LoadConfig(true)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.envVars["SERVER_PORT"], strconv.Itoa(cfg.ServerPort))
			assert.Equal(t, tt.envVars["OPENAI_sK"], cfg.OPENAI_SK)
			assert.Equal(t, tt.envVars["DEEPSEEK_SK"], cfg.DEEPSEEK_SK)
			assert.Equal(t, tt.envVars["GEMINI_SK"], cfg.GEMINI_SK)

		})
	}
}
