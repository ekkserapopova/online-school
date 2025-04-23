package provider

import (
	"context"
)

// LLMProvider определяет интерфейс для взаимодействия с языковыми моделями
type LLMProvider interface {
	// GenerateCompletion генерирует ответ на основе системного и пользовательского промптов
	GenerateCompletion(ctx context.Context, systemPrompt, userPrompt string, schema interface{}) (string, error)
}
