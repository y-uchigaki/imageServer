package postgres

import (
	"database/sql"
	"imageServer/internal/domain"
	"imageServer/internal/port"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type tagRepository struct {
	db *sql.DB
}

// NewTagRepository タグリポジトリのコンストラクタ
func NewTagRepository(db *sql.DB) port.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *domain.Tag) error {
	query := `
		INSERT INTO tag (id, name, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(
		query,
		tag.ID,
		tag.Name,
		tag.Type,
		tag.CreatedAt,
		tag.UpdatedAt,
	)
	return err
}

func (r *tagRepository) FindByID(id uuid.UUID) (*domain.Tag, error) {
	query := `
		SELECT id, name, type, created_at, updated_at
		FROM tag
		WHERE id = $1
	`
	tag := &domain.Tag{}
	var tagType string
	err := r.db.QueryRow(query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tagType,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	tag.Type = domain.TagType(tagType)
	return tag, nil
}

func (r *tagRepository) FindByName(name string) (*domain.Tag, error) {
	query := `
		SELECT id, name, type, created_at, updated_at
		FROM tag
		WHERE name = $1
	`
	tag := &domain.Tag{}
	var tagType string
	err := r.db.QueryRow(query, name).Scan(
		&tag.ID,
		&tag.Name,
		&tagType,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	tag.Type = domain.TagType(tagType)
	return tag, nil
}

func (r *tagRepository) FindAll() ([]*domain.Tag, error) {
	query := `
		SELECT id, name, type, created_at, updated_at
		FROM tag
		ORDER BY name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		tag := &domain.Tag{}
		var tagType string
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tagType,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tag.Type = domain.TagType(tagType)
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *tagRepository) Update(tag *domain.Tag) error {
	query := `
		UPDATE tag
		SET name = $2, type = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(
		query,
		tag.ID,
		tag.Name,
		tag.Type,
		tag.UpdatedAt,
	)
	return err
}

func (r *tagRepository) Delete(id uuid.UUID) error {
	// 関連するメディアタグを削除
	_, err := r.db.Exec("DELETE FROM media_tag WHERE tag_id = $1", id)
	if err != nil {
		return err
	}

	// タグを削除
	_, err = r.db.Exec("DELETE FROM tag WHERE id = $1", id)
	return err
}
