package main

import (
	"breadcrumb-backend-go/handlers/auth"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var ddbClient *dynamodb.Client
var tableName string
var cognitoClient *cognitoidentityprovider.Client
var userPoolId string

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
		panic("USER_POOL_ID environment variable not set")
	}
}

func main() {
	fd := auth.FuncDependencies{
		DdbClient:     ddbClient,
		TableName:     tableName,
		CognitoClient: cognitoClient,
		UserPoolId:    userPoolId,
	}
	lambda.Start(fd.HandleNicknameAvailable)
}
