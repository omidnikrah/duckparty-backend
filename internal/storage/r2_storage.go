package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	appconfig "github.com/omidnikrah/duckparty-backend/internal/config"
)

type R2Storage struct {
	client *s3.Client
	cfg    *appconfig.Config
}

func NewR2Storage(appConfig *appconfig.Config) (*R2Storage, error) {
	r2Config, err := appConfig.LoadR2Config()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", appConfig.R2AccountID)

	client := s3.NewFromConfig(r2Config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &R2Storage{
		client: client,
		cfg:    appConfig,
	}, nil
}

func (s *R2Storage) UploadFile(fileContent []byte, name string) (string, error) {
	uniqueID := generateUniqueID()
	filename := fmt.Sprintf("duck_%s_%s.png", name, uniqueID)
	key := fmt.Sprintf("ducks/%s", filename)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.R2Bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(fileContent),
		ContentType:   aws.String("image/png"),
		ContentLength: aws.Int64(int64(len(fileContent))),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to R2: %w", err)
	}

	imageURL := fmt.Sprintf("%s/%s", s.cfg.R2BaseURL, key)
	return imageURL, nil
}

func generateUniqueID() string {
	id := uuid.New()
	return id.String()[:8]
}
