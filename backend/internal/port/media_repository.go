package port

import (
	"imageServer/internal/domain"

	"github.com/google/uuid"
)

// MediaRepository メディアリポジトリのインターフェース
type MediaRepository interface {
	Create(media *domain.Media) error
	FindByID(id uuid.UUID) (*domain.Media, error)
	FindAll() ([]*domain.Media, error)
	FindAllWithPagination(offset, limit int) ([]*domain.Media, int, error)
	FindAllWithFilters(offset, limit int, titleSearch *string, tagIDs []uuid.UUID) ([]*domain.Media, int, error)
	FindByTagID(tagID uuid.UUID) ([]*domain.Media, error)
	Update(media *domain.Media) error
	Delete(id uuid.UUID) error
	AssociateTag(mediaID, tagID uuid.UUID) error
	RemoveTag(mediaID, tagID uuid.UUID) error
}
