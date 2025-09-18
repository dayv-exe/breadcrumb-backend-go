package models

import (
	"breadcrumb-backend-go/utils"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type friendRequest struct {
	RecipientId     string `dynamodbav:"pk" json:"recipient"`
	SenderId        string `dynamodbav:"sk" json:"sender"`
	Date            string `dynamodbav:"date" json:"date"`
	SendersName     string `dynamodbav:"name" json:"name"`
	SendersNickname string `dynamodbav:"nickname" json:"nickname"`
	SendersDpUrl    string `dynamodbav:"dp_url" json:"dpUrl"`
	SendersFgCol    string `dynamodbav:"default_pic_fg" json:"defaultPicFg"`
	SendersBgCol    string `dynamodbav:"default_pic_bg" json:"defaultPicBg"`
}

const (
	FriendRequestPkPrefix = "USER#"
	FriendRequestSkPrefix = "FRIEND_REQUEST_FROM#"
)

func NewFriendRequest(recipientUserId string, sender *User) *friendRequest {
	fr := friendRequest{
		RecipientId:     recipientUserId,
		SenderId:        sender.Userid,
		Date:            utils.GetTimeNow(),
		SendersName:     sender.Name,
		SendersNickname: sender.Nickname,
		SendersDpUrl:    sender.DpUrl,
		SendersFgCol:    sender.DefaultProfilePicFgColor,
		SendersBgCol:    sender.DefaultProfilePicBgColor,
	}

	return &fr
}

func FriendRequestKey(recipientUserId string, senderUserId string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(FriendRequestPkPrefix, recipientUserId)},
		"sk": &types.AttributeValueMemberS{Value: utils.AddPrefix(FriendRequestSkPrefix, senderUserId)},
	}
}

func (fr friendRequest) DatabaseFormat() (*map[string]types.AttributeValue, error) {
	fr.RecipientId = utils.AddPrefix(FriendRequestPkPrefix, fr.RecipientId)
	fr.SenderId = utils.AddPrefix(FriendRequestSkPrefix, fr.SenderId)
	item, err := attributevalue.MarshalMap(fr)

	if err != nil {
		log.Print("an error occurred while trying to marshal friend req item")
		return nil, err
	}

	return &item, nil
}

func ConvertToFriendRequest(item *map[string]types.AttributeValue) (*friendRequest, error) {
	var fr friendRequest
	if err := attributevalue.UnmarshalMap(*item, &fr); err != nil {
		return nil, err
	}

	fr.RecipientId = strings.TrimPrefix(fr.RecipientId, FriendRequestPkPrefix)
	fr.SenderId = strings.TrimPrefix(fr.SenderId, FriendRequestSkPrefix)

	return &fr, nil
}
