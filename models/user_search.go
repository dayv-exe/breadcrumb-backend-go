package models

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/utils"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	UserSearchIndexPkPrefix  = "USER_INDEX#"
	UserSearchIndexPrefixLen = 2
)

type UserSearch struct {
	UserId   string `dynamodbav:"userid" json:"userid"`
	Nickname string `dynamodbav:"nickname" json:"nickname"`
	Name     string `dynamodbav:"name" json:"name"`
	DpUrl    string `dynamodbav:"dp_url" json:"dpUrl"`
}

type userSearchDbItem struct {
	Pk       string `dynamodbav:"pk"`
	Sk       string `dynamodbav:"sk"`
	UserId   string `dynamodbav:"userid" json:"userid"`
	Name     string `dynamodbav:"name"`
	Nickname string `dynamodbav:"nickname"`
	DpUrl    string `dynamodbav:"dp_url" json:"dpUrl"`
}

func (u *UserSearch) BuildSearchIndexes() ([]map[string]types.AttributeValue, error) {
	if len(u.Name) < UserSearchIndexPrefixLen || len(u.Nickname) < constants.MIN_USERNAME_CHARS {
		return nil, fmt.Errorf("Name or nickname is too short!")
	}

	u.Name = strings.ToLower(u.Name)
	u.Nickname = strings.ToLower(u.Nickname)

	var tokens []string
	tokens = append(tokens, utils.SplitOnDelimiter(utils.NormalizeString(u.Name), " ", "_", ".")...)
	tokens = append(tokens, utils.SplitOnDelimiter(utils.NormalizeString(u.Nickname), " ", "_", ".")...)

	var indexes []map[string]types.AttributeValue
	seen := make(map[string]struct{})

	for _, token := range tokens {
		// get values
		if len(strings.TrimSpace(token)) < UserSearchIndexPrefixLen || len(strings.TrimSpace(token[:UserSearchIndexPrefixLen])) < UserSearchIndexPrefixLen {
			continue // skip
		}

		pk := UserSearchIndexPkPrefix + token[:UserSearchIndexPrefixLen]
		sk := token + "#" + u.UserId

		if _, ok := seen[pk+"|"+sk]; !ok {
			// seen now if not seen before
			seen[pk+"|"+sk] = struct{}{}

			// new index item
			new := userSearchDbItem{
				Pk:       pk,
				Sk:       sk,
				UserId:   u.UserId,
				Name:     u.Name,
				Nickname: u.Nickname,
				DpUrl:    u.DpUrl,
			}

			item, err := attributevalue.MarshalMap(new)
			fmt.Println(new)
			if err != nil {
				return nil, fmt.Errorf("An error occurred while marshaling user search db item: %w", err)
			}

			indexes = append(indexes, item)
		}
	}

	return indexes, nil
}

func GetUserSearchIndexesKeys(dbIndexItems []map[string]types.AttributeValue) []map[string]types.AttributeValue {
	var keys []map[string]types.AttributeValue
	seen := make(map[string]struct{})

	for _, item := range dbIndexItems {
		pk, pkOk := item["pk"].(*types.AttributeValueMemberS)
		sk, skOk := item["sk"].(*types.AttributeValueMemberS)

		if pkOk && skOk {
			// if is valid database index item
			compKey := pk.Value + sk.Value

			if _, ok := seen[compKey]; !ok {
				// not seen before
				// make it seen now
				seen[compKey] = struct{}{}

				// append
				keys = append(keys, map[string]types.AttributeValue{
					"pk": pk,
					"sk": sk,
				})
			}
		}
	}

	return keys
}
