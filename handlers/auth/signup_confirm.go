package auth

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	db         *dynamodb.Client
	usersTable = os.Getenv("USERS_TABLE")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}
	db = dynamodb.NewFromConfig(cfg)
}

func HandlePostConfirmation(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (interface{}, error) {
	// runs after user has validated their email

	userID := event.Request.UserAttributes["sub"]
	nickName := event.Request.UserAttributes["nickname"]
	name := event.Request.UserAttributes["name"]

	// create new user
	newUser, uErr := models.NewUser(userID, nickName, name, false).DatabaseFormat()
	if uErr != nil {
		log.Fatalf("AN ERROR OCCURRED WHILE ADDING NEW USER! %v", uErr)
	}

	// add to db
	_, err := utils.PutItemInDDbTable(newUser, usersTable, db, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to write user to DynamoDB: %w", err)
	}
	return event, nil
}
