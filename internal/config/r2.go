package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func (c *Config) LoadR2Config() (aws.Config, error) {
	cfgOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				c.R2AccessKeyID,
				c.R2SecretAccessKey,
				"",
			),
		),
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(), cfgOptions...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load R2 config: %w", err)
	}

	return awsConfig, nil
}
