package strategy

import "context"

type AIStrategyInterface interface {
	GenerateResponse(ctx context.Context, platform string, model string,
		prompt string, callback func(string)) error
}
