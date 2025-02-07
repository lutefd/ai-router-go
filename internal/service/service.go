package service

import "context"

type AIServiceInterface interface {
	GenerateResponse(ctx context.Context, model string, prompt string,
		callback func(string)) error
	GenerateOpenAIResponse(ctx context.Context, model string, prompt string,
		callback func(string)) error
	GenerateDeepSeekResponse(ctx context.Context, model string, prompt string,
		callback func(string)) error
}
