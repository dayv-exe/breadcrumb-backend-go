package auth

import (
	"breadcrumb-backend-go/models"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FuncDependencies struct {
	DdbClient     *dynamodb.Client
	TableName     string
	CognitoClient *cognitoidentityprovider.Client
	UserPoolId    string
}

func (fd *FuncDependencies) HandleNicknameAvailable(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// if nickname is invalid return false immediately
	nickname := strings.ToLower(req.PathParameters["nickname"])
	if nickname == "" || !isNicknameValid(nickname) {
		return models.SuccessfulRequestResponse(fmt.Sprintf("%v", false), false), nil
	}

	takenInDb, dbErr := isNicknameTakenInDynamodb(func() (*dynamodb.QueryOutput, error) {
		return fd.DdbClient.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(fd.TableName),
			IndexName:              aws.String("NicknameIndex"),
			KeyConditionExpression: aws.String("nickname = :nick"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":nick": &types.AttributeValueMemberS{Value: nickname},
			},
			Limit: aws.Int32(1),
		})
	})

	cognitoFilter := fmt.Sprintf("nickname = \"%s\"", nickname)

	takenInCognito, cognitoErr := isNickNameTakenInCognito(func() (*cognitoidentityprovider.ListUsersOutput, error) {
		return fd.CognitoClient.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(fd.UserPoolId),
			Filter:     aws.String(cognitoFilter),
			Limit:      aws.Int32(1),
		})
	})

	if dbErr != nil {
		return models.ServerSideErrorResponse("An error has occurred, try again.", dbErr.Error()), nil
	}

	if cognitoErr != nil {
		return models.ServerSideErrorResponse("An error has occurred, try again.", cognitoErr.Error()), nil
	}

	available := !takenInCognito && !takenInDb

	return models.SuccessfulRequestResponse(fmt.Sprintf("%v", available), false), nil
}

func isNickNameTakenInCognito(queryFn func() (*cognitoidentityprovider.ListUsersOutput, error)) (bool, error) {
	out, err := queryFn()
	if err != nil {
		return false, err
	}

	return len(out.Users) > 0, nil
}

func isNicknameTakenInDynamodb(queryFn func() (*dynamodb.QueryOutput, error)) (bool, error) {
	out, err := queryFn()

	if err != nil {
		return false, err
	}

	isTaken := len(out.Items) > 0

	return isTaken, nil
}

func isNicknameValid(nickname string) bool {
	if len(nickname) < 3 || len(nickname) > 20 {
		return false
	}

	match, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9._]*[a-zA-Z0-9]$`, nickname)
	if !match {
		return false
	}

	if strings.Contains(nickname, "..") || strings.Contains(nickname, "__") {
		return false
	}

	if strings.Contains(nickname, "_.") || strings.Contains(nickname, "._") {
		return false
	}

	return true
}
