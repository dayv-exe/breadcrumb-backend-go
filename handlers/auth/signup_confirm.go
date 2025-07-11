package auth

import (
	"breadcrumb-backend-go/models"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	userID := event.Request.UserAttributes["sub"]
	email := event.Request.UserAttributes["email"]
	nickName := event.Request.UserAttributes["nickname"]
	name := event.Request.UserAttributes["name"]
	birthdate := event.Request.UserAttributes["birthdate"]

	newUser, uErr := models.NewUser(userID, email, nickName, name, birthdate).DatabaseFormat()
	if uErr != nil {
		log.Fatalf("AN ERROR OCCURRED WHILE ADDING NEW USER! %v", uErr)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(usersTable),
		Item:      newUser,
	}

	_, err := db.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to write user to DynamoDB: %w", err)
	}

	return event, nil
}
