package models

import (
	"breadcrumb-backend-go/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	friendItemPk = "USER#"
	friendItemSk = "FRIEND#"
)

type friend struct {
	User1Id string `dynamodbav:"pk"`
	User2Id string `dynamodbav:"sk"`
	Date    string `dynamodbav:"date"`
}

func NewFriendship(user1id, user2id string) *friend {
	// Returns 2 friendship items
	return &friend{
		User1Id: user1id,
		User2Id: user2id,
		Date:    utils.GetTimeNow(),
	}
}

func FriendKey(user1Id string, user2Id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(friendItemPk, user1Id)},
		"sk": &types.AttributeValueMemberS{Value: utils.AddPrefix(friendItemSk, user2Id)},
	}
}

func (f *friend) DatabaseFormat() (*map[string]types.AttributeValue, error) {
	f.User1Id = utils.AddPrefix(friendItemPk, f.User1Id)
	f.User2Id = utils.AddPrefix(friendItemSk, f.User2Id)

	item, err := attributevalue.MarshalMap(f)

	if err != nil {
		log.Println("An error occurred while trying to marshal friendship item.")
		return nil, err
	}

	return &item, nil
}
