package main

import (
	apis "breadcrumb-backend-go/handlers/account"
	"breadcrumb-backend-go/utils"
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
	tableNames    *utils.TableNames
	starter       apis.GetUserDetailsDependencies
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	ddbClient = dynamodb.NewFromConfig(cfg)

	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	userPoolId = os.Getenv("USER_POOL_ID")
	if userPoolId == "" {
		log.Fatal("USER_POOL_ID environment variable is required")
	}

	tableNames = utils.GetAllTableNames()

	starter = apis.GetUserDetailsDependencies{
		UserPoolId:    userPoolId,
		CognitoClient: cognitoClient,
		DdbClient:     ddbClient,
		TableNames:    tableNames,
	}
}

func main() {
	lambda.Start(starter.HandleGetUserDetails)
}
