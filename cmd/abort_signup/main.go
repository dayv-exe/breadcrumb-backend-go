package main

import (
	"breadcrumb-backend-go/handlers/auth"
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	client     *cognitoidentityprovider.Client
	userPoolId string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	client = cognitoidentityprovider.NewFromConfig(cfg)

	userPoolId = os.Getenv("COGNITO_USER_POOL_ID")
	if userPoolId == "" {
		log.Fatal("COGNITO_USER_POOL_ID environment variable is required")
	}

	asd := auth.AbortSignupDependencies{
		Client:     client,
		UserPoolId: userPoolId,
	}

	lambda.Start(asd.AbortSignupHandler)
}
