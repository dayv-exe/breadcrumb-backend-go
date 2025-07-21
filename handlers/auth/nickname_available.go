package auth

import (
	"breadcrumb-backend-go/models"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ddbClient *dynamodb.Client
var tableName string

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	ddbClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("USERS_TABLE")
	if tableName == "" {
		panic("USERS_TABLE environment variable not set")
	}
}

func HandleNicknameAvailable(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	nickname := req.PathParameters["nickname"]
	if nickname == "" {
		return models.InvalidRequestErrorResponse(""), nil
	}

	out, err := ddbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("NicknameIndex"),
		KeyConditionExpression: aws.String("nickname = :nick"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nick": &types.AttributeValueMemberS{Value: nickname},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return models.ServerSideErrorResponse("An error has occurred, try again."), nil
	}

	exists := len(out.Items) < 1

	return models.SuccessfulRequestResponse(fmt.Sprintf("%v", exists), false), nil
}
