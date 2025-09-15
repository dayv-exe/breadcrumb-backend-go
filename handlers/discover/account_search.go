package discover

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type AccountSearchDependencies struct {
	Client    *dynamodb.Client
	TableName string
}

type SearchResult struct {
	models.UserSearch
	Type string `json:"type"`
}

func (deps *AccountSearchDependencies) HandleAccountSearch(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// /search/{search_string}

	searchStr := req.PathParameters["search_string"]
	if len(searchStr) < models.UserSearchIndexPrefixLen {
		// return empty array if string length is too short
		return models.SuccessfulGetRequestResponse([]models.UserSearch{}), nil
	}

	searchHelper := helpers.SearchDynamoHelper{
		DbClient:  deps.Client,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	results, err := searchHelper.SearchUser(searchStr)

	if err != nil {
		return models.ServerSideErrorResponse("", err, "error from dynamo helper searching users"), nil
	}

	return models.SuccessfulGetRequestResponse(results), nil
}
