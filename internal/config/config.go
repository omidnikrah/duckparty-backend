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
	R2AccountID        string
	R2Bucket           string
	R2BaseURL          string
	R2AccessKeyID      string
	R2SecretAccessKey  string
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
		R2AccountID:        os.Getenv("R2_ACCOUNT_ID"),
		R2Bucket:           os.Getenv("R2_BUCKET"),
		R2BaseURL:          os.Getenv("R2_BASE_URL"),
		R2AccessKeyID:      os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey:  os.Getenv("R2_SECRET_ACCESS_KEY"),
		RedisHost:          os.Getenv("REDIS_HOST"),
		RedisPort:          os.Getenv("REDIS_PORT"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AuthSenderEmail:    os.Getenv("AUTH_SENDER_EMAIL"),
		ResendAPIKey:       os.Getenv("RESEND_API_KEY"),
	}

	return config, nil
}
