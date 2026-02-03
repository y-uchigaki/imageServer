package main

import (
	"database/sql"
	"fmt"
	"imageServer/internal/application"
	"imageServer/internal/infrastructure/http"
	"imageServer/internal/infrastructure/postgres"
	"imageServer/internal/infrastructure/s3"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	
	_ "imageServer/docs" // Swagger docs
)

// @title           Image Server API
// @version         1.0
// @description     画像・動画保存・表示サービスのAPI
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https
func main() {
	// 環境変数の読み込み
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// データベース接続
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// マイグレーション実行
	if err := postgres.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 初期タグの投入
	if err := postgres.SeedInitialTags(db); err != nil {
		log.Fatalf("Failed to seed initial tags: %v", err)
	}

	// S3サービスの初期化
	s3Service, err := s3.NewS3Service()
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}

	// リポジトリの初期化
	mediaRepo := postgres.NewMediaRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	todoRepo := postgres.NewTodoRepository(db)

	// サービスの初期化
	mediaService := application.NewMediaService(mediaRepo, tagRepo, s3Service)
	tagService := application.NewTagService(tagRepo)
	todoService := application.NewTodoService(todoRepo)

	// HTTPハンドラーの初期化
	handler := http.NewHandler(mediaService, tagService, todoService)

	// ルーターのセットアップ
	router := http.SetupRouter(handler)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
