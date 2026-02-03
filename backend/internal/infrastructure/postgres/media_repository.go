package postgres

import (
	"database/sql"
	"fmt"
	"imageServer/internal/domain"
	"imageServer/internal/port"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type mediaRepository struct {
	db *sql.DB
}

// NewMediaRepository メディアリポジトリのコンストラクタ
func NewMediaRepository(db *sql.DB) port.MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(media *domain.Media) error {
	// トランザクションを開始
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// メディアをINSERT
	query := `
		INSERT INTO media (id, type, s3_key, youtube_url, cloudfront_url, title, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(
		query,
		media.ID,
		media.Type,
		media.S3Key,
		media.YouTubeURL,
		media.CloudFrontURL,
		media.Title,
		media.Description,
		media.CreatedAt,
		media.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// タグの関連付け（同じトランザクション内で実行）
	for _, tag := range media.Tags {
		// 既存の関連付けをチェック
		var count int
		err = tx.QueryRow(
			"SELECT COUNT(*) FROM media_tag WHERE media_id = $1 AND tag_id = $2",
			media.ID, tag.ID,
		).Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			// 関連付けが存在しない場合のみINSERT
			_, err = tx.Exec(
				"INSERT INTO media_tag (media_id, tag_id) VALUES ($1, $2)",
				media.ID, tag.ID,
			)
			if err != nil {
				return err
			}
		}
	}

	// トランザクションをコミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *mediaRepository) FindByID(id uuid.UUID) (*domain.Media, error) {
	query := `
		SELECT id, type, s3_key, youtube_url, cloudfront_url, title, description, created_at, updated_at
		FROM media
		WHERE id = $1
	`
	media := &domain.Media{}
	var s3Key, youtubeURL, cloudfrontURL, description sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&media.ID,
		&media.Type,
		&s3Key,
		&youtubeURL,
		&cloudfrontURL,
		&media.Title,
		&description,
		&media.CreatedAt,
		&media.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if s3Key.Valid {
		media.S3Key = &s3Key.String
	}
	if youtubeURL.Valid {
		media.YouTubeURL = &youtubeURL.String
	}
	if cloudfrontURL.Valid {
		media.CloudFrontURL = &cloudfrontURL.String
	}
	if description.Valid {
		media.Description = &description.String
	}

	// タグを取得
	tags, err := r.getTagsByMediaID(id)
	if err != nil {
		return nil, err
	}
	media.Tags = tags

	return media, nil
}

func (r *mediaRepository) FindAll() ([]*domain.Media, error) {
	query := `
		SELECT id, type, s3_key, youtube_url, cloudfront_url, title, description, created_at, updated_at
		FROM media
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []*domain.Media
	for rows.Next() {
		media := &domain.Media{}
		var s3Key, youtubeURL, cloudfrontURL, description sql.NullString

		err := rows.Scan(
			&media.ID,
			&media.Type,
			&s3Key,
			&youtubeURL,
			&cloudfrontURL,
			&media.Title,
			&description,
			&media.CreatedAt,
			&media.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if s3Key.Valid {
			media.S3Key = &s3Key.String
		}
		if youtubeURL.Valid {
			media.YouTubeURL = &youtubeURL.String
		}
		if cloudfrontURL.Valid {
			media.CloudFrontURL = &cloudfrontURL.String
		}
		if description.Valid {
			media.Description = &description.String
		}

		// タグを取得
		tags, err := r.getTagsByMediaID(media.ID)
		if err != nil {
			return nil, err
		}
		media.Tags = tags

		mediaList = append(mediaList, media)
	}

	return mediaList, nil
}

func (r *mediaRepository) FindAllWithPagination(offset, limit int) ([]*domain.Media, int, error) {
	// 総件数を取得
	var totalCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM media").Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// ページネーション付きでメディアを取得
	query := `
		SELECT id, type, s3_key, youtube_url, cloudfront_url, title, description, created_at, updated_at
		FROM media
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var mediaList []*domain.Media
	for rows.Next() {
		media := &domain.Media{}
		var s3Key, youtubeURL, cloudfrontURL, description sql.NullString

		err := rows.Scan(
			&media.ID,
			&media.Type,
			&s3Key,
			&youtubeURL,
			&cloudfrontURL,
			&media.Title,
			&description,
			&media.CreatedAt,
			&media.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if s3Key.Valid {
			media.S3Key = &s3Key.String
		}
		if youtubeURL.Valid {
			media.YouTubeURL = &youtubeURL.String
		}
		if cloudfrontURL.Valid {
			media.CloudFrontURL = &cloudfrontURL.String
		}
		if description.Valid {
			media.Description = &description.String
		}

		// タグを取得
		tags, err := r.getTagsByMediaID(media.ID)
		if err != nil {
			return nil, 0, err
		}
		media.Tags = tags

		mediaList = append(mediaList, media)
	}

	return mediaList, totalCount, nil
}

func (r *mediaRepository) FindAllWithFilters(offset, limit int, titleSearch *string, tagIDs []uuid.UUID) ([]*domain.Media, int, error) {
	// WHERE句を構築
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if titleSearch != nil && *titleSearch != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("m.title ILIKE $%d", argIndex))
		args = append(args, "%"+*titleSearch+"%")
		argIndex++
	}

	if len(tagIDs) > 0 {
		// タグIDのプレースホルダーを生成
		placeholders := []string{}
		for _, tagID := range tagIDs {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, tagID)
			argIndex++
		}
		whereConditions = append(whereConditions, fmt.Sprintf(
			"m.id IN (SELECT DISTINCT media_id FROM media_tag WHERE tag_id IN (%s))",
			strings.Join(placeholders, ","),
		))
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// 総件数を取得
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM media m %s", whereClause)
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// ページネーション付きでメディアを取得
	query := fmt.Sprintf(`
		SELECT m.id, m.type, m.s3_key, m.youtube_url, m.cloudfront_url, m.title, m.description, m.created_at, m.updated_at
		FROM media m
		%s
		ORDER BY m.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	
	args = append(args, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	mediaList, err := r.scanMediaList(rows)
	if err != nil {
		return nil, 0, err
	}

	return mediaList, totalCount, nil
}

func (r *mediaRepository) scanMediaList(rows *sql.Rows) ([]*domain.Media, error) {
	var mediaList []*domain.Media
	for rows.Next() {
		media := &domain.Media{}
		var s3Key, youtubeURL, cloudfrontURL, description sql.NullString

		err := rows.Scan(
			&media.ID,
			&media.Type,
			&s3Key,
			&youtubeURL,
			&cloudfrontURL,
			&media.Title,
			&description,
			&media.CreatedAt,
			&media.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if s3Key.Valid {
			media.S3Key = &s3Key.String
		}
		if youtubeURL.Valid {
			media.YouTubeURL = &youtubeURL.String
		}
		if cloudfrontURL.Valid {
			media.CloudFrontURL = &cloudfrontURL.String
		}
		if description.Valid {
			media.Description = &description.String
		}

		// タグを取得
		tags, err := r.getTagsByMediaID(media.ID)
		if err != nil {
			return nil, err
		}
		media.Tags = tags

		mediaList = append(mediaList, media)
	}

	return mediaList, nil
}

func (r *mediaRepository) FindByTagID(tagID uuid.UUID) ([]*domain.Media, error) {
	query := `
		SELECT m.id, m.type, m.s3_key, m.youtube_url, m.cloudfront_url, m.title, m.description, m.created_at, m.updated_at
		FROM media m
		INNER JOIN media_tag mt ON m.id = mt.media_id
		WHERE mt.tag_id = $1
		ORDER BY m.created_at DESC
	`
	rows, err := r.db.Query(query, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []*domain.Media
	for rows.Next() {
		media := &domain.Media{}
		var s3Key, youtubeURL, cloudfrontURL, description sql.NullString

		err := rows.Scan(
			&media.ID,
			&media.Type,
			&s3Key,
			&youtubeURL,
			&cloudfrontURL,
			&media.Title,
			&description,
			&media.CreatedAt,
			&media.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if s3Key.Valid {
			media.S3Key = &s3Key.String
		}
		if youtubeURL.Valid {
			media.YouTubeURL = &youtubeURL.String
		}
		if cloudfrontURL.Valid {
			media.CloudFrontURL = &cloudfrontURL.String
		}
		if description.Valid {
			media.Description = &description.String
		}

		// タグを取得
		tags, err := r.getTagsByMediaID(media.ID)
		if err != nil {
			return nil, err
		}
		media.Tags = tags

		mediaList = append(mediaList, media)
	}

	return mediaList, nil
}

func (r *mediaRepository) Update(media *domain.Media) error {
	query := `
		UPDATE media
		SET type = $2, s3_key = $3, youtube_url = $4, cloudfront_url = $5, title = $6, description = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.Exec(
		query,
		media.ID,
		media.Type,
		media.S3Key,
		media.YouTubeURL,
		media.CloudFrontURL,
		media.Title,
		media.Description,
		time.Now(),
	)
	return err
}

func (r *mediaRepository) Delete(id uuid.UUID) error {
	// 関連するタグを削除
	_, err := r.db.Exec("DELETE FROM media_tag WHERE media_id = $1", id)
	if err != nil {
		return err
	}

	// メディアを削除
	_, err = r.db.Exec("DELETE FROM media WHERE id = $1", id)
	return err
}

func (r *mediaRepository) AssociateTag(mediaID, tagID uuid.UUID) error {
	// 既存の関連付けをチェック
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM media_tag WHERE media_id = $1 AND tag_id = $2",
		mediaID, tagID,
	).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // 既に存在する場合はスキップ
	}

	_, err = r.db.Exec(
		"INSERT INTO media_tag (media_id, tag_id) VALUES ($1, $2)",
		mediaID, tagID,
	)
	return err
}

func (r *mediaRepository) RemoveTag(mediaID, tagID uuid.UUID) error {
	_, err := r.db.Exec(
		"DELETE FROM media_tag WHERE media_id = $1 AND tag_id = $2",
		mediaID, tagID,
	)
	return err
}

func (r *mediaRepository) getTagsByMediaID(mediaID uuid.UUID) ([]domain.Tag, error) {
	query := `
		SELECT t.id, t.name, t.type, t.created_at, t.updated_at
		FROM tag t
		INNER JOIN media_tag mt ON t.id = mt.tag_id
		WHERE mt.media_id = $1
	`
	rows, err := r.db.Query(query, mediaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var tag domain.Tag
		var tagType string
		err := rows.Scan(&tag.ID, &tag.Name, &tagType, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tag.Type = domain.TagType(tagType)
		tags = append(tags, tag)
	}

	return tags, nil
}
