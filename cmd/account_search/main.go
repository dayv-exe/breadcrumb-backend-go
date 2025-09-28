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
	dbClient        *dynamodb.Client
	searchTableName string
	starter         discover.AccountSearchDependencies
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// load dynamodb stuff
	dbClient = dynamodb.NewFromConfig(cfg)
	searchTableName = os.Getenv("SEARCH_TABLE")
	if searchTableName == "" {
		panic("USERS_TABLE environment variable not set")
	}

	starter = discover.AccountSearchDependencies{
		Client:          dbClient,
		SearchTableName: searchTableName,
	}
}

func main() {
	lambda.Start(starter.HandleAccountSearch)
}
