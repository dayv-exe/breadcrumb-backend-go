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
	// TODO: write a unit test for this function
	// takes friends request items from the database and converts them to user display info
	// user id, nickname, name and display picture

	var fr []friendRequest
	if err := attributevalue.UnmarshalListOfMaps(*item, &fr); err != nil {
		return nil, err
	}

	var users []UserDisplayInfo

	for index, request := range fr {
		fr[index].RecipientId = strings.TrimPrefix(request.RecipientId, FriendRequestPkPrefix)
		fr[index].SenderId = strings.TrimPrefix(request.SenderId, FriendRequestSkPrefix)

		users = append(users, UserDisplayInfo{
			Userid:                  fr[index].SenderId,
			Nickname:                fr[index].Nickname,
			Name:                    fr[index].Name,
			DpUrl:                   fr[index].DpUrl,
			DefaultProfilePicColors: fr[index].DefaultProfilePicColors,
		})
	}

	return &users, nil
}
