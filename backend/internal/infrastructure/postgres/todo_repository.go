package postgres

import (
	"database/sql"
	"imageServer/internal/domain"
	"imageServer/internal/port"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type todoRepository struct {
	db *sql.DB
}

// NewTodoRepository TODOリポジトリのコンストラクタ
func NewTodoRepository(db *sql.DB) port.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *domain.Todo) error {
	query := `
		INSERT INTO todo (id, title, description, start_date, end_date, due_date, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(
		query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.StartDate,
		todo.EndDate,
		todo.DueDate,
		todo.Completed,
		todo.CreatedAt,
		todo.UpdatedAt,
	)
	return err
}

func (r *todoRepository) FindByID(id uuid.UUID) (*domain.Todo, error) {
	query := `
		SELECT id, title, description, start_date, end_date, due_date, completed, created_at, updated_at
		FROM todo
		WHERE id = $1
	`
	todo := &domain.Todo{}
	var description sql.NullString
	var startDate, endDate, dueDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&description,
		&startDate,
		&endDate,
		&dueDate,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if description.Valid {
		todo.Description = &description.String
	}
	if startDate.Valid {
		todo.StartDate = &startDate.Time
	}
	if endDate.Valid {
		todo.EndDate = &endDate.Time
	}
	if dueDate.Valid {
		todo.DueDate = &dueDate.Time
	}

	return todo, nil
}

func (r *todoRepository) FindAll() ([]*domain.Todo, error) {
	query := `
		SELECT id, title, description, start_date, end_date, due_date, completed, created_at, updated_at
		FROM todo
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

func (r *todoRepository) FindByDateRange(startDate, endDate time.Time) ([]*domain.Todo, error) {
	query := `
		SELECT id, title, description, start_date, end_date, due_date, completed, created_at, updated_at
		FROM todo
		WHERE (
			(start_date IS NOT NULL AND end_date IS NOT NULL AND start_date <= $2 AND end_date >= $1)
			OR (due_date IS NOT NULL AND due_date >= $1 AND due_date <= $2)
		)
		AND completed = FALSE
		ORDER BY COALESCE(start_date, due_date) ASC
	`
	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

func (r *todoRepository) FindByDate(date time.Time) ([]*domain.Todo, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-1 * time.Second)

	return r.FindByDateRange(startOfDay, endOfDay)
}

func (r *todoRepository) FindWithoutDueDate(offset, limit int) ([]*domain.Todo, int, error) {
	// 総件数を取得（期限未設定で未完了、または期限切れでないもの）
	var totalCount int
	countQuery := `
		SELECT COUNT(*) FROM todo
		WHERE start_date IS NULL AND end_date IS NULL AND due_date IS NULL
		AND completed = FALSE
	`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// ページネーション付きで取得
	query := `
		SELECT id, title, description, start_date, end_date, due_date, completed, created_at, updated_at
		FROM todo
		WHERE start_date IS NULL AND end_date IS NULL AND due_date IS NULL
		AND completed = FALSE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	todos, err := r.scanTodos(rows)
	if err != nil {
		return nil, 0, err
	}

	return todos, totalCount, nil
}

func (r *todoRepository) Update(todo *domain.Todo) error {
	query := `
		UPDATE todo
		SET title = $2, description = $3, start_date = $4, end_date = $5, due_date = $6, completed = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.Exec(
		query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.StartDate,
		todo.EndDate,
		todo.DueDate,
		todo.Completed,
		todo.UpdatedAt,
	)
	return err
}

func (r *todoRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM todo WHERE id = $1", id)
	return err
}

func (r *todoRepository) scanTodos(rows *sql.Rows) ([]*domain.Todo, error) {
	var todos []*domain.Todo
	for rows.Next() {
		todo := &domain.Todo{}
		var description sql.NullString
		var startDate, endDate, dueDate sql.NullTime

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&description,
			&startDate,
			&endDate,
			&dueDate,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			todo.Description = &description.String
		}
		if startDate.Valid {
			todo.StartDate = &startDate.Time
		}
		if endDate.Valid {
			todo.EndDate = &endDate.Time
		}
		if dueDate.Valid {
			todo.DueDate = &dueDate.Time
		}

		todos = append(todos, todo)
	}

	return todos, nil
}
