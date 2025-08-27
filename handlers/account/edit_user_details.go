package account

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type EditUserDetailsDependency struct {
	DdbClient *dynamodb.Client
	TableName string
}

type editUserDetailsReq struct {
	Target  string `json:"target"`
	Payload string `json:"payload"`
}

func (deps *EditUserDetailsDependency) HandleEditUserDetails(ctx context.Context, req *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var reqBody editUserDetailsReq
	if umErr := json.Unmarshal([]byte(req.Body), &reqBody); umErr != nil {
		return models.InvalidRequestErrorResponse("")
	}

	userId := utils.GetAuthUserId(req)

	if userId == "" {
		return models.UnauthorizedErrorResponse("")
	}

	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	// allows update of nickname, full name, bio

	switch reqBody.Target {
	case "nickname":
		// update nickname (its called username in the front end)
		newNickname := reqBody.Payload
		nnErr := dbHelper.UpdateNickname(userId, newNickname)
		if nnErr != nil {
			return models.ServerSideErrorResponse("Something went wrong while trying to update nickname!", nnErr, "trying to update nickname")
		}

	case "name":
		// update full name
		newName := reqBody.Payload
		nErr := dbHelper.UpdateName(userId, newName)

		if nErr != nil {
			return models.ServerSideErrorResponse("Something went wrong while trying to update full name!", nErr, "trying to update full name")
		}

	default:
		// invalid target
		return models.InvalidRequestErrorResponse("")
	}

	return models.SuccessfulRequestResponse("Resource updated successfully!", false)
}
