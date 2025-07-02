package store

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3BlobStore struct {
	client *s3.Client
}

func NewS3BlobStore() (*S3BlobStore, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("load default config: %w", err)
	}

	return &S3BlobStore{
		client: s3.NewFromConfig(cfg),
	}, nil
}

func (s *S3BlobStore) GetFileContent(ctx context.Context, bucket, fileKey string) (string, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return "", fmt.Errorf("get object from S3: %w", err)
	}
	defer out.Body.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, out.Body)
	if err != nil {
		return "", fmt.Errorf("read output body: %w", err)
	}

	return buf.String(), nil
}
