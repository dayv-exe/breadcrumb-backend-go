package utils

import (
	"breadcrumb-backend-go/constants"
	"regexp"
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

func NicknameValid(nickname string) bool {
	if len(nickname) < constants.MIN_USERNAME_CHARS || len(nickname) > constants.MAX_USERNAME_CHARS {
		return false
	}

	// rules: must start and end with letter or number, only one of dot or underscore, must be 3 to 15 characters long
	match, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9._]*[a-zA-Z0-9]$`, nickname)
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
