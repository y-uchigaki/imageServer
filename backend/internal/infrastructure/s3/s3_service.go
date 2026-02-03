package s3

import (
	"bytes"
	"fmt"
	"imageServer/internal/port"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Service struct {
	s3Client     *s3.S3
	bucketName   string
	cloudFrontURL string
}

// NewS3Service S3サービスのコンストラクタ
func NewS3Service() (port.S3Service, error) {
	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_S3_BUCKET")
	cloudFrontURL := os.Getenv("AWS_CLOUDFRONT_URL")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if region == "" || bucketName == "" || cloudFrontURL == "" {
		return nil, fmt.Errorf("AWS configuration is missing")
	}

	config := &aws.Config{
		Region: aws.String(region),
	}

	// LocalStackエンドポイントの設定（環境変数で指定可能）
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint != "" {
		config.Endpoint = aws.String(endpoint)
		config.S3ForcePathStyle = aws.Bool(true) // LocalStackはパススタイルを要求
	}

	// 認証情報が提供されている場合
	if accessKeyID != "" && secretAccessKey != "" {
		config.Credentials = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &s3Service{
		s3Client:      s3.New(sess),
		bucketName:    bucketName,
		cloudFrontURL: cloudFrontURL,
	}, nil
}

func (s *s3Service) UploadImage(key string, data []byte, contentType string) error {
	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	return err
}

func (s *s3Service) GetCloudFrontURL(key string) string {
	return fmt.Sprintf("%s/%s", s.cloudFrontURL, key)
}

func (s *s3Service) DeleteImage(key string) error {
	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}
