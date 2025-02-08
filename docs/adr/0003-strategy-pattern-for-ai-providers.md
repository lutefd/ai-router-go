# 3. Strategy Pattern for AI Providers

## Status

Accepted

## Context

We needed to support multiple AI providers (OpenAI, Google Gemini, DeepSeek) while:

- Maintaining clean code organization
- Making it easy to add new providers
- Supporting streaming responses
- Handling provider-specific configurations
- Managing rate limits and quotas
- Providing fallback mechanisms

Key requirements:

- Runtime provider selection
- Consistent error handling
- Provider-specific optimizations
- Easy testing and mocking
- Cost optimization capabilities

## Decision

We implemented the Strategy pattern for AI provider selection and interaction with the following structure:

1. **Interface Definition**

   ```go
   type AIStrategyInterface interface {
       GenerateResponse(ctx context.Context, platform string, model string,
           prompt string, callback func(string)) error
   }
   ```

2. **Provider Implementations**

   ```go
   type GeminiRepository struct {
       client *genai.Client
   }

   type OpenAIRepository struct {
       client *openai.Client
   }

   type DeepSeekRepository struct {
       client *openai.Client
   }
   ```

3. **Strategy Selection**

   ```go
   func (s *AIStrategy) GenerateResponse(ctx context.Context, platform string,
       model string, prompt string, callback func(string)) error {
       switch platform {
       case "gemini":
           return s.aiService.GenerateResponse(ctx, model, prompt, callback)
       case "openai":
           return s.aiService.GenerateOpenAIResponse(ctx, model, prompt, callback)
       case "deepseek":
           return s.aiService.GenerateDeepSeekResponse(ctx, model, prompt, callback)
       default:
           return fmt.Errorf("unsupported platform: %s", platform)
       }
   }
   ```

4. **Error Handling**
   ```go
   func (r *Repository) GenerateContentStream(ctx context.Context,
       modelName string, prompt string, callback func(string)) error {
       if strings.TrimSpace(prompt) == "" {
           return fmt.Errorf("empty prompt")
       }
   }
   ```

## Consequences

### Positive

- Clean separation of provider-specific logic
- Easy to add new AI providers
- Runtime provider selection
- Consistent interface for all providers
- Simplified testing through interface mocking
- Independent scaling of providers
- Isolated provider configurations
- Centralized error handling
- Easy to implement fallback strategies

### Negative

- Additional abstraction layer
- Slight performance overhead
- Need to maintain separate implementations for each provider
- Complexity in handling provider-specific features
- Need to normalize responses across providers

### Mitigations

1. For provider-specific features:

   ```go
   type ProviderConfig struct {
       MaxTokens    int
       Temperature  float32
       TopP        float32
       ContextSize int
   }

   func NewProvider(config ProviderConfig) AIProvider {
       // Initialize provider with specific configuration
   }
   ```

2. For response normalization:

   ```go
   func normalizeResponse(providerResponse interface{}) (string, error) {
       switch resp := providerResponse.(type) {
       case *openai.ChatCompletionStreamResponse:
           return resp.Choices[0].Delta.Content, nil
       case *genai.GenerateContentResponse:
           return resp.Candidates[0].Content.Parts[0].Text, nil
       default:
           return "", fmt.Errorf("unsupported response type")
       }
   }
   ```

3. For error handling:
   ```go
   func handleProviderError(err error, provider string) error {
       switch {
       case errors.Is(err, context.DeadlineExceeded):
           return fmt.Errorf("%s provider timeout: %w", provider, err)
       case errors.Is(err, ErrRateLimit):
           return fmt.Errorf("%s rate limit exceeded: %w", provider, err)
       default:
           return fmt.Errorf("%s provider error: %w", provider, err)
       }
   }
   ```
