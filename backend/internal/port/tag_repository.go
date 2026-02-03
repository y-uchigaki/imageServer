package port

import (
	"imageServer/internal/domain"

	"github.com/google/uuid"
)

// TagRepository タグリポジトリのインターフェース
type TagRepository interface {
	Create(tag *domain.Tag) error
	FindByID(id uuid.UUID) (*domain.Tag, error)
	FindByName(name string) (*domain.Tag, error)
	FindAll() ([]*domain.Tag, error)
	Update(tag *domain.Tag) error
	Delete(id uuid.UUID) error
}
