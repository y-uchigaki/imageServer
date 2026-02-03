package domain

import (
	"time"

	"github.com/google/uuid"
)

// MediaType メディアの種類
type MediaType string

const (
	MediaTypeImage  MediaType = "image"
	MediaTypeVideo  MediaType = "video"
	MediaTypeAudio  MediaType = "audio"
)

// Media メディアエンティティ
type Media struct {
	ID          uuid.UUID
	Type        MediaType
	S3Key       *string // 画像の場合のS3キー
	YouTubeURL  *string // YouTube動画の場合のURL
	CloudFrontURL *string // CloudFront経由のURL
	Title       string
	Description *string
	Tags        []Tag
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IsImage 画像かどうか
func (m *Media) IsImage() bool {
	return m.Type == MediaTypeImage
}

// IsVideo 動画かどうか
func (m *Media) IsVideo() bool {
	return m.Type == MediaTypeVideo
}

// IsAudio 音声かどうか
func (m *Media) IsAudio() bool {
	return m.Type == MediaTypeAudio
}
