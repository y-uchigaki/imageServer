package domain

import (
	"time"

	"github.com/google/uuid"
)

// Todo TODOエンティティ
type Todo struct {
	ID          uuid.UUID
	Title       string
	Description *string
	StartDate   *time.Time // 開始日（期間指定の場合）
	EndDate     *time.Time // 終了日（期間指定の場合）
	DueDate     *time.Time // 期限日（単体指定の場合）
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// HasPeriod 期間が設定されているか
func (t *Todo) HasPeriod() bool {
	return t.StartDate != nil && t.EndDate != nil
}

// HasDueDate 期限日が設定されているか
func (t *Todo) HasDueDate() bool {
	return t.DueDate != nil
}

// IsOverdue 期限切れかどうか
func (t *Todo) IsOverdue() bool {
	now := time.Now()
	if t.DueDate != nil {
		return t.DueDate.Before(now) && !t.Completed
	}
	if t.EndDate != nil {
		return t.EndDate.Before(now) && !t.Completed
	}
	return false
}
