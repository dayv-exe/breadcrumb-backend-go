package discover

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type AccountSearchDependencies struct {
	Client          *dynamodb.Client
	SearchTableName string
}

type SearchResult struct {
	models.UserSearch
	Type string `json:"type"`
}

func (deps *AccountSearchDependencies) HandleAccountSearch(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// /search/{search_string}

	limit := 30 // limit to 30 results
	searchStr := req.PathParameters["search_string"]
	if len(searchStr) < models.UserSearchIndexPrefixLen {
		// return empty array if string length is too short
		return models.SuccessfulGetRequestResponse([]models.UserSearch{}), nil
	}

	// only takes the first 2 - 20 chars in the search string to process and then search the database for
	maxChars := len((searchStr))

	if len(searchStr) > constants.MAX_SEARCH_STRING_CHARS {
		maxChars = constants.MAX_SEARCH_STRING_CHARS
	}

	searchStr = searchStr[:maxChars]

	searchHelper := helpers.SearchDynamoHelper{
		DbClient:        deps.Client,
		SearchTableName: deps.SearchTableName,
		Ctx:             ctx,
	}

	results, err := searchHelper.SearchUser(strings.ToLower(searchStr), int32(limit))

	if err != nil {
		return models.ServerSideErrorResponse("", err, "error from dynamo helper searching users"), nil
	}

	return models.SuccessfulGetRequestResponse(results), nil
}
