package main

import (
	"breadcrumb-backend-go/handlers/account"
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dbClient      *dynamodb.Client
	userTable     string
	searchTable   string
	userPoolId    string
	cognitoClient *cognitoidentityprovider.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// load cognito stuff
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)

	userPoolId = os.Getenv("USER_POOL_ID")
	if userPoolId == "" {
		log.Fatal("USER_POOL_ID environment variable is required")
	}

	// load dynamodb stuff
	dbClient = dynamodb.NewFromConfig(cfg)
	userTable = os.Getenv("USERS_TABLE")
	if userTable == "" {
		panic("USERS_TABLE environment variable not set")
	}

	searchTable = os.Getenv("SEARCH_TABLE")
	if searchTable == "" {
		panic("SEARCH_TABLE environment variable not set")
	}
}

func main() {
	starter := account.DeleteUserDependencies{
		DbClient:        dbClient,
		UserTableName:   userTable,
		SearchTableName: searchTable,
		CognitoClient:   cognitoClient,
		UserPoolId:      userPoolId,
	}

	lambda.Start(starter.HandleDeleteUser)
}
