package port

import (
	"imageServer/internal/domain"
	"time"

	"github.com/google/uuid"
)

// TodoRepository TODOリポジトリのインターフェース
type TodoRepository interface {
	Create(todo *domain.Todo) error
	FindByID(id uuid.UUID) (*domain.Todo, error)
	FindAll() ([]*domain.Todo, error)
	FindByDateRange(startDate, endDate time.Time) ([]*domain.Todo, error)
	FindWithoutDueDate(offset, limit int) ([]*domain.Todo, int, error)
	FindByDate(date time.Time) ([]*domain.Todo, error)
	Update(todo *domain.Todo) error
	Delete(id uuid.UUID) error
}
