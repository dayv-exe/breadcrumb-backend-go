package auth

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PostConfirmationDependencies struct {
	DdbClient *dynamodb.Client
	TableName string
}

func (deps PostConfirmationDependencies) HandlePostConfirmation(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (interface{}, error) {
	// runs after user has validated their email

	userID := event.Request.UserAttributes["sub"]
	nickName := event.Request.UserAttributes["nickname"]
	name := event.Request.UserAttributes["name"]

	// create new user
	newUser := models.NewUser(userID, nickName, name, false).DatabaseFormat()

	database := utils.DynamoDbHelper{
		TableName: deps.TableName,
		Client:    deps.DdbClient,
		Ctx:       &ctx,
	}

	// add to db
	_, err := database.PutItem(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to write user to DynamoDB: %w", err)
	}
	return event, nil
}
