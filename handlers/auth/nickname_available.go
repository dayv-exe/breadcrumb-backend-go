package auth

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FuncDependencies struct {
	DdbClient *dynamodb.Client
	TableName string
}

func (fd *FuncDependencies) HandleNicknameAvailable(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// if nickname is invalid return false immediately
	nickname := strings.ToLower(req.PathParameters["nickname"])
	if nickname == "" || !utils.IsNicknameValid(nickname) {
		return models.SuccessfulRequestResponse(fmt.Sprintf("%v", false), false), nil
	}

	log.Print(nickname)

	takenInDb, dbErr := isNicknameTakenInDynamodb(func() (*dynamodb.QueryOutput, error) {
		return fd.DdbClient.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(fd.TableName),
			IndexName:              aws.String("NicknameIndex"),
			KeyConditionExpression: aws.String("nickname = :nick"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":nick": &types.AttributeValueMemberS{Value: nickname},
			},
			// Limit: aws.Int32(1),
		})
	})

	if dbErr != nil {
		return models.ServerSideErrorResponse("An error has occurred, try again.", dbErr.Error()), nil
	}

	available := !takenInDb

	return models.SuccessfulRequestResponse(fmt.Sprintf("%v", available), false), nil
}

func isNicknameTakenInDynamodb(queryFn func() (*dynamodb.QueryOutput, error)) (bool, error) {
	log.Print("nickname here")
	out, err := queryFn()

	if err != nil {
		return false, err
	}

	isTaken := len(out.Items) > 0
	for _, name := range out.Items {
		log.Print(name)
	}

	return isTaken, nil
}
