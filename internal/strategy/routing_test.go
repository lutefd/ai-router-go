package strategy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lutefd/ai-router-go/internal/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockAIService struct {
	generateFunc         func(ctx context.Context, model string, prompt string, callback func(string)) error
	generateOpenAIFunc   func(ctx context.Context, model string, prompt string, callback func(string)) error
	generateDeepSeekFunc func(ctx context.Context, model string, prompt string, callback func(string)) error
}

func (m *MockAIService) GenerateResponse(ctx context.Context, model string, prompt string, callback func(string)) error {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, model, prompt, callback)
	}
	return nil
}

func (m *MockAIService) GenerateOpenAIResponse(ctx context.Context, model string, prompt string, callback func(string)) error {
	if m.generateOpenAIFunc != nil {
		return m.generateOpenAIFunc(ctx, model, prompt, callback)
	}
	return nil
}

func (m *MockAIService) GenerateDeepSeekResponse(ctx context.Context, model string, prompt string, callback func(string)) error {
	if m.generateDeepSeekFunc != nil {
		return m.generateDeepSeekFunc(ctx, model, prompt, callback)
	}
	return nil
}

func TestAIStrategy_GenerateResponse(t *testing.T) {
	mockService := &MockAIService{}
	aiStrategy := strategy.NewAIStrategy(mockService)

	tests := []struct {
		name      string
		platform  string
		model     string
		prompt    string
		setupMock func()
		wantResp  string
		wantErr   bool
	}{
		{
			name:     "successful Gemini generation",
			platform: "gemini",
			model:    "gemini-pro",
			prompt:   "test prompt",
			setupMock: func() {
				mockService.generateFunc = func(ctx context.Context, model string, prompt string, callback func(string)) error {
					callback("gemini response")
					return nil
				}
			},
			wantResp: "gemini response",
			wantErr:  false,
		},
		{
			name:     "successful OpenAI generation",
			platform: "openai",
			model:    "gpt-3.5-turbo",
			prompt:   "test prompt",
			setupMock: func() {
				mockService.generateOpenAIFunc = func(ctx context.Context, model string, prompt string, callback func(string)) error {
					callback("openai response")
					return nil
				}
			},
			wantResp: "openai response",
			wantErr:  false,
		},
		{
			name:     "successful DeepSeek generation",
			platform: "deepseek",
			model:    "deepseek-chat",
			prompt:   "test prompt",
			setupMock: func() {
				mockService.generateDeepSeekFunc = func(ctx context.Context, model string, prompt string, callback func(string)) error {
					callback("deepseek response")
					return nil
				}
			},
			wantResp: "deepseek response",
			wantErr:  false,
		},
		{
			name:      "unsupported platform",
			platform:  "unsupported",
			model:     "some-model",
			prompt:    "test prompt",
			setupMock: func() {},
			wantErr:   true,
		},
		{
			name:     "service error",
			platform: "gemini",
			model:    "gemini-pro",
			prompt:   "test prompt",
			setupMock: func() {
				mockService.generateFunc = func(ctx context.Context, model string, prompt string, callback func(string)) error {
					return fmt.Errorf("service error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var response string
			callback := func(resp string) {
				response = resp
			}

			err := aiStrategy.GenerateResponse(context.Background(), tt.platform, tt.model, tt.prompt, callback)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantResp != "" {
				assert.Equal(t, tt.wantResp, response)
			}
		})
	}
}
