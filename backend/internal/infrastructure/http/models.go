package http

// MediaResponse メディアレスポンス
// @Description メディア情報
type MediaResponse struct {
	ID            string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type          string         `json:"type" example:"image" enums:"image,video,audio"`
	Title         string         `json:"title" example:"サンプル画像"`
	Description   *string        `json:"description" example:"これはサンプル画像です"`
	S3Key         *string        `json:"s3_key,omitempty" example:"images/550e8400-e29b-41d4-a716-446655440000.jpg"`
	CloudFrontURL *string        `json:"cloudfront_url,omitempty" example:"https://cloudfront.net/images/550e8400-e29b-41d4-a716-446655440000.jpg"`
	YouTubeURL    *string        `json:"youtube_url,omitempty" example:"https://www.youtube.com/watch?v=dQw4w9WgXcQ"`
	Tags          []TagResponse  `json:"tags"`
	CreatedAt     string         `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     string         `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// TagResponse タグレスポンス
// @Description タグ情報
type TagResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"name" example:"画像"`
	Type      string `json:"type" example:"all" enums:"all,image,audio,video"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// MediaListResponse メディア一覧レスポンス
// @Description メディア一覧
type MediaListResponse struct {
	Media []MediaResponse `json:"media"`
}

// TagListResponse タグ一覧レスポンス
// @Description タグ一覧
type TagListResponse struct {
	Tags []TagResponse `json:"tags"`
}

// ErrorResponse エラーレスポンス
// @Description エラー情報
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// MessageResponse メッセージレスポンス
// @Description メッセージ
type MessageResponse struct {
	Message string `json:"message" example:"success message"`
}

// TodoResponse TODOレスポンス
// @Description TODO情報
type TodoResponse struct {
	ID          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string  `json:"title" example:"サンプルTODO"`
	Description *string `json:"description,omitempty" example:"これはサンプルTODOです"`
	StartDate   *string `json:"start_date,omitempty" example:"2024-01-01T00:00:00Z"`
	EndDate     *string `json:"end_date,omitempty" example:"2024-01-31T23:59:59Z"`
	DueDate     *string `json:"due_date,omitempty" example:"2024-01-15T00:00:00Z"`
	Completed   bool    `json:"completed" example:"false"`
	CreatedAt   string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// TodoListResponse TODO一覧レスポンス
// @Description TODO一覧
type TodoListResponse struct {
	Todos   []TodoResponse `json:"todos"`
	Total   *int            `json:"total,omitempty" example:"100"`
	Offset  *int            `json:"offset,omitempty" example:"0"`
	Limit   *int            `json:"limit,omitempty" example:"20"`
	HasMore *bool           `json:"has_more,omitempty" example:"true"`
}