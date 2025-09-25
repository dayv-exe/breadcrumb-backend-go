package main

import (
	"breadcrumb-backend-go/handlers/discover"
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	userTable string
	dbClient  *dynamodb.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// load dynamodb stuff
	dbClient = dynamodb.NewFromConfig(cfg)
	userTable = os.Getenv("USERS_TABLE")
	if userTable == "" {
		panic("USERS_TABLE environment variable not set")
	}
}

func main() {
	starter := discover.FriendRequestDependencies{
		DbClient:      dbClient,
		UserTableName: userTable,
	}

	lambda.Start(starter.HandleFriendshipAction)
}
