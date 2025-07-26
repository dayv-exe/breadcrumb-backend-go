package auth

import (
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PreSignupDependencies struct {
	DdbClient *dynamodb.Client
	TableName string
}

func (deps *PreSignupDependencies) PreSignupHandler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	nickname := event.Request.UserAttributes["nickname"]

	if !utils.IsNicknameValid(nickname) {
		return event, fmt.Errorf("invalid nickname")
	}

	nicknameAvail, err := utils.IsNicknameAvailableInDynamodb(nickname, deps.TableName, deps.DdbClient, ctx)

	if err != nil {
		return event, fmt.Errorf("error checking nickname availability %w", err)
	}

	if !nicknameAvail {
		return event, fmt.Errorf("nickname taken")
	}

	return event, nil
}
