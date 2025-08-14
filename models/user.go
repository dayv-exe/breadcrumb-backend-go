package models

// standard user db model

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Userid        string
	Nickname      string
	Name          string
	Bio           string
	DpUrl         string
	isSuspended   bool
	isDeactivated bool
	UserLogs      UserLogs
}

func NewUser(userid string, nickname string, name string, isSuspended bool) User {
	return User{
		Userid:        userid,
		Nickname:      nickname,
		Name:          name,
		Bio:           "",
		isSuspended:   isSuspended,
		isDeactivated: false,
		UserLogs:      NewUserLogs(),
	}
}

func (u User) DatabaseFormat() map[string]types.AttributeValue {
	logsJSON := u.UserLogs.DatabaseFormat()

	return map[string]types.AttributeValue{
		"pk":             &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", u.Userid)},
		"sk":             &types.AttributeValueMemberS{Value: "PROFILE"},
		"name":           &types.AttributeValueMemberS{Value: u.Name},
		"nickname":       &types.AttributeValueMemberS{Value: strings.ToLower(u.Nickname)},
		"bio":            &types.AttributeValueMemberS{Value: ""},
		"dpUrl":          &types.AttributeValueMemberS{Value: u.DpUrl},
		"is_suspended":   &types.AttributeValueMemberBOOL{Value: u.isSuspended},
		"is_deactivated": &types.AttributeValueMemberBOOL{Value: u.isDeactivated},
		"user_logs":      &types.AttributeValueMemberM{Value: logsJSON},
	}
}
