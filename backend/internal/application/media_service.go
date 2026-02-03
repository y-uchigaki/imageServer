package application

import (
	"fmt"
	"imageServer/internal/domain"
	"imageServer/internal/port"
	"time"

	"github.com/google/uuid"
)

// MediaService メディアサービスのユースケース
type MediaService struct {
	mediaRepo port.MediaRepository
	tagRepo   port.TagRepository
	s3Service port.S3Service
}

// NewMediaService メディアサービスのコンストラクタ
func NewMediaService(mediaRepo port.MediaRepository, tagRepo port.TagRepository, s3Service port.S3Service) *MediaService {
	return &MediaService{
		mediaRepo: mediaRepo,
		tagRepo:   tagRepo,
		s3Service: s3Service,
	}
}

// CreateImageMedia 画像メディアを作成
func (s *MediaService) CreateImageMedia(s3Key, title string, description *string, tagIDs []uuid.UUID) (*domain.Media, error) {
	now := time.Now()
	media := &domain.Media{
		ID:          uuid.New(),
		Type:        domain.MediaTypeImage,
		S3Key:       &s3Key,
		CloudFrontURL: stringPtr(s.s3Service.GetCloudFrontURL(s3Key)),
		Title:       title,
		Description: description,
		Tags:        []domain.Tag{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// タグを取得してメディアオブジェクトに追加
	// 実際のDBへの関連付けはCreateメソッド内で行われる
	for _, tagID := range tagIDs {
		tag, err := s.tagRepo.FindByID(tagID)
		if err != nil {
			return nil, fmt.Errorf("failed to find tag: %w", err)
		}
		media.Tags = append(media.Tags, *tag)
	}

	// メディアを作成（タグの関連付けもCreateメソッド内で行われる）
	if err := s.mediaRepo.Create(media); err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	return media, nil
}

// CreateYouTubeMedia YouTube動画メディアを作成
func (s *MediaService) CreateYouTubeMedia(youtubeURL, title string, description *string, tagIDs []uuid.UUID) (*domain.Media, error) {
	now := time.Now()
	media := &domain.Media{
		ID:          uuid.New(),
		Type:        domain.MediaTypeVideo,
		YouTubeURL:  &youtubeURL,
		Title:       title,
		Description: description,
		Tags:        []domain.Tag{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// タグを取得してメディアオブジェクトに追加
	// 実際のDBへの関連付けはCreateメソッド内で行われる
	for _, tagID := range tagIDs {
		tag, err := s.tagRepo.FindByID(tagID)
		if err != nil {
			return nil, fmt.Errorf("failed to find tag: %w", err)
		}
		media.Tags = append(media.Tags, *tag)
	}

	// メディアを作成（タグの関連付けもCreateメソッド内で行われる）
	if err := s.mediaRepo.Create(media); err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	return media, nil
}

// GetMedia メディアを取得
func (s *MediaService) GetMedia(id uuid.UUID) (*domain.Media, error) {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}

	// CloudFront経由のS3呼び出しだった場合、URLを更新
	if (media.IsImage() || media.IsAudio()) && media.S3Key != nil {
		media.CloudFrontURL = stringPtr(s.s3Service.GetCloudFrontURL(*media.S3Key))
	}

	return media, nil
}

// ListMedia メディア一覧を取得
func (s *MediaService) ListMedia() ([]*domain.Media, error) {
	mediaList, err := s.mediaRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}

	// CloudFront URLを更新
	for _, media := range mediaList {
		if (media.IsImage() || media.IsAudio()) && media.S3Key != nil {
			media.CloudFrontURL = stringPtr(s.s3Service.GetCloudFrontURL(*media.S3Key))
		}
	}

	return mediaList, nil
}

// ListMediaWithPagination ページネーション付きでメディア一覧を取得
func (s *MediaService) ListMediaWithPagination(offset, limit int) ([]*domain.Media, int, error) {
	mediaList, totalCount, err := s.mediaRepo.FindAllWithPagination(offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list media with pagination: %w", err)
	}

	// CloudFront URLを更新
	for _, media := range mediaList {
		if (media.IsImage() || media.IsAudio()) && media.S3Key != nil {
			media.CloudFrontURL = stringPtr(s.s3Service.GetCloudFrontURL(*media.S3Key))
		}
	}

	return mediaList, totalCount, nil
}

// ListMediaWithFilters フィルター付きでメディア一覧を取得
func (s *MediaService) ListMediaWithFilters(offset, limit int, titleSearch *string, tagIDs []uuid.UUID) ([]*domain.Media, int, error) {
	mediaList, totalCount, err := s.mediaRepo.FindAllWithFilters(offset, limit, titleSearch, tagIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list media with filters: %w", err)
	}

	// CloudFront URLを更新
	for _, media := range mediaList {
		if (media.IsImage() || media.IsAudio()) && media.S3Key != nil {
			media.CloudFrontURL = stringPtr(s.s3Service.GetCloudFrontURL(*media.S3Key))
		}
	}

	return mediaList, totalCount, nil
}

// GetMediaByTag タグでメディアを取得
func (s *MediaService) GetMediaByTag(tagID uuid.UUID) ([]*domain.Media, error) {
	mediaList, err := s.mediaRepo.FindByTagID(tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to find media by tag: %w", err)
	}

	// CloudFront URLを更新
	for _, media := range mediaList {
		if (media.IsImage() || media.IsAudio()) && media.S3Key != nil {
			media.CloudFrontURL = stringPtr(s.s3Service.GetCloudFrontURL(*media.S3Key))
		}
	}

	return mediaList, nil
}

// DeleteMedia メディアを削除
func (s *MediaService) DeleteMedia(id uuid.UUID) error {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find media: %w", err)
	}

	// S3からファイルを削除
	if (media.IsImage() || media.IsAudio()) && media.S3Key != nil {
		if err := s.s3Service.DeleteImage(*media.S3Key); err != nil {
			return fmt.Errorf("failed to delete file from S3: %w", err)
		}
	}

	if err := s.mediaRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}

	return nil
}

