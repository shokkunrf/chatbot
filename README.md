# Discord Chatbot (バージョン 1.0)

Gemini APIを使用したシンプルなDiscord用チャットボットです。

## 特徴

- 🤖 Google Gemini 2.5 Flash による自然な会話
- 💬 メモリ内での会話履歴管理（最大20メッセージ）
- 🏰 ギルド間の会話履歴の完全な分離
- ⚡ 簡潔な応答（175文字以内）
- 🔧 環境変数による柔軟な設定

## 必要な環境

- Go 1.24.5
- Docker & Docker Compose
- Discord Bot トークン
- Google Gemini API キー

## セットアップ

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd chatbot
```

### 2. 環境変数の設定

```bash
cp .env.template .env
```

`.env` ファイルを編集して必要な値を設定：

```env
# Discord設定 (必須)
DISCORD_BOT_TOKEN=your_discord_bot_token_here

# Gemini API設定 (必須)
GEMINI_API_KEY=your_gemini_api_key_here

# アプリケーション設定 (オプション)
LOG_LEVEL=info                # デフォルト: info
MAX_HISTORY=20                # デフォルト: 20
RESPONSE_TIMEOUT=5s           # デフォルト: 5s
```

### 3. Discord Bot の作成

1. [Discord Developer Portal](https://discord.com/developers/applications) にアクセス
2. 新しいアプリケーションを作成
3. Bot セクションでトークンを生成
4. 必要な権限を設定：
   - Send Messages
   - Read Message History
   - Use Slash Commands

### 4. 起動

```bash
# Docker Compose で起動
docker compose up -d

# ログの確認
docker compose logs -f chatbot
```

## 使い方

### 基本的な使い方

- **メンション**: `@BotName こんにちは`
- **リプライ**: Bot のメッセージに返信

### コマンド

- `/reset` - 現在のチャンネルの会話履歴をリセット

## アーキテクチャ

```
Discord Users
    ↓
Discord Gateway
    ↓
Chatbot Application
├── Discord Handler (イベント処理)
├── Memory Store (会話履歴管理)
└── Gemini Client (AI応答生成)
```

## 開発

### ローカル開発

```bash
cd app
go mod download
go test ./...
go run .
```

### テストの実行

```bash
cd app
go test -v ./...
```

### ディレクトリ構造

```
chatbot/
├── app/                    # アプリケーションコード
│   ├── ai/                # Gemini API クライアント
│   ├── bot/               # Discord bot 実装
│   ├── config/            # 設定管理
│   ├── store/             # メモリストレージ
│   └── main.go            # エントリポイント
├── docs/                   # 設計ドキュメント
├── docker-compose.yml      # Docker Compose設定
├── Dockerfile              # Dockerイメージ定義
└── README.md               # このファイル
```

## トラブルシューティング

### Bot が応答しない

1. Bot がオンラインか確認
2. 適切な権限があるか確認
3. ログでエラーを確認: `docker compose logs chatbot`

### エラーメッセージ

- "申し訳ありません。現在応答できません。" - Gemini API エラー
- "エラーが発生しました。" - その他のシステムエラー

## 今後の機能（バージョン 2.0）

- Firestore による会話履歴の永続化
- 複数のキャラクター対応
- 過去の会話を参照した応答生成

## ライセンス

[LICENSE](LICENSE) ファイルを参照してください。
