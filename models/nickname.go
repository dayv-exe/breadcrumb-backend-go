package models

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type nicknameType struct {
	PK     string `dynamodbav:"pk" json:"nickname"`
	Sk     string `dynamodbav:"sk"`
	UserId string `dynamodbav:"userId" json:"userId"`
}

func GetNicknameDbItem(user *User) map[string]types.AttributeValue {
	n := nicknameType{
		PK:     user.Nickname,
		Sk:     "NICKNAME",
		UserId: user.Userid,
	}

	item, err := attributevalue.MarshalMap(n)

	if err != nil {
		return nil
	}

	return item
}
