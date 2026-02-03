package application

import (
	"fmt"
	"imageServer/internal/domain"
	"imageServer/internal/port"
	"time"

	"github.com/google/uuid"
)

// TagService タグサービスのユースケース
type TagService struct {
	tagRepo port.TagRepository
}

// NewTagService タグサービスのコンストラクタ
func NewTagService(tagRepo port.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// CreateTag タグを作成
func (s *TagService) CreateTag(name string, tagType domain.TagType) (*domain.Tag, error) {
	// 既存のタグをチェック
	existing, err := s.tagRepo.FindByName(name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("tag with name '%s' already exists", name)
	}

	// デフォルト値の設定
	if tagType == "" {
		tagType = domain.TagTypeAll
	}

	now := time.Now()
	tag := &domain.Tag{
		ID:        uuid.New(),
		Name:      name,
		Type:      tagType,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

// GetTag タグを取得
func (s *TagService) GetTag(id uuid.UUID) (*domain.Tag, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	return tag, nil
}

// ListTags タグ一覧を取得
func (s *TagService) ListTags() ([]*domain.Tag, error) {
	tags, err := s.tagRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// UpdateTag タグを更新
func (s *TagService) UpdateTag(id uuid.UUID, name string, tagType domain.TagType) (*domain.Tag, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	tag.Name = name
	if tagType != "" {
		tag.Type = tagType
	}
	tag.UpdatedAt = time.Now()

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return tag, nil
}

// DeleteTag タグを削除
func (s *TagService) DeleteTag(id uuid.UUID) error {
	if err := s.tagRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}
