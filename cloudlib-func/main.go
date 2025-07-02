package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgdraganov/cloudlib/store"
)

type BlobStore interface {
	GetFileContent(ctx context.Context, bucket, fileKey string) (string, error)
}

type LambdaHandler struct {
	store      BlobStore
	bucketName string
}

func (l *LambdaHandler) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileKey := request.QueryStringParameters["file"]
	if fileKey == "" {
		log.Println("file query parameter is missing")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: "Missing file parameter",
		}, nil
	}

	content, err := l.store.GetFileContent(ctx, l.bucketName, fileKey)
	if err != nil {
		log.Printf("fetching file content from s3: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: content,
	}, nil
}

func main() {
	s3store, err := store.NewS3BlobStore()
	if err != nil {
		log.Fatalf("failed to create S3 blob store: %v", err)
	}

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("BUCKET_NAME environment variable is not set")
	}

	lh := &LambdaHandler{
		store:      s3store,
		bucketName: bucketName,
	}

	lambda.Start(lh.handler)
}
