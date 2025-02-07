package repository

import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type DeepSeekRepository struct {
	client *openai.Client
}

func NewDeepSeekRepository(apiKey string) *DeepSeekRepository {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.deepseek.com/v1"
	client := openai.NewClientWithConfig(config)

	return &DeepSeekRepository{client: client}
}

func (r *DeepSeekRepository) GenerateContentStream(ctx context.Context,
	modelName string, prompt string, callback func(string)) error {
	req := openai.ChatCompletionRequest{
		Model: modelName,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}

	streamer, err := r.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("error creating stream: %w", err)
	}
	defer streamer.Close()

	for {
		response, err := streamer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving stream data: %w", err)
		}

		callback(response.Choices[0].Delta.Content)
	}

	return nil
}
