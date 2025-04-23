package llm

import "os"

type Config struct {
	APIKey string
	Model  string
	//MaxTokens int
	//Temperature float32
}

// NewDefaultConfig создает конфигурацию по умолчанию
func NewDefaultConfig() *Config {
	return &Config{
		APIKey: os.Getenv("API_KEY"),
		Model:  "gpt-4o",
		//MaxTokens:   2048,
		//Temperature: 0.0, // Для оценки кода лучше использовать низкую температуру
	}
}
