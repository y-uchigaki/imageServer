package postgres

import (
	"database/sql"
	"fmt"
)

// Migrate データベースマイグレーションを実行
func Migrate(db *sql.DB) error {
	queries := []string{
		// タグテーブル
		`CREATE TABLE IF NOT EXISTS tag (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			type VARCHAR(50) NOT NULL DEFAULT 'all',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		// タグテーブルにtypeカラムを追加（既存テーブル用 - 安全な方法）
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'tag' AND column_name = 'type'
			) THEN
				ALTER TABLE tag ADD COLUMN type VARCHAR(50);
				UPDATE tag SET type = 'all' WHERE type IS NULL;
				ALTER TABLE tag ALTER COLUMN type SET DEFAULT 'all';
				ALTER TABLE tag ALTER COLUMN type SET NOT NULL;
			END IF;
		END $$`,
		// メディアテーブル
		`CREATE TABLE IF NOT EXISTS media (
			id UUID PRIMARY KEY,
			type VARCHAR(50) NOT NULL,
			s3_key VARCHAR(500),
			youtube_url TEXT,
			cloudfront_url TEXT,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		// メディアとタグの関連テーブル
		`CREATE TABLE IF NOT EXISTS media_tag (
			media_id UUID NOT NULL,
			tag_id UUID NOT NULL,
			PRIMARY KEY (media_id, tag_id),
			FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tag(id) ON DELETE CASCADE
		)`,
		// インデックス
		`CREATE INDEX IF NOT EXISTS idx_media_type ON media(type)`,
		`CREATE INDEX IF NOT EXISTS idx_media_created_at ON media(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_media_tag_media_id ON media_tag(media_id)`,
		`CREATE INDEX IF NOT EXISTS idx_media_tag_tag_id ON media_tag(tag_id)`,
		// タグテーブルのインデックス
		`CREATE INDEX IF NOT EXISTS idx_tag_type ON tag(type)`,
		`CREATE INDEX IF NOT EXISTS idx_tag_created_at ON tag(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_tag_name ON tag(name)`,
		// TODOテーブル
		`CREATE TABLE IF NOT EXISTS todo (
			id UUID PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			start_date TIMESTAMP,
			end_date TIMESTAMP,
			due_date TIMESTAMP,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		// TODOテーブルのインデックス
		`CREATE INDEX IF NOT EXISTS idx_todo_start_date ON todo(start_date)`,
		`CREATE INDEX IF NOT EXISTS idx_todo_end_date ON todo(end_date)`,
		`CREATE INDEX IF NOT EXISTS idx_todo_due_date ON todo(due_date)`,
		`CREATE INDEX IF NOT EXISTS idx_todo_completed ON todo(completed)`,
		`CREATE INDEX IF NOT EXISTS idx_todo_created_at ON todo(created_at)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

// SeedInitialTags 初期タグ（動画、画像）を投入
func SeedInitialTags(db *sql.DB) error {
	tags := []struct {
		id   string
		name string
		tagType string
	}{
		{"00000000-0000-0000-0000-000000000001", "画像", "image"},
		{"00000000-0000-0000-0000-000000000002", "動画", "video"},
	}

	for _, tag := range tags {
		query := `
			INSERT INTO tag (id, name, type, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type
		`
		if _, err := db.Exec(query, tag.id, tag.name, tag.tagType); err != nil {
			return fmt.Errorf("failed to seed tag %s: %w", tag.name, err)
		}
	}

	return nil
}
