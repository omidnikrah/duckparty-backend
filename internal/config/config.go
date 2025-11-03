package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	S3Bucket   string
	S3Region   string
	S3BaseURL  string
	AWSAccessKeyID string
	AWSSecretAccessKey string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	config := &Config{
		AppPort:           os.Getenv("APP_PORT"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		S3Bucket:          os.Getenv("S3_BUCKET"),
		S3Region:          os.Getenv("S3_REGION"),
		S3BaseURL:         os.Getenv("S3_BASE_URL"),
		AWSAccessKeyID:    os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	return config, nil
}
