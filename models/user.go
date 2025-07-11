package models

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Userid    string
	Nickname  string
	Name      string
	Email     string
	Birthdate string
	DpUrl     string
	UserLogs  UserLogs
}

func NewUser(userid string, nickname string, name string, email string, birthdate string) User {
	return User{
		Userid:    userid,
		Nickname:  nickname,
		Name:      name,
		Email:     email,
		Birthdate: birthdate,
		UserLogs:  NewUserLogs(),
	}
}

func (u User) DatabaseFormat() (map[string]types.AttributeValue, error) {
	logsJSON, err := json.Marshal(u.UserLogs.DatabaseFormat())
	if err != nil {
		return nil, err
	}

	return map[string]types.AttributeValue{
		"pk":        &types.AttributeValueMemberS{Value: fmt.Sprintf("USERS#%s", u.Userid)},
		"sk":        &types.AttributeValueMemberS{Value: "PROFILE"},
		"name":      &types.AttributeValueMemberS{Value: u.Name},
		"nickname":  &types.AttributeValueMemberS{Value: u.Nickname},
		"birthdate": &types.AttributeValueMemberS{Value: u.Birthdate},
		"email":     &types.AttributeValueMemberS{Value: u.Email},
		"dpUrl":     &types.AttributeValueMemberS{Value: u.DpUrl},
		"userLogs":  &types.AttributeValueMemberS{Value: string(logsJSON)},
	}, nil
}
