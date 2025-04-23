package provider

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIProvider реализует интерфейс LLMProvider для работы с OpenAI API
type OpenAIProvider struct {
	client openai.Client
	model  string
}

// NewOpenAIProvider создает новый провайдер для OpenAI
func NewOpenAIProvider(apiKey string, model string) *OpenAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey))

	if model == "" {
		model = openai.ChatModelGPT4o
	}

	return &OpenAIProvider{
		client: client,
		model:  model,
	}
}

// GenerateCompletion реализует метод интерфейса LLMProvider
func (p *OpenAIProvider) GenerateCompletion(ctx context.Context, systemPrompt, userPrompt string, schema interface{}) (string, error) {
	// Создаем параметры для JSON-схемы
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "code_evaluation",
		Description: openai.String("Оценка кода с детальным анализом по различным аспектам"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	// Создаем запрос к API OpenAI
	chatParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
		Model: p.model,
	}

	// Отправляем запрос
	chat, err := p.client.Chat.Completions.New(ctx, chatParams)
	if err != nil {
		return "", err
	}

	// Возвращаем контент первого сообщения
	return chat.Choices[0].Message.Content, nil
}
