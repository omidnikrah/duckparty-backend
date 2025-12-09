package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	S3Bucket           string
	S3Region           string
	S3BaseURL          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	RedisHost          string
	RedisPassword      string
	RedisPort          string
	JWTSecret          string
	AuthSenderEmail    string
	ResendAPIKey       string
	ApiPrefix          string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		AppPort:            os.Getenv("APP_PORT"),
		ApiPrefix:          os.Getenv("API_PREFIX"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		S3Bucket:           os.Getenv("S3_BUCKET"),
		S3Region:           os.Getenv("S3_REGION"),
		S3BaseURL:          os.Getenv("S3_BASE_URL"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		RedisHost:          os.Getenv("REDIS_HOST"),
		RedisPort:          os.Getenv("REDIS_PORT"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AuthSenderEmail:    os.Getenv("AUTH_SENDER_EMAIL"),
		ResendAPIKey:       os.Getenv("RESEND_API_KEY"),
	}

	return config, nil
}
