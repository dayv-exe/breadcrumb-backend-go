package main

import (
	apis "breadcrumb-backend-go/handlers/account"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	ddbClient     *dynamodb.Client
	cognitoClient *cognitoidentityprovider.Client
	userPoolId    string
	tableName     string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	ddbClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("USERS_TABLE")
	if tableName == "" {
		panic("USERS_TABLE environment variable not set")
	}

	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	userPoolId = os.Getenv("USER_POOL_ID")
	if userPoolId == "" {
		log.Fatal("USER_POOL_ID environment variable is required")
	}
}

func main() {
	gud := apis.GetUserDetailsDependencies{
		UserPoolId:    userPoolId,
		CognitoClient: cognitoClient,
		TableName:     tableName,
		DdbClient:     ddbClient,
	}

	lambda.Start(gud.HandleGetUserDetails)
}
