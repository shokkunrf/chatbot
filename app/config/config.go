package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Config はアプリケーションの設定を保持する
type Config struct {
	// Discord設定
	DiscordToken string

	// Gemini API設定
	GeminiAPIKey string

	// アプリケーション設定
	LogLevel        string
	MaxHistory      int
	ResponseTimeout time.Duration
	SystemPrompt    string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	cfg := &Config{
		// 必須項目
		DiscordToken: os.Getenv("DISCORD_BOT_TOKEN"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),

		// デフォルト値を持つ項目
		LogLevel:        getEnvOrDefault("LOG_LEVEL", "info"),
		MaxHistory:      getEnvAsIntOrDefault("MAX_HISTORY", 20),
		ResponseTimeout: getEnvAsDurationOrDefault("RESPONSE_TIMEOUT", 5*time.Second),
	}

	// システムプロンプトの設定
	cfg.SystemPrompt = getEnvOrDefault("SYSTEM_PROMPT", getDefaultSystemPrompt())

	// 必須項目の検証
	if cfg.DiscordToken == "" {
		return nil, errors.New("DISCORD_BOT_TOKEN環境変数が設定されていません")
	}
	if cfg.GeminiAPIKey == "" {
		return nil, errors.New("GEMINI_API_KEY環境変数が設定されていません")
	}

	// 値の検証
	if cfg.MaxHistory < 1 {
		return nil, errors.New("MAX_HISTORYは1以上である必要があります")
	}
	if cfg.ResponseTimeout < 1*time.Second {
		return nil, errors.New("RESPONSE_TIMEOUTは1秒以上である必要があります")
	}

	return cfg, nil
}

// getEnvOrDefault は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault は環境変数を整数として取得し、存在しない場合はデフォルト値を返す
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsDurationOrDefault は環境変数を時間として取得し、存在しない場合はデフォルト値を返す
func getEnvAsDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getDefaultSystemPrompt はデフォルトのシステムプロンプトを返す
func getDefaultSystemPrompt() string {
	return `あなたはDiscordサーバーで活動するフレンドリーなアシスタントです。
以下のルールに従って応答してください：

1. 応答は必ず175文字以内に収めてください
2. 自然な会話を心がけ、Discord特有の表現も適切に使用してください
3. ユーザーのメッセージの文脈を理解し、適切に応答してください
4. 要点を簡潔にまとめ、必要な情報を優先して伝えてください`
}
