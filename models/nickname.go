package models

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Nickname struct {
	Nickname      string `dynamodbav:"pk" json:"nickname"`
	dbDescription string `dynamodbav:"sk"`
	UserId        string `dynamodbav:"userId" json:"userId"`
	Fullname      string `dynamodbav:"name" json:"fullname"`
}

func NewNickname(nicknameStr string, name string, userid string) *Nickname {
	return &Nickname{
		Nickname:      "NICKNAME#" + nicknameStr,
		dbDescription: "NICKNAME",
		UserId:        userid,
		Fullname:      name,
	}
}

func NicknameKey(nickname string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "NICKNAME#" + nickname},
		"sk": &types.AttributeValueMemberS{Value: "NICKNAME"},
	}
}

func (n *Nickname) DatabaseFormat() *map[string]types.AttributeValue {
	item, err := attributevalue.MarshalMap(n)

	if err != nil {
		return nil
	}

	return &item
}
