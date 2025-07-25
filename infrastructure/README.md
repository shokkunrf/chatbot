# Chatbot Infrastructure

このディレクトリには、Discord ChatbotをGCP上にデプロイするためのTerraform構成が含まれています。

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

すべてのリソースは以下の命名規則に従います：
```
chatbot-{環境名}-{リソース名}
```

例：
- `chatbot-dev-instance`
- `chatbot-prd-network`
- `chatbot-dev-discord-bot-token`

## 前提条件

1. **Terraform**: バージョン1.12以上
2. **Google Cloud SDK**: 認証とプロジェクトの設定
3. **GCSバケット**: Terraformステート用（事前作成が必要）
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

### 3. 環境変数の設定

各環境のディレクトリで、`terraform.tfvars.template`を`terraform.tfvars`にコピーして値を設定：

```bash
cd environments/dev
cp terraform.tfvars.template terraform.tfvars
# terraform.tfvarsを編集して適切な値を設定
```

### 4. Terraformの初期化と適用

```bash
# 開発環境の場合
cd environments/dev
terraform init -backend-config="bucket=YOUR_PROJECT_ID-chatbot-tfstate"
terraform plan
terraform apply

# 本番環境の場合
cd environments/prd
terraform init -backend-config="bucket=YOUR_PROJECT_ID-chatbot-tfstate"
terraform plan
terraform apply
```

### 5. シークレットの設定

Terraform適用後、以下のコマンドで認証情報を設定：

```bash
# 開発環境
echo -n "YOUR_DISCORD_BOT_TOKEN" | gcloud secrets versions add chatbot-dev-discord-bot-token --data-file=-
echo -n "YOUR_GEMINI_API_KEY" | gcloud secrets versions add chatbot-dev-gemini-api-key --data-file=-

# 本番環境
echo -n "YOUR_DISCORD_BOT_TOKEN" | gcloud secrets versions add chatbot-prd-discord-bot-token --data-file=-
echo -n "YOUR_GEMINI_API_KEY" | gcloud secrets versions add chatbot-prd-gemini-api-key --data-file=-
```

## インフラストラクチャ詳細

### Compute Engine

- **OS**: Container-Optimized OS
- **マシンタイプ**: e2-micro
- **ディスク**: 10GB標準永続ディスク
- **ネットワーク**: エフェメラル外部IP、インバウンド拒否、アウトバウンド許可
- **起動スクリプト**: Secret Managerから認証情報を取得してDockerコンテナを起動

### ネットワーク

- **VPC**: カスタムVPC（`chatbot-{env}-network`）
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

## 環境の削除

```bash
# 環境を削除する場合（注意：すべてのリソースが削除されます）
cd environments/dev  # または prd
terraform destroy
```

## 注意事項

- 本番環境への適用前は必ず`terraform plan`で変更内容を確認
- シークレットは絶対にコードにハードコーディングしない
- 定期的にTerraformとプロバイダーのバージョンを更新
