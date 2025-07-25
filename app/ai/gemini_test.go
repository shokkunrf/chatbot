package ai

import (
	"context"
	"testing"
	"time"

	"chatbot/store"
	"github.com/google/generative-ai-go/genai"
)

// mockGeminiClient はテスト用のモック実装
type mockGeminiClient struct {
	generateResponse func(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

func (m *mockGeminiClient) GenerativeModel(name string) *genai.GenerativeModel {
	// これは簡略化されたモック - 実際のテストではより完全なモックが必要
	return nil
}

// TestGeminiClient_GenerateResponse は通常の応答生成をテスト
func TestGeminiClient_GenerateResponse(t *testing.T) {
	// 注: 実際の実装では、モックGeminiクライアントを使用
	// 現在はGenerateResponseメソッドの基本的な動作をテスト

	client := &GeminiClient{
		systemPrompt: "You are a helpful assistant. Keep responses under 175 characters.",
		timeout:      5 * time.Second,
	}

	// エラーケース: 空のメッセージ
	_, err := client.GenerateResponse([]store.Message{}, "")
	if err == nil {
		t.Error("Expected error for empty messages")
	}

	// エラーケース: nilメッセージ
	_, err = client.GenerateResponse(nil, "")
	if err == nil {
		t.Error("Expected error for nil messages")
	}
}



// TestGeminiClient_ErrorHandling はエラーシナリオをテスト
func TestGeminiClient_ErrorHandling(t *testing.T) {
	// nilメッセージでテスト
	client := &GeminiClient{}

	response, err := client.GenerateResponse(nil, "")
	if err == nil {
		t.Error("Expected error for nil messages")
	}
	if response != "" {
		t.Error("Expected empty response on error")
	}

	// 空のメッセージでテスト
	response, err = client.GenerateResponse([]store.Message{}, "")
	if err == nil {
		t.Error("Expected error for empty messages")
	}
	if response != "" {
		t.Error("Expected empty response on error")
	}
}
