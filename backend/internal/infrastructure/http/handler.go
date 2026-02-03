package http

import (
	"fmt"
	"imageServer/internal/application"
	"imageServer/internal/domain"
	"imageServer/internal/port"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	mediaService *application.MediaService
	tagService   *application.TagService
	todoService  *application.TodoService
}

// NewHandler HTTPハンドラーのコンストラクタ
func NewHandler(mediaService *application.MediaService, tagService *application.TagService, todoService *application.TodoService) port.HTTPHandler {
	return &handler{
		mediaService: mediaService,
		tagService:   tagService,
		todoService:  todoService,
	}
}

// UploadImage 画像をアップロード
func (h *handler) UploadImage(ctx interface{}) error {
	c := ctx.(*gin.Context)

	// マルチパートフォームを解析
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return err
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return fmt.Errorf("title is required")
	}

	description := c.PostForm("description")
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	// タグIDを取得
	tagIDsStr := c.PostFormArray("tag_ids")
	var tagIDs []uuid.UUID
	for _, tagIDStr := range tagIDsStr {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid tag_id: %s", tagIDStr)})
			return err
		}
		tagIDs = append(tagIDs, tagID)
	}

	// ファイルを開く
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return err
	}
	defer src.Close()

	// ファイル内容を読み込む
	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return err
	}

	// ファイルタイプを判定
	ext := filepath.Ext(file.Filename)
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	var media *domain.Media
	var s3Key string

	// 音楽ファイルかどうかを判定
	isAudio := ext == ".mp3" || ext == ".wav" || ext == ".wave" ||
		contentType == "audio/mpeg" || contentType == "audio/mp3" ||
		contentType == "audio/wav" || contentType == "audio/wave" ||
		contentType == "audio/x-wav"

	if isAudio {
		// 音楽ファイルの場合
		s3Key = fmt.Sprintf("audio/%s%s", uuid.New().String(), ext)
		if err := h.mediaService.UploadImageToS3(s3Key, data, contentType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload to S3: %v", err)})
			return err
		}
		media, err = h.mediaService.CreateAudioMedia(s3Key, title, descPtr, tagIDs)
	} else {
		// 画像ファイルの場合
		s3Key = fmt.Sprintf("images/%s%s", uuid.New().String(), ext)
		if err := h.mediaService.UploadImageToS3(s3Key, data, contentType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload to S3: %v", err)})
			return err
		}
		media, err = h.mediaService.CreateImageMedia(s3Key, title, descPtr, tagIDs)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create media: %v", err)})
		return err
	}

	c.JSON(http.StatusCreated, toMediaResponse(media))
	return nil
}

// CreateMediaWithYouTube YouTube URLでメディアを作成
func (h *handler) CreateMediaWithYouTube(ctx interface{}) error {
	c := ctx.(*gin.Context)

	var req port.CreateMediaWithYouTubeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	// タグIDをパース
	var tagIDs []uuid.UUID
	for _, tagIDStr := range req.TagIDs {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid tag_id: %s", tagIDStr)})
			return err
		}
		tagIDs = append(tagIDs, tagID)
	}

	media, err := h.mediaService.CreateYouTubeMedia(req.YouTubeURL, req.Title, req.Description, tagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create media: %v", err)})
		return err
	}

	c.JSON(http.StatusCreated, toMediaResponse(media))
	return nil
}

// GetMedia メディアを取得
func (h *handler) GetMedia(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	media, err := h.mediaService.GetMedia(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return err
	}

	c.JSON(http.StatusOK, toMediaResponse(media))
	return nil
}

// ListMedia メディア一覧を取得
func (h *handler) ListMedia(ctx interface{}) error {
	c := ctx.(*gin.Context)

	// ページネーションパラメータを取得
	offset := 0
	limit := 20 // デフォルト値
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	// フィルターパラメータを取得
	titleSearch := c.Query("title")
	var titleSearchPtr *string
	if titleSearch != "" {
		titleSearchPtr = &titleSearch
	}

	tagIDsStr := c.QueryArray("tag_ids")
	var tagIDs []uuid.UUID
	for _, tagIDStr := range tagIDsStr {
		if tagID, err := uuid.Parse(tagIDStr); err == nil {
			tagIDs = append(tagIDs, tagID)
		}
	}

	// フィルターまたはページネーションが指定されている場合
	hasFilters := titleSearchPtr != nil || len(tagIDs) > 0
	hasPagination := offset > 0 || limit != 20 || c.Query("offset") != "" || c.Query("limit") != ""

	if hasFilters || hasPagination {
		var mediaList []*domain.Media
		var totalCount int
		var err error

		if hasFilters {
			mediaList, totalCount, err = h.mediaService.ListMediaWithFilters(offset, limit, titleSearchPtr, tagIDs)
		} else {
			mediaList, totalCount, err = h.mediaService.ListMediaWithPagination(offset, limit)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to list media: %v", err)})
			return err
		}

		responses := make([]map[string]interface{}, len(mediaList))
		for i, media := range mediaList {
			responses[i] = toMediaResponse(media)
		}

		hasMore := offset+limit < totalCount
		c.JSON(http.StatusOK, gin.H{
			"media":     responses,
			"total":     totalCount,
			"offset":    offset,
			"limit":     limit,
			"has_more":  hasMore,
		})
		return nil
	}

	// フィルターもページネーションも指定されていない場合は全件取得（後方互換性のため）
	mediaList, err := h.mediaService.ListMedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to list media: %v", err)})
		return err
	}

	responses := make([]map[string]interface{}, len(mediaList))
	for i, media := range mediaList {
		responses[i] = toMediaResponse(media)
	}

	c.JSON(http.StatusOK, gin.H{"media": responses})
	return nil
}

