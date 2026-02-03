package port

import (
	"imageServer/internal/domain"
)

// HTTPHandler HTTPハンドラーのインターフェース
type HTTPHandler interface {
	// メディア関連
	UploadImage(ctx interface{}) error
	CreateMediaWithYouTube(ctx interface{}) error
	GetMedia(ctx interface{}) error
	ListMedia(ctx interface{}) error
	DeleteMedia(ctx interface{}) error
	
	// タグ関連
	CreateTag(ctx interface{}) error
	GetTag(ctx interface{}) error
	ListTags(ctx interface{}) error
	UpdateTag(ctx interface{}) error
	DeleteTag(ctx interface{}) error
	
	// メディアとタグの関連付け
	AssociateMediaTag(ctx interface{}) error
	RemoveMediaTag(ctx interface{}) error
	GetMediaByTag(ctx interface{}) error
	
	// TODO関連
	CreateTodo(ctx interface{}) error
	GetTodo(ctx interface{}) error
	ListTodos(ctx interface{}) error
	GetTodosByDateRange(ctx interface{}) error
	GetTodosByDate(ctx interface{}) error
	GetTodosWithoutDueDate(ctx interface{}) error
	UpdateTodo(ctx interface{}) error
	DeleteTodo(ctx interface{}) error
}

// CreateMediaRequest メディア作成リクエスト
type CreateMediaRequest struct {
	Type        domain.MediaType `json:"type" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Description *string          `json:"description"`
	TagIDs      []string         `json:"tag_ids"`
}

// CreateMediaWithYouTubeRequest YouTube URL付きメディア作成リクエスト
// @Description YouTube URLでメディアを作成するリクエスト
type CreateMediaWithYouTubeRequest struct {
	YouTubeURL  string   `json:"youtube_url" binding:"required" example:"https://www.youtube.com/watch?v=dQw4w9WgXcQ"`
	Title       string   `json:"title" binding:"required" example:"サンプル動画"`
	Description *string  `json:"description" example:"これはサンプル動画です"`
	TagIDs      []string `json:"tag_ids" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// CreateTagRequest タグ作成リクエスト
// @Description タグを作成するリクエスト
type CreateTagRequest struct {
	Name string            `json:"name" binding:"required" example:"新規タグ"`
	Type domain.TagType    `json:"type" binding:"required" example:"all"`
}

// UpdateTagRequest タグ更新リクエスト
// @Description タグを更新するリクエスト
type UpdateTagRequest struct {
	Name string         `json:"name" binding:"required" example:"更新されたタグ名"`
	Type domain.TagType `json:"type" binding:"required" example:"all"`
}

// AssociateTagRequest タグ関連付けリクエスト
// @Description メディアにタグを関連付けるリクエスト
type AssociateTagRequest struct {
	TagID string `json:"tag_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// CreateTodoRequest TODO作成リクエスト
// @Description TODOを作成するリクエスト
type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required" example:"サンプルTODO"`
	Description *string `json:"description" example:"これはサンプルTODOです"`
	StartDate   *string `json:"start_date" example:"2024-01-01T00:00:00Z"`
	EndDate     *string `json:"end_date" example:"2024-01-31T23:59:59Z"`
	DueDate     *string `json:"due_date" example:"2024-01-15T00:00:00Z"`
}

// UpdateTodoRequest TODO更新リクエスト
// @Description TODOを更新するリクエスト
type UpdateTodoRequest struct {
	Title       string  `json:"title" binding:"required" example:"更新されたTODO"`
	Description *string `json:"description" example:"更新された説明"`
	StartDate   *string `json:"start_date" example:"2024-01-01T00:00:00Z"`
	EndDate     *string `json:"end_date" example:"2024-01-31T23:59:59Z"`
	DueDate     *string `json:"due_date" example:"2024-01-15T00:00:00Z"`
	Completed   bool    `json:"completed" example:"false"`
}
