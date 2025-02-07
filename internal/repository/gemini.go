package repository

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiRepository struct {
	client *genai.Client
}

func NewGeminiRepository(ctx context.Context, geminiSK string) *GeminiRepository {
	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiSK,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create gemini client: %w", err))
	}
	return &GeminiRepository{client: geminiClient}
}

func (r *GeminiRepository) GenerateContentStream(ctx context.Context,
	modelName string, prompt string, callback func(string)) error {
	model := r.client.Models
	for result, err := range model.GenerateContentStream(
		ctx,
		modelName,
		genai.Text(prompt),
		nil,
	) {
		if err != nil {
			return err
		}
		callback(result.Candidates[0].Content.Parts[0].Text)
	}

	return nil
}