// DeleteMedia メディアを削除
func (h *handler) DeleteMedia(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	if err := h.mediaService.DeleteMedia(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete media: %v", err)})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "media deleted successfully"})
	return nil
}

// CreateTag タグを作成
func (h *handler) CreateTag(ctx interface{}) error {
	c := ctx.(*gin.Context)

	var req port.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	tagType := req.Type
	if tagType == "" {
		tagType = domain.TagTypeAll
	}
	tag, err := h.tagService.CreateTag(req.Name, tagType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create tag: %v", err)})
		return err
	}

	c.JSON(http.StatusCreated, toTagResponse(tag))
	return nil
}

// GetTag タグを取得
func (h *handler) GetTag(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	tag, err := h.tagService.GetTag(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return err
	}

	c.JSON(http.StatusOK, toTagResponse(tag))
	return nil
}

// ListTags タグ一覧を取得
func (h *handler) ListTags(ctx interface{}) error {
	c := ctx.(*gin.Context)

	tags, err := h.tagService.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to list tags: %v", err)})
		return err
	}

	responses := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		responses[i] = toTagResponse(tag)
	}

	c.JSON(http.StatusOK, gin.H{"tags": responses})
	return nil
}

// UpdateTag タグを更新
func (h *handler) UpdateTag(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	var req port.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	tagType := req.Type
	if tagType == "" {
		// 既存のタグを取得してタイプを保持
		existingTag, err := h.tagService.GetTag(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get tag: %v", err)})
			return err
		}
		tagType = existingTag.Type
	}
	tag, err := h.tagService.UpdateTag(id, req.Name, tagType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update tag: %v", err)})
		return err
	}

	c.JSON(http.StatusOK, toTagResponse(tag))
	return nil
}

// DeleteTag タグを削除
func (h *handler) DeleteTag(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	if err := h.tagService.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete tag: %v", err)})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})
	return nil
}

// AssociateMediaTag メディアにタグを関連付け
func (h *handler) AssociateMediaTag(ctx interface{}) error {
	c := ctx.(*gin.Context)

	mediaIDStr := c.Param("id")
	mediaID, err := uuid.Parse(mediaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media id"})
		return err
	}

	var req port.AssociateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	tagID, err := uuid.Parse(req.TagID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return err
	}

	if err := h.mediaService.AssociateTag(mediaID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to associate tag: %v", err)})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag associated successfully"})
	return nil
}

// RemoveMediaTag メディアからタグを削除
func (h *handler) RemoveMediaTag(ctx interface{}) error {
	c := ctx.(*gin.Context)

	mediaIDStr := c.Param("id")
	mediaID, err := uuid.Parse(mediaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media id"})
		return err
	}

	tagIDStr := c.Param("tag_id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return err
	}

	if err := h.mediaService.RemoveTag(mediaID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to remove tag: %v", err)})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag removed successfully"})
	return nil
}

// GetMediaByTag タグでメディアを取得
func (h *handler) GetMediaByTag(ctx interface{}) error {
	c := ctx.(*gin.Context)

	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return err
	}

	mediaList, err := h.mediaService.GetMediaByTag(tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get media by tag: %v", err)})
		return err
	}

	responses := make([]map[string]interface{}, len(mediaList))
	for i, media := range mediaList {
		responses[i] = toMediaResponse(media)
	}

	c.JSON(http.StatusOK, gin.H{"media": responses})
	return nil
}

// レスポンス変換関数
func toMediaResponse(media *domain.Media) map[string]interface{} {
	tags := make([]map[string]interface{}, len(media.Tags))
	for i, tag := range media.Tags {
		tags[i] = toTagResponse(&tag)
	}

	resp := map[string]interface{}{
		"id":          media.ID.String(),
		"type":        string(media.Type),
		"title":       media.Title,
		"description": media.Description,
		"tags":        tags,
		"created_at":  media.CreatedAt.Format(time.RFC3339),
		"updated_at":  media.UpdatedAt.Format(time.RFC3339),
	}

	if media.S3Key != nil {
		resp["s3_key"] = *media.S3Key
	}
	if media.CloudFrontURL != nil {
		resp["cloudfront_url"] = *media.CloudFrontURL
	}
	if media.YouTubeURL != nil {
		resp["youtube_url"] = *media.YouTubeURL
	}

	return resp
}

