package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/lutefd/ai-router-go/internal/mocks"
	"github.com/lutefd/ai-router-go/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAIService_GenerateResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	geminiMock := mocks.NewMockAIRepositoryInterface(ctrl)
	openaiMock := mocks.NewMockAIRepositoryInterface(ctrl)
	deepseekMock := mocks.NewMockAIRepositoryInterface(ctrl)

	aiService := service.NewAIService(geminiMock, openaiMock, deepseekMock)

	tests := []struct {
		name      string
		prompt    string
		setupMock func()
		wantErr   bool
	}{
		{
			name:   "successful generation",
			prompt: "test prompt",
			setupMock: func() {
				geminiMock.EXPECT().GenerateContentStream(gomock.Any(), "gemini-pro", "test prompt", gomock.Any()).DoAndReturn(func(ctx context.Context, model string, prompt string, callback func(string)) error {
					callback("test response")
					return nil
				})
			},
			wantErr: false,
		},
		{
			name:      "empty prompt",
			prompt:    "",
			setupMock: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var responses []string
			callback := func(response string) {
				responses = append(responses, response)
			}

			tt.setupMock()
			err := aiService.GenerateResponse(context.Background(), "gemini-pro", tt.prompt, callback)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, responses)
		})
	}
}

func TestAIService_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	geminiMock := mocks.NewMockAIRepositoryInterface(ctrl)
	openaiMock := mocks.NewMockAIRepositoryInterface(ctrl)
	deepseekMock := mocks.NewMockAIRepositoryInterface(ctrl)

	aiService := service.NewAIService(geminiMock, openaiMock, deepseekMock)

	geminiMock.EXPECT().GenerateContentStream(gomock.Any(), "gemini-pro", "test prompt", gomock.Any()).DoAndReturn(func(ctx context.Context, model string, prompt string, callback func(string)) error {
		select {
		case <-time.After(10 * time.Millisecond):
			callback("too late")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	err := aiService.GenerateResponse(ctx, "gemini-pro", "test prompt", func(string) {})
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
