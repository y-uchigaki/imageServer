package domain

import (
	"time"

	"github.com/google/uuid"
)

// TagType タグの適用可能なメディアタイプ
type TagType string

const (
	TagTypeAll   TagType = "all"   // すべてのメディアタイプで使用可能
	TagTypeImage TagType = "image" // 画像のみ
	TagTypeAudio TagType = "audio" // 音楽のみ
	TagTypeVideo TagType = "video" // YouTube動画のみ
)

// Tag タグエンティティ
type Tag struct {
	ID        uuid.UUID
	Name      string
	Type      TagType // 適用可能なメディアタイプ
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MediaTag メディアとタグの関連エンティティ
type MediaTag struct {
	MediaID uuid.UUID
	TagID   uuid.UUID
}
