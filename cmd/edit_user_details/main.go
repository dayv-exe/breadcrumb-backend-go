package main

import (
	"breadcrumb-backend-go/handlers/account"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dbClient    *dynamodb.Client
	userTable   string
	searchTable string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

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
	deps := account.EditUserDetailsDependency{
		DdbClient:   dbClient,
		UserTable:   userTable,
		SearchTable: searchTable,
	}

	lambda.Start(deps.HandleEditUserDetails)
}
