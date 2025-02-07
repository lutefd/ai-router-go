package strategy

import (
	"context"
	"fmt"

	"github.com/lutefd/ai-router-go/internal/service"
)

type AIStrategy struct {
	aiService service.AIServiceInterface
}

func NewAIStrategy(aiService service.AIServiceInterface) *AIStrategy {
	return &AIStrategy{aiService: aiService}
}

func (s *AIStrategy) GenerateResponse(ctx context.Context, platform string,
	model string, prompt string, callback func(string)) error {
	switch platform {
	case "gemini":
		return s.aiService.GenerateResponse(ctx, model, prompt, callback)
	case "openai":
		return s.aiService.GenerateOpenAIResponse(ctx, model, prompt, callback)
	case "deepseek":
		return s.aiService.GenerateDeepSeekResponse(ctx, model, prompt,
			callback)
	default:
		return fmt.Errorf("unsupported platform: %s", platform)
	}
}
