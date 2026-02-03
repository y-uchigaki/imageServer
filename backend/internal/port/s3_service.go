package port

// S3Service S3サービスのインターフェース
type S3Service interface {
	UploadImage(key string, data []byte, contentType string) error
	GetCloudFrontURL(key string) string
	DeleteImage(key string) error
}
