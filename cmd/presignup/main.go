package main

import (
	"breadcrumb-backend-go/handlers/auth"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var ddbClient *dynamodb.Client
var usersTable string

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	ddbClient = dynamodb.NewFromConfig(cfg)
	usersTable = os.Getenv("USERS_TABLE")
	if usersTable == "" {
		panic("USERS_TABLE environment variable not set")
	}
}

func main() {
	fd := auth.PreSignupDependencies{
		DdbClient:     ddbClient,
		UserTableName: usersTable,
	}

	lambda.Start(fd.PreSignupHandler)
}
