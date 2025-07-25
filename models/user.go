package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Userid      string
	Nickname    string
	Name        string
	Email       string
	Birthdate   string
	DpUrl       string
	isSuspended bool
	UserLogs    UserLogs
}

func NewUser(userid string, nickname string, name string, birthdate string, isSuspended bool) User {
	return User{
		Userid:      userid,
		Nickname:    nickname,
		Name:        name,
		Birthdate:   birthdate,
		isSuspended: isSuspended,
		UserLogs:    NewUserLogs(),
	}
}

func (u User) DatabaseFormat() (map[string]types.AttributeValue, error) {
	logsJSON, err := json.Marshal(u.UserLogs.DatabaseFormat())
	if err != nil {
		return nil, err
	}

	return map[string]types.AttributeValue{
		"pk":           &types.AttributeValueMemberS{Value: fmt.Sprintf("USERS#%s", u.Userid)},
		"sk":           &types.AttributeValueMemberS{Value: "PROFILE"},
		"name":         &types.AttributeValueMemberS{Value: u.Name},
		"nickname":     &types.AttributeValueMemberS{Value: strings.ToLower(u.Nickname)},
		"birthdate":    &types.AttributeValueMemberS{Value: u.Birthdate},
		"dpUrl":        &types.AttributeValueMemberS{Value: u.DpUrl},
		"is_suspended": &types.AttributeValueMemberS{Value: strings.ToLower(fmt.Sprint(u.isSuspended))},
		"userLogs":     &types.AttributeValueMemberS{Value: string(logsJSON)},
	}, nil
}
