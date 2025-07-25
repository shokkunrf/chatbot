package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"chatbot/store"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiClient はGoogle Gemini APIを使用してAI Clientインターフェースを実装する
type GeminiClient struct {
	client       *genai.Client
	model        *genai.GenerativeModel
	systemPrompt string
	timeout      time.Duration
}

// NewGeminiClient は新しいGeminiクライアントインスタンスを作成する
func NewGeminiClient(apiKey string, systemPrompt string) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, errors.New("APIキーが必要です")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Geminiクライアントの作成に失敗しました: %w", err)
	}

	// Gemini 2.5 Flashモデルを使用
	model := client.GenerativeModel("gemini-2.5-flash")

	// モデルパラメータを設定
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)

	return &GeminiClient{
		client:       client,
		model:        model,
		systemPrompt: systemPrompt,
		timeout:      60 * time.Second,
	}, nil
}

// GenerateResponse は会話履歴に基づいて応答を生成する
func (c *GeminiClient) GenerateResponse(messages []store.Message, customSystemPrompt string) (string, error) {
	if len(messages) == 0 {
		return "", errors.New("メッセージが提供されていません")
	}

	// カスタムシステムプロンプトが提供された場合はそれを使用、そうでなければデフォルトを使用
	systemPrompt := c.systemPrompt
	if customSystemPrompt != "" {
		systemPrompt = customSystemPrompt
	}

	// タイムアウト付きコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// チャットセッションを作成
	cs := c.model.StartChat()

	// システムプロンプトを設定（175文字制限を含む）
	fullPrompt := strings.ToValidUTF8(systemPrompt+"\n\n重要: 応答は必ず175文字以内で簡潔にまとめてください。", "")
	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text(fullPrompt),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("了解しました。その設定で会話を進めます。"),
			},
			Role: "model",
		},
	}

	// 会話履歴を追加（最後のメッセージを除く）
	for i := 0; i < len(messages)-1; i++ {
		msg := messages[i]
		role := "user"
		if msg.Role != "user" {
			role = "model"
		}

		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(strings.ToValidUTF8(msg.Content, "")),
			},
			Role: role,
		})
	}

	// 最後のメッセージを送信
	lastMessage := strings.ToValidUTF8(messages[len(messages)-1].Content, "")
	resp, err := cs.SendMessage(ctx, genai.Text(lastMessage))
	if err != nil {
		return "", fmt.Errorf("応答の生成に失敗しました: %w", err)
	}

	// 応答を抽出
	return extractResponse(resp)
}

// extractResponse はGeminiのレスポンスからテキストを抽出する
func extractResponse(resp *genai.GenerateContentResponse) (string, error) {
	if resp == nil {
		return "", fmt.Errorf("レスポンスがnilです")
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("候補がありません")
	}

	candidate := resp.Candidates[0]

	// コンテンツが存在しない場合のエラーハンドリング
	if candidate.Content == nil {
		switch candidate.FinishReason {
		case genai.FinishReasonSafety:
			return "申し訳ありません。安全性の理由により応答できません。", nil
		case genai.FinishReasonRecitation:
			return "申し訳ありません。著作権の問題により応答できません。", nil
		case genai.FinishReasonMaxTokens:
			return "申し訳ありません。応答が長すぎます。", nil
		default:
			return "", fmt.Errorf("コンテンツが空です: FinishReason=%v", candidate.FinishReason)
		}
	}

	// Partsからテキストを抽出
	var responseText string
	for _, part := range candidate.Content.Parts {
		if text, ok := part.(genai.Text); ok {
			responseText += string(text)
		}
	}

	if responseText == "" {
		return "申し訳ありません。応答を生成できませんでした。", nil
	}

	return responseText, nil
}

// Close はGeminiクライアント接続を閉じる
func (c *GeminiClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
