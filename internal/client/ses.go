package client

import (
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/omidnikrah/duckparty-backend/internal/config"
)

func NewSESClient(appConfig *config.Config) (*ses.Client, error) {
	awsConfig, err := appConfig.LoadAWSConfig()
	if err != nil {
		return nil, err
	}

	return ses.NewFromConfig(awsConfig), nil
}
