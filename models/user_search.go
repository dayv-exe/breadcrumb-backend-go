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
	UserDisplayInfo
	Rating int `json:"rating"`
}

type userSearchDbItem struct {
	Pk     string `dynamodbav:"pk"`
	Sk     string `dynamodbav:"sk"`
	Userid string `dynamodbav:"userid" json:"userid"`
	UserDisplayInfoNoId
}

func (u *UserSearch) BuildSearchIndexes() ([]map[string]types.AttributeValue, error) {
	// returns items to be put in the database that contains search index
	// for example
	// john.test this function will return
	// pk: jo, sk: john
	// pk: te, sk: test

	if len(u.Nickname) < constants.MIN_USERNAME_CHARS {
		return nil, fmt.Errorf("Name or nickname is too short!")
	}

	u.Nickname = strings.ToLower(u.Nickname)

	var tokens []string
	tokens = append(tokens, utils.SplitOnDelimiter(strings.ToLower(utils.NormalizeString(u.Name)), " ", "_", ".")...)
	tokens = append(tokens, utils.SplitOnDelimiter(utils.NormalizeString(u.Nickname), " ", "_", ".")...)

	var indexes []map[string]types.AttributeValue
	seen := make(map[string]struct{})

	for _, token := range tokens {
		// get values
		if len(strings.TrimSpace(token)) < UserSearchIndexPrefixLen || len(strings.TrimSpace(token[:UserSearchIndexPrefixLen])) < UserSearchIndexPrefixLen {
			continue // skip
		}

		pk := UserSearchIndexPkPrefix + token[:UserSearchIndexPrefixLen]
		sk := token + "#" + strings.Trim(u.Userid, UserPkPrefix)

		if _, ok := seen[pk+"|"+sk]; !ok {
			// seen now if not seen before
			seen[pk+"|"+sk] = struct{}{}

			// new index item
			new := userSearchDbItem{
				Pk:     pk,
				Sk:     sk,
				Userid: u.Userid,
				UserDisplayInfoNoId: UserDisplayInfoNoId{
					Nickname:                u.Nickname,
					Name:                    u.Name,
					DpUrl:                   u.DpUrl,
					DefaultProfilePicColors: u.DefaultProfilePicColors,
				},
			}

			item, err := attributevalue.MarshalMap(new)
			if err != nil {
				return nil, fmt.Errorf("An error occurred while marshaling user search db item: %w", err)
			}

			indexes = append(indexes, item)
		}
	}

	return indexes, nil
}

func GetUserSearchIndexesKeys(dbIndexItems []map[string]types.AttributeValue) []map[string]types.AttributeValue {
	// this function returns a slice containing the pk and sk of database items parsed into it
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
