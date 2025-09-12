package models

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/utils"
	"fmt"
	"log"
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
	var indexes []map[string]types.AttributeValue

	nicknameItems, nnErr := u.buildNicknameIndexItem()
	if nnErr != nil {
		return nil, nnErr
	}

	fullnameItems, fnErr := u.buildFullnameIndexItem()
	if fnErr != nil {
		return nil, fnErr
	}

	indexes = append(indexes, nicknameItems...)
	indexes = append(indexes, fullnameItems...)
	return indexes, nil
}

func (u *UserSearch) buildNicknameIndexItem() ([]map[string]types.AttributeValue, error) {
	u.Nickname = strings.ToLower(u.Nickname)
	chunks := utils.SplitOnDelimiter(u.Nickname, "_", ".")
	return u.buildIndexItem(chunks, UserSearchIndexPkPrefix)
}

func (u *UserSearch) buildFullnameIndexItem() ([]map[string]types.AttributeValue, error) {
	u.Name = strings.ToLower(u.Name)
	u.Name = strings.TrimSpace(u.Name)
	chunks := utils.SplitOnDelimiter(utils.NormalizeString(u.Name), " ", "_", ".")
	// adds id as suffix to make each db item unique
	return u.buildIndexItem(chunks, UserSearchIndexPkPrefix)
}

func DeconstructSearchIndexItem(item *map[string]types.AttributeValue) (*UserSearch, error) {
	var i UserSearch
	if err := attributevalue.UnmarshalMap(*item, &i); err != nil {
		return nil, err
	}

	return &i, nil
}

func GetUserSearchIndexesKeys(indexItems []map[string]types.AttributeValue) []map[string]types.AttributeValue {
	var keys []map[string]types.AttributeValue
	for _, item := range indexItems {
		pkAttr, pkOk := item["pk"].(*types.AttributeValueMemberS)
		skAttr, skOk := item["sk"].(*types.AttributeValueMemberS)

		if pkOk && skOk {
			keys = append(keys, map[string]types.AttributeValue{
				"pk": pkAttr,
				"sk": skAttr,
			})
		}
	}

	return keys
}

func (u *UserSearch) buildIndexItem(tokens []string, pkPrefix string) ([]map[string]types.AttributeValue, error) {
	// returns a list of search index items based on name provided
	var items []map[string]types.AttributeValue
	for _, token := range tokens {
		if len(token) < UserSearchIndexPrefixLen || len(strings.ReplaceAll(token[:UserSearchIndexPrefixLen], " ", "")) < UserSearchIndexPrefixLen {
			// only index when the name is 2 or more chars long
			continue
		}
		new := userSearchDbItem{
			Pk:       utils.AddPrefix(pkPrefix, token[:UserSearchIndexPrefixLen]), // pk is "NICKNAME#{first to chars of nickname}"
			Sk:       token + utils.AddPrefix("#", u.UserId),
			UserId:   u.UserId,
			Name:     u.Name,
			Nickname: u.Nickname,
			DpUrl:    u.DpUrl,
		}
		log.Println(new)
		item, err := attributevalue.MarshalMap(new)

		if err != nil {
			return nil, fmt.Errorf("Failed to marshal search index item: %w", err)
		}

		items = append(items, item)
	}
	return items, nil
}

// not in use yet
func chunkifyName(name string) []string {
	// david
	name = utils.NormalizeString(name)
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, ".", "")
	n := len(name)
	var chunks []string

	if n == UserSearchIndexPrefixLen {
		// if the name is just short enough to already be a chunk
		return []string{name}
	} else if n < UserSearchIndexPrefixLen {
		// if the name is too short to chunkify
		return []string{}
	}

	var i = 0

	s := []rune(name)
	n = len(s)
	for i < n-(UserSearchIndexPrefixLen-1) && i < constants.MAX_CHUNKABLE_LEN {
		chunks = append(chunks, string(s[i:]))
		i++
	}
	return chunks
}