// AssociateTag メディアにタグを関連付け
func (s *MediaService) AssociateTag(mediaID, tagID uuid.UUID) error {
	// タグの存在確認
	_, err := s.tagRepo.FindByID(tagID)
	if err != nil {
		return fmt.Errorf("failed to find tag: %w", err)
	}

	if err := s.mediaRepo.AssociateTag(mediaID, tagID); err != nil {
		return fmt.Errorf("failed to associate tag: %w", err)
	}

	return nil
}

// RemoveTag メディアからタグを削除
func (s *MediaService) RemoveTag(mediaID, tagID uuid.UUID) error {
	if err := s.mediaRepo.RemoveTag(mediaID, tagID); err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}

// CreateAudioMedia 音声メディアを作成
func (s *MediaService) CreateAudioMedia(s3Key, title string, description *string, tagIDs []uuid.UUID) (*domain.Media, error) {
	now := time.Now()
	media := &domain.Media{
		ID:            uuid.New(),
		Type:          domain.MediaTypeAudio,
		S3Key:         &s3Key,
		CloudFrontURL: stringPtr(s.s3Service.GetCloudFrontURL(s3Key)),
		Title:         title,
		Description:   description,
		Tags:          []domain.Tag{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// タグを取得してメディアオブジェクトに追加
	// 実際のDBへの関連付けはCreateメソッド内で行われる
	for _, tagID := range tagIDs {
		tag, err := s.tagRepo.FindByID(tagID)
		if err != nil {
			return nil, fmt.Errorf("failed to find tag: %w", err)
		}
		media.Tags = append(media.Tags, *tag)
	}

	// メディアを作成（タグの関連付けもCreateメソッド内で行われる）
	if err := s.mediaRepo.Create(media); err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	return media, nil
}

// UploadImageToS3 S3に画像をアップロード
func (s *MediaService) UploadImageToS3(key string, data []byte, contentType string) error {
	return s.s3Service.UploadImage(key, data, contentType)
}

func stringPtr(s string) *string {
	return &s
}
