package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/lutefd/ai-router-go/internal/repository"
)

type AIService struct {
	geminiRepo   repository.AIRepositoryInterface
	openaiRepo   repository.AIRepositoryInterface
	deepseekRepo repository.AIRepositoryInterface
}

func NewAIService(geminiRepo repository.AIRepositoryInterface,
	openaiRepo repository.AIRepositoryInterface,
	deepseekRepo repository.AIRepositoryInterface) *AIService {
	return &AIService{
		geminiRepo:   geminiRepo,
		openaiRepo:   openaiRepo,
		deepseekRepo: deepseekRepo,
	}
}

func (s *AIService) GenerateResponse(ctx context.Context, model string,
	prompt string, callback func(string)) error {
	if strings.TrimSpace(prompt) == "" {
		return fmt.Errorf("empty prompt")
	}

	if s.geminiRepo == nil {
		return fmt.Errorf("gemini repository not initialized")
	}
	return s.geminiRepo.GenerateContentStream(ctx, model, prompt, callback)
}

func (s *AIService) GenerateOpenAIResponse(ctx context.Context, model string,
	prompt string, callback func(string)) error {
	if s.openaiRepo == nil {
		return fmt.Errorf("openai repository not initialized")
	}
	return s.openaiRepo.GenerateContentStream(ctx, model, prompt, callback)
}

func (s *AIService) GenerateDeepSeekResponse(ctx context.Context, model string,
	prompt string, callback func(string)) error {
	if s.deepseekRepo == nil {
		return fmt.Errorf("deepseek repository not initialized")
	}
	return s.deepseekRepo.GenerateContentStream(ctx, model, prompt, callback)
}
