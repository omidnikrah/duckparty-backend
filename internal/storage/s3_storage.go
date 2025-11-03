package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	appconfig "github.com/omidnikrah/duckparty-backend/internal/config"
)

type S3Storage struct {
	client *s3.Client
	cfg    *appconfig.Config
}

func NewS3Storage(appConfig *appconfig.Config) (*S3Storage, error) {
	var cfgOptions []func(*awsconfig.LoadOptions) error

	cfgOptions = append(cfgOptions, awsconfig.WithRegion(appConfig.S3Region))

	if appConfig.AWSAccessKeyID != "" && appConfig.AWSSecretAccessKey != "" {
		cfgOptions = append(cfgOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				appConfig.AWSAccessKeyID,
				appConfig.AWSSecretAccessKey,
				"",
			),
		))
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(), cfgOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsConfig)

	return &S3Storage{
		client: client,
		cfg:    appConfig,
	}, nil
}

func (s *S3Storage) UploadFile(fileContent []byte, name string) (string, error) {
	uniqueID := generateUniqueID()
	filename := fmt.Sprintf("duck_%s_%s.png", name, uniqueID)
	key := fmt.Sprintf("ducks/%s", filename)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.S3Bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(fileContent),
		ContentType:   aws.String("image/png"),
		ContentLength: aws.Int64(int64(len(fileContent))),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	imageURL := fmt.Sprintf("%s/%s", s.cfg.S3BaseURL, key)
	return imageURL, nil
}

func generateUniqueID() string {
	id := uuid.New()
	return id.String()[:8]
}
