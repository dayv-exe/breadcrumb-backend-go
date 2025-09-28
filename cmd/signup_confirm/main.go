package main

import (
	"breadcrumb-backend-go/handlers/auth"
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	db              *dynamodb.Client
	cognitoClient   *cognitoidentityprovider.Client
	tableName       string
	searchTableName string
	starter         auth.PostConfirmationDependencies
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}
	db = dynamodb.NewFromConfig(cfg)
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)

	tableName = os.Getenv("USERS_TABLE")
	if tableName == "" {
		panic("USERS_TABLE not set")
	}

	searchTableName = os.Getenv("SEARCH_TABLE")
	if searchTableName == "" {
		panic("SEARCH_TABLE environment not set")
	}

	starter = auth.PostConfirmationDependencies{
		DdbClient:     db,
		TableName:     tableName,
		CognitoClient: cognitoClient,
	}
}

func main() {
	lambda.Start(starter.HandlePostConfirmation)
}
