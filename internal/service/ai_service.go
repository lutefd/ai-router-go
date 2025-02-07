package service

import (
	"context"

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
	return s.geminiRepo.GenerateContentStream(ctx, model, prompt, callback)
}

func (s *AIService) GenerateOpenAIResponse(ctx context.Context, model string,
	prompt string, callback func(string)) error {
	return s.openaiRepo.GenerateContentStream(ctx, model, prompt, callback)
}

func (s *AIService) GenerateDeepSeekResponse(ctx context.Context, model string,
	prompt string, callback func(string)) error {
	return s.deepseekRepo.GenerateContentStream(ctx, model, prompt, callback)
}
