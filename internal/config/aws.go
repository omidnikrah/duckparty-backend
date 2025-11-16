package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func (c *Config) LoadAWSConfig() (aws.Config, error) {
	var cfgOptions []func(*awsconfig.LoadOptions) error

	cfgOptions = append(cfgOptions, awsconfig.WithRegion(c.S3Region))

	if c.AWSAccessKeyID != "" && c.AWSSecretAccessKey != "" {
		cfgOptions = append(cfgOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				c.AWSAccessKeyID,
				c.AWSSecretAccessKey,
				"",
			),
		))
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(), cfgOptions...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return awsConfig, nil
}
