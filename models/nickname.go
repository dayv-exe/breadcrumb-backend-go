package models

import (
	"breadcrumb-backend-go/utils"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Nickname struct {
	Nickname    string `dynamodbav:"pk" json:"nickname"`
	Description string `dynamodbav:"sk"`
	UserId      string `dynamodbav:"userId" json:"userId"`
}

const (
	nicknamePkPrefix = "NICKNAME#"
	nicknameSkPrefix = "NICKNAME"
)

// david kalu arubuike
// 1.) david = pk: fullname#da sk: fullname#david#123 normalized: david kalu arubuike raw: david kalu arubuike
// 2.) david = pk: fullname#ar sk: fullname#arubuike raw: david kalu arubuike
// 3.) david = pk: fullname#ka sk: fullname#kalu raw: david kalu arubuike

func NewNickname(nicknameStr string, name string, userid string) *Nickname {
	return &Nickname{
		Nickname:    nicknameStr,
		Description: nicknameSkPrefix,
		UserId:      userid,
	}
}

func NicknameKey(nickname string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(nicknamePkPrefix, nickname)},
		"sk": &types.AttributeValueMemberS{Value: nicknameSkPrefix},
	}
}

func (n Nickname) DatabaseFormat() *map[string]types.AttributeValue {
	n.Nickname = utils.AddPrefix(nicknamePkPrefix, n.Nickname)
	item, err := attributevalue.MarshalMap(n)

	if err != nil {
		return nil
	}

	return &item
}

func ConvertToNickname(item *map[string]types.AttributeValue) (*Nickname, error) {
	var n Nickname
	if err := attributevalue.UnmarshalMap(*item, &n); err != nil {
		return nil, err
	}

	// clean up db tags
	n.Nickname = strings.TrimPrefix(n.Nickname, nicknamePkPrefix)
	return &n, nil
}
