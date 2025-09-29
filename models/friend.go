package models

import (
	"breadcrumb-backend-go/utils"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	FriendItemPk = "USER#"
	FriendItemSk = "FRIEND#"
)

type friend struct {
	ThisUserId          string `dynamodbav:"pk"`
	OtherUserID         string `dynamodbav:"sk"`
	UserDisplayInfoNoId        // this denormalized info of the other user
	Date                string `dynamodbav:"date"`
}

func NewFriendship(thisUserId string, otherUser *User) *friend {
	// Returns 2 friendship items
	return &friend{
		ThisUserId:          thisUserId,
		OtherUserID:         otherUser.Userid,
		UserDisplayInfoNoId: *GetUserDisplayInfoNoId(otherUser),
		Date:                utils.GetTimeNow(),
	}
}

func FriendKey(thisUserId string, otherUserId string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(FriendItemPk, thisUserId)},
		"sk": &types.AttributeValueMemberS{Value: utils.AddPrefix(FriendItemSk, otherUserId)},
	}
}

func (f friend) DatabaseFormat() (*map[string]types.AttributeValue, error) {
	f.ThisUserId = utils.AddPrefix(FriendItemPk, f.ThisUserId)
	f.OtherUserID = utils.AddPrefix(FriendItemSk, f.OtherUserID)

	item, err := attributevalue.MarshalMap(f)

	if err != nil {
		log.Println("An error occurred while trying to marshal friendship item.")
		return nil, err
	}

	return &item, nil
}

// TODO: write unit test for this function
func FriendItemsToUserDisplayStructs(items *[]map[string]types.AttributeValue) (*[]UserDisplayInfo, error) {
	var f []friend
	if err := attributevalue.UnmarshalListOfMaps(*items, &f); err != nil {
		log.Println("Failed to unmarshal db items to friend struct")
		return nil, err
	}

	var users []UserDisplayInfo
	for index, friend := range f {
		// clean up
		f[index].ThisUserId = strings.TrimPrefix(friend.ThisUserId, FriendItemPk)
		f[index].OtherUserID = strings.TrimPrefix(friend.OtherUserID, FriendItemSk)

		users = append(users, UserDisplayInfo{
			Userid:                  f[index].OtherUserID,
			Nickname:                f[index].Nickname,
			Name:                    f[index].Name,
			DpUrl:                   f[index].DpUrl,
			DefaultProfilePicColors: f[index].DefaultProfilePicColors,
		})
	}

	return &users, nil
}
