# ビルドステージ
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Swagger docsを生成
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# 環境変数
ENV TZ=Asia/Tokyo

EXPOSE 8080

CMD ["./main"]
