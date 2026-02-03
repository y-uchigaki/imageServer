# Image Server API

## アーキテクチャ

```
internal/
├── domain/          # ドメイン層（エンティティ、値オブジェクト）
├── application/     # アプリケーション層（ユースケース）
├── port/            # ポート層（インターフェース）
└── infrastructure/  # インフラストラクチャ層（実装）
    ├── postgres/    # PostgreSQL実装
    ├── s3/          # AWS S3実装
    └── http/        # HTTPハンドラー実装
```

## セットアップ

### ローカル開発環境（推奨）

Docker ComposeとLocalStackを使用したローカル開発環境を使用することを推奨します。

#### 1. Docker Composeでサービスを起動

プロジェクトルートで以下を実行：

```bash
docker-compose up -d
```

これにより以下が起動します：
- PostgreSQL（ポート: 5432）
- LocalStack（ポート: 4566）- AWS S3のエミュレーター

#### 2. 環境変数の設定

```bash
cp env.example .env
```

`.env`ファイルには既にLocalStack用の設定が含まれています。

#### 3. 依存関係のインストール

```bash
go mod download
```

#### 4. OpenAPI仕様の生成（初回のみ）

```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

#### 5. サーバーの起動

```bash
go run cmd/api/main.go
```

### 本番環境向けセットアップ

#### 1. 依存関係のインストール

```bash
go mod download
```

#### 2. 環境変数の設定

`.env.example`をコピーして`.env`を作成し、本番環境の値を設定してください。

```bash
cp env.example .env
# .envファイルを編集して本番環境の設定に変更
```

#### 3. データベースの準備

PostgreSQLデータベースを作成してください。

```sql
CREATE DATABASE imageserver;
```

### LocalStackの確認

LocalStackが正常に動作しているか確認：

```bash
# LocalStackのヘルスチェック
curl http://localhost:4566/_localstack/health

# S3バケット一覧
aws --endpoint-url=http://localhost:4566 s3 ls

# バケット内のファイル確認
aws --endpoint-url=http://localhost:4566 s3 ls s3://imageserver-bucket
```

### サーバー起動後のアクセス

サーバー起動後、以下のURLでアクセスできます：

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- OpenAPI JSON: http://localhost:8080/swagger/doc.json
- OpenAPI YAML: http://localhost:8080/swagger/doc.yaml
- LocalStack Dashboard: http://localhost:4566/_localstack/health
