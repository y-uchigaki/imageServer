package application

import (
	"fmt"
	"imageServer/internal/domain"
	"imageServer/internal/port"
	"time"

	"github.com/google/uuid"
)

// TodoService TODOサービスのユースケース
type TodoService struct {
	todoRepo port.TodoRepository
}

// NewTodoService TODOサービスのコンストラクタ
func NewTodoService(todoRepo port.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

// CreateTodo TODOを作成
func (s *TodoService) CreateTodo(title string, description *string, startDate, endDate, dueDate *time.Time) (*domain.Todo, error) {
	now := time.Now()
	todo := &domain.Todo{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		DueDate:     dueDate,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.todoRepo.Create(todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

// GetTodo TODOを取得
func (s *TodoService) GetTodo(id uuid.UUID) (*domain.Todo, error) {
	todo, err := s.todoRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find todo: %w", err)
	}

	return todo, nil
}

// ListTodos TODO一覧を取得
func (s *TodoService) ListTodos() ([]*domain.Todo, error) {
	todos, err := s.todoRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	return todos, nil
}

// GetTodosByDateRange 日付範囲でTODOを取得
func (s *TodoService) GetTodosByDateRange(startDate, endDate time.Time) ([]*domain.Todo, error) {
	todos, err := s.todoRepo.FindByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos by date range: %w", err)
	}

	return todos, nil
}

// GetTodosByDate 特定の日付のTODOを取得
func (s *TodoService) GetTodosByDate(date time.Time) ([]*domain.Todo, error) {
	todos, err := s.todoRepo.FindByDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos by date: %w", err)
	}

	return todos, nil
}

// GetTodosWithoutDueDate 期限未設定のTODOを取得（ページネーション付き）
func (s *TodoService) GetTodosWithoutDueDate(offset, limit int) ([]*domain.Todo, int, error) {
	todos, totalCount, err := s.todoRepo.FindWithoutDueDate(offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get todos without due date: %w", err)
	}

	return todos, totalCount, nil
}

// UpdateTodo TODOを更新
func (s *TodoService) UpdateTodo(id uuid.UUID, title string, description *string, startDate, endDate, dueDate *time.Time, completed bool) (*domain.Todo, error) {
	todo, err := s.todoRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find todo: %w", err)
	}

	todo.Title = title
	todo.Description = description
	todo.StartDate = startDate
	todo.EndDate = endDate
	todo.DueDate = dueDate
	todo.Completed = completed
	todo.UpdatedAt = time.Now()

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return todo, nil
}

// DeleteTodo TODOを削除
func (s *TodoService) DeleteTodo(id uuid.UUID) error {
	if err := s.todoRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}
