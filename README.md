# Image Server

画像・動画保存・表示サービス

## プロジェクト構成

```
imageServer/
├── backend/     # バックエンド（Go + Gin）
└── frontend/    # フロントエンド（今後追加予定）
```

## セットアップ

### ローカル開発環境（Docker Compose + LocalStack）

Docker Composeを使用して、すべてのサービス（PostgreSQL、LocalStack、バックエンドAPI）を一括で起動できます。

#### 1. Docker Composeで全サービスを起動

```bash
docker-compose up -d
```

これにより以下が起動します：
- PostgreSQL（ポート: 5432）
- LocalStack（ポート: 4566）- AWS S3のエミュレーター
- バックエンドAPI（ポート: 8080）

#### 2. サービスの確認

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- LocalStack Dashboard: http://localhost:4566/_localstack/health

#### 3. ログの確認

```bash
# すべてのサービスのログ
docker-compose logs -f

# バックエンドのみのログ
docker-compose logs -f backend
```

### バックエンドのみローカルで実行する場合

Docker ComposeでPostgreSQLとLocalStackだけを起動し、バックエンドはローカルで実行することもできます。

#### 1. Docker Composeでインフラのみ起動

```bash
docker-compose up -d postgres-imageserver localstack-imageserver localstack-init
```

#### 2. バックエンドの環境変数を設定

```bash
cd backend
cp env.example .env
```

#### 3. バックエンドの起動

```bash
cd backend
go run cmd/api/main.go
```

### バックエンド

バックエンドの詳細なセットアップについては、[backend/README.md](./backend/README.md)を参照してください。

### フロントエンド

フロントエンドは今後追加予定です。

## 開発

### Docker Composeコマンド

```bash
# 全サービス起動
docker-compose up -d

# バックエンドのみ再ビルドして起動
docker-compose up -d --build backend

# ログ確認
docker-compose logs -f

# バックエンドのログのみ確認
docker-compose logs -f backend

# サービス停止
docker-compose down

# データを削除して再起動
docker-compose down -v
docker-compose up -d

# バックエンドのコンテナに入る
docker-compose exec backend sh
```

### LocalStackのS3バケット確認

```bash
# S3バケット一覧
aws --endpoint-url=http://localhost:4566 s3 ls

# バケット内のファイル確認
aws --endpoint-url=http://localhost:4566 s3 ls s3://imageserver-bucket
```

## 注意事項

- LocalStackはAWSサービスのローカルエミュレーターです
- 本番環境では実際のAWSサービスを使用してください
- 環境変数は`.env`ファイルで管理し、Gitにコミットしないでください
