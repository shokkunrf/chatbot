package ai

import "chatbot/store"

// Client インターフェースはAIとの対話メソッドを定義する
type Client interface {
	GenerateResponse(messages []store.Message, systemPrompt string) (string, error)
}
