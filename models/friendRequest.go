package models

import (
	"breadcrumb-backend-go/utils"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type friendRequest struct {
	RecipientId string `dynamodbav:"pk"`
	SenderId    string `dynamodbav:"sk" json:"userId"`
	UserDisplayInfoNoId
	Date string `dynamodbav:"date" json:"date"`
}

const (
	FriendRequestPkPrefix = "USER#"
	FriendRequestSkPrefix = "FRIEND_REQUEST_FROM#"
)

func NewFriendRequest(recipientUserId string, sender *User) *friendRequest {
	fr := friendRequest{
		RecipientId: recipientUserId,
		SenderId:    sender.Userid,
		Date:        utils.GetTimeNow(),
		UserDisplayInfoNoId: UserDisplayInfoNoId{
			Nickname:                sender.Nickname,
			Name:                    sender.Name,
			DpUrl:                   sender.DpUrl,
			DefaultProfilePicColors: sender.DefaultProfilePicColors,
		},
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

func FriendRequestItemsToUserDisplayStructs(item *[]map[string]types.AttributeValue) (*[]UserDisplayInfo, error) {
	log.Print("inside item to struct function")
	// TODO: write a unit test for this function
	// takes friends request items from the database and converts them to user display info
	// user id, nickname, name and display picture

	var requests []friendRequest
	if err := attributevalue.UnmarshalListOfMaps(*item, &requests); err != nil {
		return nil, err
	}

	log.Print("done unmarshaling")

	var users []UserDisplayInfo

	for index, request := range requests {
		requests[index].SenderId = strings.TrimPrefix(request.SenderId, FriendRequestSkPrefix)

		users = append(users, UserDisplayInfo{
			Userid:                  requests[index].SenderId,
			Nickname:                request.Nickname,
			Name:                    request.Name,
			DpUrl:                   request.DpUrl,
			DefaultProfilePicColors: request.DefaultProfilePicColors,
		})
	}

	log.Print("done looping through and extracting")

	return &users, nil
}
