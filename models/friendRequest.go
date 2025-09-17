package models

import (
	"breadcrumb-backend-go/utils"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type friendRequest struct {
	RecipientId string `dynamodbav:"pk" json:"recipient"`
	SenderId    string `dynamodbav:"sk" json:"sender"`
	Date        string `dynamodbav:"date" json:"date"`
}

const (
	friendRequestPkPrefix = "USER#"
	friendRequestSkPrefix = "FRIEND_REQUEST_FROM#"
)

func NewFriendRequest(recipientUserId string, senderUserId string) *friendRequest {
	fr := friendRequest{
		RecipientId: recipientUserId,
		SenderId:    senderUserId,
		Date:        utils.GetTimeNow(),
	}

	return &fr
}

func FriendRequestKey(recipientUserId string, senderUserId string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(friendRequestPkPrefix, recipientUserId)},
		"sk": &types.AttributeValueMemberS{Value: utils.AddPrefix(friendRequestSkPrefix, senderUserId)},
	}
}

func (fr friendRequest) DatabaseFormat() (*map[string]types.AttributeValue, error) {
	fr.RecipientId = utils.AddPrefix(friendRequestPkPrefix, fr.RecipientId)
	fr.SenderId = utils.AddPrefix(friendRequestSkPrefix, fr.SenderId)
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

	fr.RecipientId = strings.TrimPrefix(fr.RecipientId, friendRequestPkPrefix)
	fr.SenderId = strings.TrimPrefix(fr.SenderId, friendRequestSkPrefix)

	return &fr, nil
}
