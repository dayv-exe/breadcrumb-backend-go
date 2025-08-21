package models

// standard user db model

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type User struct {
	Userid        string                          `dynamodbav:"pk" json:"userId"`
	Nickname      string                          `dynamodbav:"nickname" json:"nickname"`
	Name          string                          `dynamodbav:"name" json:"name"`
	Bio           string                          `dynamodbav:"bio" json:"bio"`
	DpUrl         string                          `dynamodbav:"dpUrl" json:"dpUrl"`
	IsSuspended   bool                            `dynamodbav:"is_suspended" json:"isSuspended"`
	IsDeactivated bool                            `dynamodbav:"is_deactivated" json:"isDeactivated"`
	UserLogs      map[string]types.AttributeValue `dynamodbav:"user_logs" json:"userLogs"`
}

type CognitoInfo struct {
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
}

func NewUser(userid string, nickname string, name string, isSuspended bool) *User {
	return &User{
		Userid:        userid,
		Nickname:      nickname,
		Name:          name,
		Bio:           "",
		IsSuspended:   isSuspended,
		IsDeactivated: false,
		UserLogs:      NewUserLogs().DatabaseFormat(),
	}
}

func (u *User) DatabaseFormat() map[string]types.AttributeValue {
	item, err := attributevalue.MarshalMap(u)

	if err != nil {
		return nil
	}

	item["pk"] = &types.AttributeValueMemberS{Value: "USER#" + u.Userid}
	item["sk"] = &types.AttributeValueMemberS{Value: "PROFILE"}

	return item
}

func GetCognitoInfo(sub string, cognitoClient *cognitoidentityprovider.Client, ctx *context.Context, userPoolId string) (*CognitoInfo, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(userPoolId),
		Username:   aws.String(sub),
	}

	user, err := getCognitoInfoQueryRunner(func() (*cognitoidentityprovider.AdminGetUserOutput, error) {
		return cognitoClient.AdminGetUser(*ctx, input)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func getCognitoInfoQueryRunner(queryFn func() (*cognitoidentityprovider.AdminGetUserOutput, error)) (*CognitoInfo, error) {
	// THIS CUSTOM FUNCTION TO BE TESTED
	output, err := queryFn()

	if err != nil {
		return nil, err
	}

	userInfo := map[string]string{}

	for _, attr := range output.UserAttributes {
		userInfo[*attr.Name] = *attr.Value
	}

	return &CognitoInfo{
		Email:     userInfo["email"],
		Birthdate: userInfo["birthdate"],
	}, nil
}

func convertToUser(item map[string]types.AttributeValue) (*User, error) {
	var u User
	if err := attributevalue.UnmarshalMap(item, &u); err != nil {
		return nil, err
	}

	u.Userid = strings.TrimPrefix(u.Userid, "USER#")

	return &u, nil
}

func FindByNickname(nickname string, tableName string, ctx *context.Context, ddbClient *dynamodb.Client) (*User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("NicknameIndex"),
		KeyConditionExpression: aws.String("nickname = :nick"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nick": &types.AttributeValueMemberS{Value: nickname},
		},
		Limit: aws.Int32(1),
	}

	output, err := ddbClient.Query(*ctx, input)

	if err != nil {
		return nil, err
	}

	if output.Count < 1 {
		return nil, nil
	}

	return convertToUser(output.Items[0])
}
