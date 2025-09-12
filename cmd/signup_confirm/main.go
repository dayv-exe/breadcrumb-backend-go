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
	db            *dynamodb.Client
	usersTable    string
	searchTable   string
	cognitoClient *cognitoidentityprovider.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}
	db = dynamodb.NewFromConfig(cfg)
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)

	usersTable = os.Getenv("USERS_TABLE")
	if usersTable == "" {
		panic("USERS_TABLE environment variable not set")
	}

	searchTable = os.Getenv("SEARCH_TABLE")
	if searchTable == "" {
		panic("SEARCH_TABLE environment variable not set")
	}
}

func main() {

	hpc := auth.PostConfirmationDependencies{
		DdbClient:       db,
		TableName:       usersTable,
		SearchTableName: searchTable,
		CognitoClient:   cognitoClient,
	}

	lambda.Start(hpc.HandlePostConfirmation)
}
