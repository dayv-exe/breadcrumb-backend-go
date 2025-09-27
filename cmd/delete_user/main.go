package main

import (
	"breadcrumb-backend-go/handlers/account"
	"breadcrumb-backend-go/utils"
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
	tableNames    *utils.TableNames
	userPoolId    string
	cognitoClient *cognitoidentityprovider.Client
	starter       account.DeleteUserDependencies
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
	tableNames = utils.GetAllTableNames()

	starter = account.DeleteUserDependencies{
		DbClient:      dbClient,
		TableNames:    tableNames,
		CognitoClient: cognitoClient,
		UserPoolId:    userPoolId,
	}
}

func main() {
	lambda.Start(starter.HandleDeleteUser)
}