func toTagResponse(tag *domain.Tag) map[string]interface{} {
	return map[string]interface{}{
		"id":         tag.ID.String(),
		"name":       tag.Name,
		"type":       string(tag.Type),
		"created_at": tag.CreatedAt.Format(time.RFC3339),
		"updated_at": tag.UpdatedAt.Format(time.RFC3339),
	}
}

// parseTime 文字列をtime.Timeに変換
func parseTime(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// formatTime time.Timeを文字列に変換
func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// CreateTodo TODOを作成
func (h *handler) CreateTodo(ctx interface{}) error {
	c := ctx.(*gin.Context)

	var req port.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	startDate, err := parseTime(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return err
	}
	endDate, err := parseTime(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return err
	}
	dueDate, err := parseTime(req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format"})
		return err
	}

	todo, err := h.todoService.CreateTodo(req.Title, req.Description, startDate, endDate, dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create todo: %v", err)})
		return err
	}

	c.JSON(http.StatusCreated, toTodoResponse(todo))
	return nil
}

// GetTodo TODOを取得
func (h *handler) GetTodo(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	todo, err := h.todoService.GetTodo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return err
	}

	c.JSON(http.StatusOK, toTodoResponse(todo))
	return nil
}

// ListTodos TODO一覧を取得
func (h *handler) ListTodos(ctx interface{}) error {
	c := ctx.(*gin.Context)

	todos, err := h.todoService.ListTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to list todos: %v", err)})
		return err
	}

	responses := make([]map[string]interface{}, len(todos))
	for i, todo := range todos {
		responses[i] = toTodoResponse(todo)
	}

	c.JSON(http.StatusOK, gin.H{"todos": responses})
	return nil
}

// GetTodosByDateRange 日付範囲でTODOを取得
func (h *handler) GetTodosByDateRange(ctx interface{}) error {
	c := ctx.(*gin.Context)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return fmt.Errorf("start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return err
	}

	todos, err := h.todoService.GetTodosByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get todos: %v", err)})
		return err
	}

	responses := make([]map[string]interface{}, len(todos))
	for i, todo := range todos {
		responses[i] = toTodoResponse(todo)
	}

	c.JSON(http.StatusOK, gin.H{"todos": responses})
	return nil
}

// GetTodosByDate 特定の日付のTODOを取得
func (h *handler) GetTodosByDate(ctx interface{}) error {
	c := ctx.(*gin.Context)

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return err
	}

	todos, err := h.todoService.GetTodosByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get todos: %v", err)})
		return err
	}

	responses := make([]map[string]interface{}, len(todos))
	for i, todo := range todos {
		responses[i] = toTodoResponse(todo)
	}

	c.JSON(http.StatusOK, gin.H{"todos": responses})
	return nil
}

// GetTodosWithoutDueDate 期限未設定のTODOを取得（ページネーション付き）
func (h *handler) GetTodosWithoutDueDate(ctx interface{}) error {
	c := ctx.(*gin.Context)

	offset := 0
	limit := 20
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	todos, totalCount, err := h.todoService.GetTodosWithoutDueDate(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get todos: %v", err)})
		return err
	}

	responses := make([]map[string]interface{}, len(todos))
	for i, todo := range todos {
		responses[i] = toTodoResponse(todo)
	}

	hasMore := offset+limit < totalCount
	c.JSON(http.StatusOK, gin.H{
		"todos":    responses,
		"total":     totalCount,
		"offset":    offset,
		"limit":     limit,
		"has_more":  hasMore,
	})
	return nil
}

// UpdateTodo TODOを更新
func (h *handler) UpdateTodo(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	var req port.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	startDate, err := parseTime(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return err
	}
	endDate, err := parseTime(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return err
	}
	dueDate, err := parseTime(req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format"})
		return err
	}

	todo, err := h.todoService.UpdateTodo(id, req.Title, req.Description, startDate, endDate, dueDate, req.Completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update todo: %v", err)})
		return err
	}

	c.JSON(http.StatusOK, toTodoResponse(todo))
	return nil
}

// DeleteTodo TODOを削除
func (h *handler) DeleteTodo(ctx interface{}) error {
	c := ctx.(*gin.Context)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return err
	}

	if err := h.todoService.DeleteTodo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete todo: %v", err)})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo deleted successfully"})
	return nil
}

func toTodoResponse(todo *domain.Todo) map[string]interface{} {
	return map[string]interface{}{
		"id":          todo.ID.String(),
		"title":       todo.Title,
		"description": todo.Description,
		"start_date":  formatTime(todo.StartDate),
		"end_date":    formatTime(todo.EndDate),
		"due_date":    formatTime(todo.DueDate),
		"completed":   todo.Completed,
		"created_at":  todo.CreatedAt.Format(time.RFC3339),
		"updated_at":  todo.UpdatedAt.Format(time.RFC3339),
	}
}
