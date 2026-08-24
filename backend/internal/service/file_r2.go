package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *FileService) uploadR2(objectKey string, data []byte) error {
	endpoint := s.config.GetString("data.storage.r2.endpoint")
	bucket := s.config.GetString("data.storage.r2.bucket")
	region := s.config.GetString("data.storage.r2.region")
	accessKeyID := s.config.GetString("data.storage.r2.access_key_id")
	secretAccessKey := s.config.GetString("data.storage.r2.secret_access_key")

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("put R2 object: %w", err)
	}
	return nil
}
