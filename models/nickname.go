package models

type nicknameType struct {
	Nickname string `dynamodbav:"pk" json:"nickname"`
	UserId   string `dynamodbav:"sk" json:"userId"`
}

func GetNicknameDbItem(u *User) nicknameType {
	return nicknameType{
		u.Nickname,
		u.Userid,
	}
}
