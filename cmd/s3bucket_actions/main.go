package main

import (
	"breadcrumb-backend-go/handlers/media"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client   *s3.Client
	bucketName string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	s3Client = s3.NewFromConfig(cfg)
	bucketName = os.Getenv("MEDIA_BUCKET")
	if bucketName == "" {
		panic("MEDIA_BUCKET environment variable not set")
	}
}

func main() {
	handler := media.S3BucketActionsDependencies{
		S3Client:   s3Client,
		BucketName: bucketName,
	}

	lambda.Start(handler.HandleStorageActions)
}
