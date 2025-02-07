package repository

import (
	"context"
)

type AIRepositoryInterface interface {
	GenerateContentStream(ctx context.Context, model string, prompt string,
		callback func(string)) error
}
