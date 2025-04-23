package service

import (
	"context"
	"encoding/json"
	"errors"

	"onlineschool/internal/llm/provider"
	"onlineschool/internal/llm/schema"
)

type CodeEvaluator struct {
	provider      provider.LLMProvider
	promptBuilder *PromptBuilder
}

func NewCodeEvaluator(provider provider.LLMProvider) *CodeEvaluator {
	return &CodeEvaluator{
		provider:      provider,
		promptBuilder: NewPromptBuilder(),
	}
}

func (e *CodeEvaluator) EvaluateCode(ctx context.Context, code, taskDescription string) (*schema.CodeEvaluation, error) {

	codeEvaluationSchema := schema.GenerateSchema[schema.CodeEvaluation]()

	systemPrompt := e.promptBuilder.BuildSystemPrompt()
	userPrompt := e.promptBuilder.BuildEvaluationPrompt(taskDescription, code)

	// Отправляем запрос к провайдеру
	content, err := e.provider.GenerateCompletion(ctx, systemPrompt, userPrompt, codeEvaluationSchema)
	if err != nil {
		return nil, err
	}

	// Разбираем JSON-ответ
	var evalResponse schema.CodeEvaluation
	err = json.Unmarshal([]byte(content), &evalResponse)
	if err != nil {
		return nil, errors.New("ошибка при разборе JSON-ответа: " + err.Error())
	}

	return &evalResponse, nil
}
