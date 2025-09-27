package main

import (
	"breadcrumb-backend-go/handlers/auth"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	ddbClient  *dynamodb.Client
	tableNames *utils.TableNames
	starter    auth.HandleNicknameAvailableDependencies
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	ddbClient = dynamodb.NewFromConfig(cfg)
	tableNames = utils.GetAllTableNames()

	starter = auth.HandleNicknameAvailableDependencies{
		DdbClient:  ddbClient,
		TableNames: tableNames,
	}
}

func main() {
	lambda.Start(starter.HandleNicknameAvailable)
}
