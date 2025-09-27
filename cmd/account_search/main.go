package main

import (
	"breadcrumb-backend-go/handlers/discover"
	"breadcrumb-backend-go/utils"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dbClient   *dynamodb.Client
	tableNames *utils.TableNames
	starter    discover.AccountSearchDependencies
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// load dynamodb stuff
	dbClient = dynamodb.NewFromConfig(cfg)
	tableNames = utils.GetAllTableNames()

	starter = discover.AccountSearchDependencies{
		Client:     dbClient,
		TableNames: tableNames,
	}
}

func main() {
	lambda.Start(starter.HandleAccountSearch)
}
