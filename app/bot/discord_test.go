package bot

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"chatbot/store"
)

// MockStore はテスト用のモックストア
type MockStore struct {
	messages  map[string][]store.Message
	addCalled int
}

func NewMockStore() *MockStore {
	return &MockStore{
		messages: make(map[string][]store.Message),
	}
}

func (m *MockStore) GetHistory(guildID, channelID string) []store.Message {
	key := guildID + ":" + channelID
	return m.messages[key]
}

func (m *MockStore) AddMessage(guildID, channelID string, message store.Message) error {
	m.addCalled++
	key := guildID + ":" + channelID
	m.messages[key] = append(m.messages[key], message)
	return nil
}

func (m *MockStore) ClearHistory(guildID, channelID string) error {
	key := guildID + ":" + channelID
	delete(m.messages, key)
	return nil
}

// MockAIClient はテスト用のモックAIクライアント
type MockAIClient struct {
	response string
	err      error
}

func (m *MockAIClient) GenerateResponse(messages []store.Message, systemPrompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

// TestShouldRespond はボットが応答すべきかどうかの判定をテスト
func TestShouldRespond(t *testing.T) {
	botID := "bot123"

	tests := []struct {
		name          string
		message       *discordgo.Message
		shouldRespond bool
	}{
		{
			name: "ボットへのメンション",
			message: &discordgo.Message{
				Content:  "<@bot123> こんにちは",
				Author:   &discordgo.User{ID: "user1", Bot: false},
				GuildID:  "test-guild",
				Mentions: []*discordgo.User{{ID: botID}},
			},
			shouldRespond: true,
		},
		{
			name: "ボットへの返信",
			message: &discordgo.Message{
				Content: "ありがとう！",
				Author:  &discordgo.User{ID: "user1", Bot: false},
				GuildID: "test-guild",
				ReferencedMessage: &discordgo.Message{
					Author: &discordgo.User{ID: botID, Bot: true},
				},
			},
			shouldRespond: true,
		},
		{
			name: "通常のメッセージ",
			message: &discordgo.Message{
				Content:  "ただの会話",
				Author:   &discordgo.User{ID: "user1", Bot: false},
				GuildID:  "test-guild",
				Mentions: []*discordgo.User{},
			},
			shouldRespond: false,
		},
		{
			name: "ボットからのメッセージ",
			message: &discordgo.Message{
				Content: "私はボットです",
				Author:  &discordgo.User{ID: "bot456", Bot: true},
				GuildID: "test-guild",
			},
			shouldRespond: false,
		},
		{
			name: "DMメッセージ",
			message: &discordgo.Message{
				Content: "DMでこんにちは",
				Author:  &discordgo.User{ID: "user1", Bot: false},
				GuildID: "",
			},
			shouldRespond: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{botID: botID}
			result := b.shouldRespond(tt.message)
			if result != tt.shouldRespond {
				t.Errorf("期待値 shouldRespond=%v, 実際の値 %v", tt.shouldRespond, result)
			}
		})
	}
}

// TestCleanContent はボットメンションの除去をテスト
func TestCleanContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		botID    string
		expected string
	}{
		{
			name:     "先頭のボットメンションを削除",
			content:  "<@bot123> こんにちは世界",
			botID:    "bot123",
			expected: "こんにちは世界",
		},
		{
			name:     "途中のボットメンションを削除",
			content:  "ねえ <@bot123> 手伝って？",
			botID:    "bot123",
			expected: "ねえ 手伝って？",
		},
		{
			name:     "メンションなし",
			content:  "ただの通常のメッセージ",
			botID:    "bot123",
			expected: "ただの通常のメッセージ",
		},
		{
			name:     "メンション後の複数スペース",
			content:  "<@bot123>   こんにちは",
			botID:    "bot123",
			expected: "こんにちは",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{botID: tt.botID}
			result := b.cleanContent(tt.content)
			if result != tt.expected {
				t.Errorf("期待値 '%s', 実際の値 '%s'", tt.expected, result)
			}
		})
	}
}

// TestProcessMessage はメッセージ処理の統合テスト
func TestProcessMessage(t *testing.T) {
	mockStore := NewMockStore()
	mockAI := &MockAIClient{response: "テスト応答です"}

	bot := &Bot{
		botID:        "bot123",
		store:        mockStore,
		aiClient:     mockAI,
		systemPrompt: "テスト用システムプロンプト",
	}

	// テストメッセージ
	message := &discordgo.Message{
		ID:        "msg123",
		ChannelID: "channel123",
		GuildID:   "guild123",
		Content:   "<@bot123> こんにちは",
		Author: &discordgo.User{
			ID:       "user123",
			Username: "TestUser",
			Bot:      false,
		},
		Mentions:  []*discordgo.User{{ID: "bot123"}},
		Timestamp: time.Now(),
	}

	// メッセージを処理
	response, err := bot.processMessage(message)
	if err != nil {
		t.Fatalf("メッセージ処理でエラー: %v", err)
	}

	if response != "テスト応答です" {
		t.Errorf("期待値 'テスト応答です', 実際の値 '%s'", response)
	}

	// ストアに2つのメッセージが追加されたことを確認（ユーザーとアシスタント）
	if mockStore.addCalled != 2 {
		t.Errorf("期待値 2回のAddMessage呼び出し, 実際の値 %d", mockStore.addCalled)
	}

	// 履歴を確認
	history := mockStore.GetHistory("guild123", "channel123")
	if len(history) != 2 {
		t.Errorf("期待値 2メッセージ, 実際の値 %d", len(history))
	}
}

// TestHandleResetCommand は/resetコマンドのテスト
func TestHandleResetCommand(t *testing.T) {
	mockStore := NewMockStore()

	// 事前にメッセージを追加
	mockStore.AddMessage("guild123", "channel123", store.Message{
		Role:    "user",
		Content: "テストメッセージ",
	})

	bot := &Bot{
		store: mockStore,
	}

	// リセット前の確認
	history := mockStore.GetHistory("guild123", "channel123")
	if len(history) != 1 {
		t.Errorf("リセット前: 期待値 1メッセージ, 実際の値 %d", len(history))
	}

	// リセットコマンドを実行
	err := bot.handleResetCommand("guild123", "channel123")
	if err != nil {
		t.Fatalf("リセットコマンドでエラー: %v", err)
	}

	// リセット後の確認
	history = mockStore.GetHistory("guild123", "channel123")
	if len(history) != 0 {
		t.Errorf("リセット後: 期待値 0メッセージ, 実際の値 %d", len(history))
	}
}
