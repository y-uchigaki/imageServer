package http

import (
	"imageServer/internal/port"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// CreateMediaWithYouTubeRequest YouTube URL付きメディア作成リクエスト（Swagger用エイリアス）
type CreateMediaWithYouTubeRequest = port.CreateMediaWithYouTubeRequest

// CreateTagRequest タグ作成リクエスト（Swagger用エイリアス）
type CreateTagRequest = port.CreateTagRequest

// UpdateTagRequest タグ更新リクエスト（Swagger用エイリアス）
type UpdateTagRequest = port.UpdateTagRequest

// AssociateTagRequest タグ関連付けリクエスト（Swagger用エイリアス）
type AssociateTagRequest = port.AssociateTagRequest

// CreateTodoRequest TODO作成リクエスト（Swagger用エイリアス）
type CreateTodoRequest = port.CreateTodoRequest

// UpdateTodoRequest TODO更新リクエスト（Swagger用エイリアス）
type UpdateTodoRequest = port.UpdateTodoRequest

// SetupRouter ルーターをセットアップ
func SetupRouter(handler port.HTTPHandler) *gin.Engine {
	router := gin.Default()

	// CORS設定（必要に応じて調整）
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// メディア関連エンドポイント
	api := router.Group("/api/v1")
	{
		api.POST("/media/upload", UploadImageHandler(handler))
		api.POST("/media/youtube", CreateMediaWithYouTubeHandler(handler))
		api.GET("/media", ListMediaHandler(handler))
		api.GET("/media/:id", GetMediaHandler(handler))
		api.DELETE("/media/:id", DeleteMediaHandler(handler))

		api.POST("/tags", CreateTagHandler(handler))
		api.GET("/tags", ListTagsHandler(handler))
		// より具体的なパスを先に登録（競合を避けるため）
		api.GET("/tags/:id/media", GetMediaByTagHandler(handler))
		api.GET("/tags/:id", GetTagHandler(handler))
		api.PUT("/tags/:id", UpdateTagHandler(handler))
		api.DELETE("/tags/:id", DeleteTagHandler(handler))

		api.POST("/media/:id/tags", AssociateMediaTagHandler(handler))
		api.DELETE("/media/:id/tags/:tag_id", RemoveMediaTagHandler(handler))

		// TODO関連エンドポイント
		api.POST("/todos", CreateTodoHandler(handler))
		api.GET("/todos", ListTodosHandler(handler))
		api.GET("/todos/date-range", GetTodosByDateRangeHandler(handler))
		api.GET("/todos/date", GetTodosByDateHandler(handler))
		api.GET("/todos/without-due-date", GetTodosWithoutDueDateHandler(handler))
		api.GET("/todos/:id", GetTodoHandler(handler))
		api.PUT("/todos/:id", UpdateTodoHandler(handler))
		api.DELETE("/todos/:id", DeleteTodoHandler(handler))
	}

	return router
}
