package config

import (
	"os"
	"testing"
	"time"
)

// TestLoadConfig は設定の読み込みをテスト
func TestLoadConfig(t *testing.T) {
	// 環境変数を保存してテスト後に復元
	oldToken := os.Getenv("DISCORD_BOT_TOKEN")
	oldAPIKey := os.Getenv("GEMINI_API_KEY")
	defer func() {
		os.Setenv("DISCORD_BOT_TOKEN", oldToken)
		os.Setenv("GEMINI_API_KEY", oldAPIKey)
	}()

	// 必須項目を設定
	os.Setenv("DISCORD_BOT_TOKEN", "test-token")
	os.Setenv("GEMINI_API_KEY", "test-api-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("設定の読み込みに失敗: %v", err)
	}

	// 必須項目の確認
	if cfg.DiscordToken != "test-token" {
		t.Errorf("期待値 DiscordToken='test-token', 実際の値 '%s'", cfg.DiscordToken)
	}
	if cfg.GeminiAPIKey != "test-api-key" {
		t.Errorf("期待値 GeminiAPIKey='test-api-key', 実際の値 '%s'", cfg.GeminiAPIKey)
	}

	// デフォルト値の確認
	if cfg.LogLevel != "info" {
		t.Errorf("期待値 LogLevel='info', 実際の値 '%s'", cfg.LogLevel)
	}
	if cfg.MaxHistory != 20 {
		t.Errorf("期待値 MaxHistory=20, 実際の値 %d", cfg.MaxHistory)
	}
	if cfg.ResponseTimeout != 5*time.Second {
		t.Errorf("期待値 ResponseTimeout=5s, 実際の値 %v", cfg.ResponseTimeout)
	}
}

// TestLoadConfigMissingRequired は必須項目が欠けている場合のテスト
func TestLoadConfigMissingRequired(t *testing.T) {
	// 環境変数を保存してテスト後に復元
	oldToken := os.Getenv("DISCORD_BOT_TOKEN")
	oldAPIKey := os.Getenv("GEMINI_API_KEY")
	defer func() {
		os.Setenv("DISCORD_BOT_TOKEN", oldToken)
		os.Setenv("GEMINI_API_KEY", oldAPIKey)
	}()

	tests := []struct {
		name      string
		token     string
		apiKey    string
		wantError bool
	}{
		{
			name:      "トークンなし",
			token:     "",
			apiKey:    "test-api-key",
			wantError: true,
		},
		{
			name:      "APIキーなし",
			token:     "test-token",
			apiKey:    "",
			wantError: true,
		},
		{
			name:      "両方なし",
			token:     "",
			apiKey:    "",
			wantError: true,
		},
		{
			name:      "両方あり",
			token:     "test-token",
			apiKey:    "test-api-key",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DISCORD_BOT_TOKEN", tt.token)
			os.Setenv("GEMINI_API_KEY", tt.apiKey)

			_, err := Load()
			if (err != nil) != tt.wantError {
				t.Errorf("エラー期待値 %v, 実際のエラー %v", tt.wantError, err)
			}
		})
	}
}

// TestLoadConfigCustomValues はカスタム値の読み込みをテスト
func TestLoadConfigCustomValues(t *testing.T) {
	// 環境変数を保存してテスト後に復元
	oldVars := map[string]string{
		"DISCORD_BOT_TOKEN": os.Getenv("DISCORD_BOT_TOKEN"),
		"GEMINI_API_KEY":    os.Getenv("GEMINI_API_KEY"),
		"LOG_LEVEL":         os.Getenv("LOG_LEVEL"),
		"MAX_HISTORY":       os.Getenv("MAX_HISTORY"),
		"RESPONSE_TIMEOUT":  os.Getenv("RESPONSE_TIMEOUT"),
		"SYSTEM_PROMPT":     os.Getenv("SYSTEM_PROMPT"),
	}
	defer func() {
		for k, v := range oldVars {
			os.Setenv(k, v)
		}
	}()

	// カスタム値を設定
	os.Setenv("DISCORD_BOT_TOKEN", "test-token")
	os.Setenv("GEMINI_API_KEY", "test-api-key")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("MAX_HISTORY", "50")
	os.Setenv("RESPONSE_TIMEOUT", "10s")
	os.Setenv("SYSTEM_PROMPT", "カスタムプロンプト")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("設定の読み込みに失敗: %v", err)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("期待値 LogLevel='debug', 実際の値 '%s'", cfg.LogLevel)
	}
	if cfg.MaxHistory != 50 {
		t.Errorf("期待値 MaxHistory=50, 実際の値 %d", cfg.MaxHistory)
	}
	if cfg.ResponseTimeout != 10*time.Second {
		t.Errorf("期待値 ResponseTimeout=10s, 実際の値 %v", cfg.ResponseTimeout)
	}
	if cfg.SystemPrompt != "カスタムプロンプト" {
		t.Errorf("期待値 SystemPrompt='カスタムプロンプト', 実際の値 '%s'", cfg.SystemPrompt)
	}
}
