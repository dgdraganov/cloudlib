package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileKey := request.QueryStringParameters["file"]
	if fileKey == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing file parameter",
		}, nil
	}

	content, err := fetchFileContent(ctx, fileKey)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       content,
	}, nil
}

func fetchFileContent(ctx context.Context, fileKey string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return "", err
	}
	defer out.Body.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, out.Body)
	return buf.String(), err
}

func main() {
	lambda.Start(handler)
}
