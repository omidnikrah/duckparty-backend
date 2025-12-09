package client

import (
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/resend/resend-go/v3"
)

func NewResendClient(appConfig *config.Config) *resend.Client {
	return resend.NewClient(appConfig.ResendAPIKey)
}
