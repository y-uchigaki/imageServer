package http

import (
	"imageServer/internal/port"

	"github.com/gin-gonic/gin"
)

// UploadImageHandler 画像をアップロード
// @Summary      画像をアップロード
// @Description  画像ファイルをS3にアップロードし、メディア情報をDBに保存します
// @Tags         media
// @Accept       multipart/form-data
// @Produce      json
// @Param        file        formData  file    true   "画像ファイル"
// @Param        title       formData  string  true   "タイトル"
// @Param        description formData  string  false  "説明"
// @Param        tag_ids     formData  array   false  "タグIDの配列"
// @Success      201         {object}  MediaResponse
// @Failure      400         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /media/upload [post]
func UploadImageHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.UploadImage(c)
	}
}

// CreateMediaWithYouTubeHandler YouTube URLでメディアを作成
// @Summary      YouTube URLでメディアを作成
// @Description  YouTube URLを指定してメディア情報をDBに保存します
// @Tags         media
// @Accept       json
// @Produce      json
// @Param        request  body      CreateMediaWithYouTubeRequest  true  "リクエスト"
// @Success      201      {object}  MediaResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /media/youtube [post]
func CreateMediaWithYouTubeHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.CreateMediaWithYouTube(c)
	}
}

// ListMediaHandler メディア一覧を取得
// @Summary      メディア一覧を取得
// @Description  すべてのメディアの一覧を取得します
// @Tags         media
// @Produce      json
// @Success      200  {object}  MediaListResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /media [get]
func ListMediaHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.ListMedia(c)
	}
}

// GetMediaHandler メディアを取得
// @Summary      メディアを取得
// @Description  IDを指定してメディア情報を取得します
// @Tags         media
// @Produce      json
// @Param        id   path      string  true  "メディアID"
// @Success      200  {object}  MediaResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /media/{id} [get]
func GetMediaHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.GetMedia(c)
	}
}

// DeleteMediaHandler メディアを削除
// @Summary      メディアを削除
// @Description  IDを指定してメディアを削除します（S3からも削除されます）
// @Tags         media
// @Produce      json
// @Param        id   path      string  true  "メディアID"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /media/{id} [delete]
func DeleteMediaHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.DeleteMedia(c)
	}
}

// CreateTagHandler タグを作成
// @Summary      タグを作成
// @Description  新しいタグを作成します
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        request  body      CreateTagRequest  true  "リクエスト"
// @Success      201      {object}  TagResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /tags [post]
func CreateTagHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.CreateTag(c)
	}
}

// ListTagsHandler タグ一覧を取得
// @Summary      タグ一覧を取得
// @Description  すべてのタグの一覧を取得します
// @Tags         tags
// @Produce      json
// @Success      200  {object}  TagListResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tags [get]
func ListTagsHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.ListTags(c)
	}
}

// GetTagHandler タグを取得
// @Summary      タグを取得
// @Description  IDを指定してタグ情報を取得します
// @Tags         tags
// @Produce      json
// @Param        id   path      string  true  "タグID"
// @Success      200  {object}  TagResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /tags/{id} [get]
func GetTagHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.GetTag(c)
	}
}

// UpdateTagHandler タグを更新
// @Summary      タグを更新
// @Description  IDを指定してタグ情報を更新します
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id       path      string           true  "タグID"
// @Param        request  body      UpdateTagRequest  true  "リクエスト"
// @Success      200      {object}  TagResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /tags/{id} [put]
func UpdateTagHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.UpdateTag(c)
	}
}

// DeleteTagHandler タグを削除
// @Summary      タグを削除
// @Description  IDを指定してタグを削除します
// @Tags         tags
// @Produce      json
// @Param        id   path      string  true  "タグID"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tags/{id} [delete]
func DeleteTagHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.DeleteTag(c)
	}
}

// AssociateMediaTagHandler メディアにタグを関連付け
// @Summary      メディアにタグを関連付け
// @Description  メディアにタグを関連付けます
// @Tags         media-tags
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "メディアID"
// @Param        request  body      AssociateTagRequest    true  "リクエスト"
// @Success      200      {object}  MessageResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /media/{id}/tags [post]
func AssociateMediaTagHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.AssociateMediaTag(c)
	}
}

// RemoveMediaTagHandler メディアからタグを削除
// @Summary      メディアからタグを削除
// @Description  メディアからタグの関連付けを削除します
// @Tags         media-tags
// @Produce      json
// @Param        id      path      string  true  "メディアID"
// @Param        tag_id  path      string  true  "タグID"
// @Success      200     {object}  MessageResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /media/{id}/tags/{tag_id} [delete]
func RemoveMediaTagHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.RemoveMediaTag(c)
	}
}

// GetMediaByTagHandler タグでメディアを取得
// @Summary      タグでメディアを取得
// @Description  タグIDを指定して関連するメディアの一覧を取得します
// @Tags         media-tags
// @Produce      json
// @Param        id   path      string  true  "タグID"
// @Success      200  {object}  MediaListResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tags/{id}/media [get]
func GetMediaByTagHandler(handler port.HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler.GetMediaByTag(c)
	}
}
