package models

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type nicknameType struct {
	PK     string `dynamodbav:"pk" json:"nickname"`
	Sk     string `dynamodbav:"sk"`
	UserId string `dynamodbav:"userId" json:"userId"`
}

func GetNicknameDbItem(userId string, nickname string) map[string]types.AttributeValue {
	n := nicknameType{
		PK:     "NICKNAME#" + strings.ToLower(nickname),
		Sk:     "NICKNAME",
		UserId: userId,
	}

	item, err := attributevalue.MarshalMap(n)

	if err != nil {
		return nil
	}

	return item
}
