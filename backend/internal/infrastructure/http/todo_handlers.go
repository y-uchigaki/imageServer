package http

import (
	"imageServer/internal/port"

	"github.com/gin-gonic/gin"
)

// CreateTodoHandler TODOを作成
// @Summary      TODOを作成
// @Description  TODOを作成します
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        request  body      CreateTodoRequest  true  "リクエスト"
// @Success      201      {object}  TodoResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /todos [post]
func CreateTodoHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.CreateTodo(c)
	}
}

// GetTodoHandler TODOを取得
// @Summary      TODOを取得
// @Description  IDを指定してTODO情報を取得します
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "TODO ID"
// @Success      200  {object}  TodoResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /todos/{id} [get]
func GetTodoHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.GetTodo(c)
	}
}

// ListTodosHandler TODO一覧を取得
// @Summary      TODO一覧を取得
// @Description  すべてのTODOの一覧を取得します
// @Tags         todos
// @Produce      json
// @Success      200  {object}  TodoListResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /todos [get]
func ListTodosHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.ListTodos(c)
	}
}

// GetTodosByDateRangeHandler 日付範囲でTODOを取得
// @Summary      日付範囲でTODOを取得
// @Description  開始日と終了日を指定してTODOを取得します
// @Tags         todos
// @Produce      json
// @Param        start_date  query     string  true  "開始日 (YYYY-MM-DD)"
// @Param        end_date    query     string  true  "終了日 (YYYY-MM-DD)"
// @Success      200         {object}  TodoListResponse
// @Failure      400         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /todos/date-range [get]
func GetTodosByDateRangeHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.GetTodosByDateRange(c)
	}
}

// GetTodosByDateHandler 特定の日付のTODOを取得
// @Summary      特定の日付のTODOを取得
// @Description  指定した日付のTODOを取得します
// @Tags         todos
// @Produce      json
// @Param        date  query     string  false  "日付 (YYYY-MM-DD)"
// @Success      200   {object}  TodoListResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /todos/date [get]
func GetTodosByDateHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.GetTodosByDate(c)
	}
}

// GetTodosWithoutDueDateHandler 期限未設定のTODOを取得
// @Summary      期限未設定のTODOを取得
// @Description  期限が設定されていないTODOをページネーション付きで取得します
// @Tags         todos
// @Produce      json
// @Param        offset  query     int     false  "オフセット"
// @Param        limit   query     int     false  "リミット"
// @Success      200     {object}  TodoListResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /todos/without-due-date [get]
func GetTodosWithoutDueDateHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.GetTodosWithoutDueDate(c)
	}
}

// UpdateTodoHandler TODOを更新
// @Summary      TODOを更新
// @Description  IDを指定してTODO情報を更新します
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id       path      string           true  "TODO ID"
// @Param        request  body      UpdateTodoRequest  true  "リクエスト"
// @Success      200      {object}  TodoResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /todos/{id} [put]
func UpdateTodoHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.UpdateTodo(c)
	}
}

// DeleteTodoHandler TODOを削除
// @Summary      TODOを削除
// @Description  IDを指定してTODOを削除します
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "TODO ID"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /todos/{id} [delete]
func DeleteTodoHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.DeleteTodo(c)
	}
}
