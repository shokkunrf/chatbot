# Chatbot Infrastructure

このディレクトリには、Discord ChatbotをGCP上にデプロイするためのTerraform構成が含まれてる。

## 構成概要

### ディレクトリ構造

```
infrastructure/
├── environments/        # 環境別の設定
│   ├── dev/            # 開発環境
│   └── prd/            # 本番環境
└── modules/            # 再利用可能なTerraformモジュール
    ├── api/            # GCP APIの有効化
    ├── compute/        # GCEインスタンス
    ├── iam/            # IAMとサービスアカウント
    ├── network/        # VPCネットワーク
    └── secret_manager/ # Secret Manager
```

### リソース命名規則

すべてのリソースは以下の命名規則に従う:

```
chatbot-{環境名}-{リソース名}
```

例:

- `chatbot-dev-instance`
- `chatbot-prd-network`
- `chatbot-dev-discord-bot-token`

## 前提条件

1. **Terraform**: バージョン1.12以上
2. **Google Cloud SDK**: 認証とプロジェクトの設定
3. **GCSバケット**: Terraformステート用(事前作成が必要)
   - バケット名: `{プロジェクトID}-chatbot-tfstate`
4. **必要な権限**: プロジェクトオーナーまたは必要なIAMロール

## セットアップ手順

### 1. GCP認証の設定

```bash
gcloud auth application-default login
gcloud config set project YOUR_PROJECT_ID
```

### 2. Terraformステート用バケットの作成

```bash
# バケットがまだ存在しない場合のみ実行
gsutil mb gs://YOUR_PROJECT_ID-chatbot-tfstate
gsutil versioning set on gs://YOUR_PROJECT_ID-chatbot-tfstate
```

### 3. 環境設定ファイルの準備

各環境のディレクトリで、`terraform.tfvars.template`を`terraform.tfvars`にコピーして値を設定:

```bash
cd environments/dev
cp terraform.tfvars.template terraform.tfvars
# terraform.tfvarsを以下の形式で編集
```

`terraform.tfvars`の設定例:

```hcl
project_id   = "your-gcp-project-id"
environment  = "dev"
docker_image = "gcr.io/your-project/chatbot:latest"
secret_value_discord_bot_token = "your-discord-bot-token-here"
secret_value_gemini_api_key    = "your-gemini-api-key-here"
```

### 4. Terraformの初期化と適用

```bash
# 開発環境の場合
cd environments/dev
terraform init -backend-config="bucket=YOUR_PROJECT_ID-chatbot-tfstate"
terraform plan
terraform apply
```

**重要事項**:

- バケット名はinit時の`-backend-config`で指定
- 環境ごとに異なるバケットを使用することを推奨
- シークレット値は`terraform.tfvars`で管理し、tfstateには平文で保存されない(ハッシュ値のみ)

### 5. シークレット管理

**自動設定(推奨)**:
Terraformが`terraform.tfvars`の値を使用してSecret Managerに自動設定する:

- `null_resource`とlocal-execスクリプトにより実行
- シークレット値は`sensitive = true`属性で保護
- tfstateにはSHA256ハッシュのみ保存(平文は保存されない)

**手動設定(必要時のみ)**:

```bash
# 開発環境
echo -n "YOUR_DISCORD_BOT_TOKEN" | gcloud secrets versions add chatbot-dev-discord-bot-token --data-file=-
echo -n "YOUR_GEMINI_API_KEY" | gcloud secrets versions add chatbot-dev-gemini-api-key --data-file=-
```

## インフラストラクチャ詳細

### Compute Engine

- **OS**: Container-Optimized OS
- **マシンタイプ**: e2-micro
- **ディスク**: 10GB標準永続ディスク
- **ネットワーク**: エフェメラル外部IP、インバウンド拒否、アウトバウンド許可
- **起動スクリプト**:
  - Secret Managerから認証情報を安全に取得
  - Dockerコンテナを起動(標準出力を抑制)
  - エラー時は自動的にインスタンスをシャットダウン
  - GCP Logsへのログ出力設定済み

### ネットワーク

- **VPC**: カスタムVPC(`chatbot-{env}-network`)
- **サブネット**: 10.0.1.0/24
- **ファイアウォール**: エグレスのみ許可

### IAM

- **サービスアカウント**: `chatbot-{env}-compute@{project-id}.iam.gserviceaccount.com`
- **権限**:
  - `roles/secretmanager.secretAccessor`
  - `roles/logging.logWriter`
  - `roles/monitoring.metricWriter`

### Secret Manager

- `chatbot-{env}-discord-bot-token`: Discord Botトークン
- `chatbot-{env}-gemini-api-key`: Gemini APIキー

## モニタリング

- **ログ**: GCPログエクスプローラで確認可能
- **メトリクス**: Cloud Monitoringで確認可能

## トラブルシューティング

### インスタンスが起動しない場合

1. ログエクスプローラで起動スクリプトのログを確認
2. Secret Managerに値が設定されているか確認
3. サービスアカウントの権限を確認

### Dockerコンテナが起動しない場合

1. イメージのURLが正しいか確認
2. イメージがパブリックアクセス可能か確認
3. 環境変数が正しく渡されているか確認
4. startup.shのログを確認(docker run失敗時は自動シャットダウン)
5. Container-Optimized OSの制約を確認

## 環境の削除

```bash
# 環境を削除する場合(注意:すべてのリソースが削除される)
cd environments/dev  # または prd
terraform destroy
```

## セキュリティ設計

- **シークレット管理**: Secret Managerで暗号化保存
- **tfstate 保護**: シークレット値は平文で保存されず、SHA256ハッシュのみ
- **アクセス制御**: 最小権限の原則でサービスアカウントを設定
- **ネットワーク**: インバウンド通信を完全遮断

## 注意事項

- 本番環境への適用前は必ず`terraform plan`で変更内容を確認
- `terraform.tfvars`はGitにコミットしない(.gitignore設定済み)
- バックエンド設定は`terraform init`時のみ必要
- 定期的に Terraform とプロバイダーのバージョンを更新
