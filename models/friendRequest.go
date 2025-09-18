package models

import (
	"breadcrumb-backend-go/utils"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FriendRequest struct {
	RecipientId             string `dynamodbav:"pk"`
	SenderId                string `dynamodbav:"sk" json:"userId"`
	Date                    string `dynamodbav:"date" json:"date"`
	SendersName             string `dynamodbav:"name" json:"name"`
	SendersNickname         string `dynamodbav:"nickname" json:"nickname"`
	SendersDpUrl            string `dynamodbav:"dp_url" json:"dpUrl"`
	SendersDefaultPicColors string `dynamodbav:"default_pic_colors" json:"defaultPicColors"`
}

const (
	FriendRequestPkPrefix = "USER#"
	FriendRequestSkPrefix = "FRIEND_REQUEST_FROM#"
)

func NewFriendRequest(recipientUserId string, sender *User) *FriendRequest {
	fr := FriendRequest{
		RecipientId:             recipientUserId,
		SenderId:                sender.Userid,
		Date:                    utils.GetTimeNow(),
		SendersName:             sender.Name,
		SendersNickname:         sender.Nickname,
		SendersDpUrl:            sender.DpUrl,
		SendersDefaultPicColors: sender.DefaultProfilePicColors,
	}

	return &fr
}

func FriendRequestKey(recipientUserId string, senderUserId string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(FriendRequestPkPrefix, recipientUserId)},
		"sk": &types.AttributeValueMemberS{Value: utils.AddPrefix(FriendRequestSkPrefix, senderUserId)},
	}
}

func (fr FriendRequest) DatabaseFormat() (*map[string]types.AttributeValue, error) {
	fr.RecipientId = utils.AddPrefix(FriendRequestPkPrefix, fr.RecipientId)
	fr.SenderId = utils.AddPrefix(FriendRequestSkPrefix, fr.SenderId)
	item, err := attributevalue.MarshalMap(fr)

	if err != nil {
		log.Print("an error occurred while trying to marshal friend req item")
		return nil, err
	}

	return &item, nil
}

func ConvertToFriendRequest(item *map[string]types.AttributeValue) (*FriendRequest, error) {
	var fr FriendRequest
	if err := attributevalue.UnmarshalMap(*item, &fr); err != nil {
		return nil, err
	}

	fr.RecipientId = strings.TrimPrefix(fr.RecipientId, FriendRequestPkPrefix)
	fr.SenderId = strings.TrimPrefix(fr.SenderId, FriendRequestSkPrefix)

	return &fr, nil
}
