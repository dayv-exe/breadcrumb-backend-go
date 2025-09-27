package main

import (
	"breadcrumb-backend-go/handlers/auth"
	"breadcrumb-backend-go/utils"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	db            *dynamodb.Client
	cognitoClient *cognitoidentityprovider.Client
	tableNames    *utils.TableNames
	starter       auth.PostConfirmationDependencies
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}
	db = dynamodb.NewFromConfig(cfg)
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)

	tableNames = utils.GetAllTableNames()
	starter = auth.PostConfirmationDependencies{
		DdbClient:     db,
		TableNames:    tableNames,
		CognitoClient: cognitoClient,
	}
}

func main() {
	lambda.Start(starter.HandlePostConfirmation)
}
