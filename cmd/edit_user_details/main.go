package main

import (
	"breadcrumb-backend-go/handlers/account"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dbClient   *dynamodb.Client
	tableNames *utils.TableNames
	starter    account.EditUserDetailsDependency
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	dbClient = dynamodb.NewFromConfig(cfg)
	tableNames = utils.GetAllTableNames()

	starter = account.EditUserDetailsDependency{
		DdbClient:  dbClient,
		TableNames: tableNames,
	}
}

func main() {
	lambda.Start(starter.HandleEditUserDetails)
}
