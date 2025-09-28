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
	dbClient  *dynamodb.Client
	tableName string
	starter   discover.FriendRequestDependencies
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// load dynamodb stuff
	dbClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("USERS_TABLE")
	if tableName == "" {
		panic("USERS_TABLE environment not set")
	}

	starter = discover.FriendRequestDependencies{
		DbClient:  dbClient,
		TableName: tableName,
	}
}

func main() {
	lambda.Start(starter.HandleFriendshipAction)
}
