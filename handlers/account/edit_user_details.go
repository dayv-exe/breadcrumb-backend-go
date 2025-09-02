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
	Target  string          `json:"target"`
	Payload json.RawMessage `json:"payload"`
}

type editNamePayload struct {
	Nickname string `json:"nickname"`
	Name     string `json:"fullname"`
}

func (deps *EditUserDetailsDependency) HandleEditUserDetails(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody editUserDetailsReq
	if umErr := json.Unmarshal([]byte(req.Body), &reqBody); umErr != nil {
		return models.InvalidRequestErrorResponse(""), nil
	}

	userId := utils.GetAuthUserId(req)

	if userId == "" {
		return models.UnauthorizedErrorResponse(""), nil
	}

	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	// allows update of nickname, full name, bio

	switch reqBody.Target {
	case "names":
		// update nickname and fullname (its called username in the front end)
		var namePayload editNamePayload
		if umErr := json.Unmarshal([]byte(reqBody.Payload), &namePayload); umErr != nil {
			return models.InvalidRequestErrorResponse(""), nil
		}

		nnErr := dbHelper.UpdateNicknameAndFullname(userId, namePayload.Nickname, namePayload.Name)
		if nnErr != nil {
			return models.ServerSideErrorResponse("Something went wrong while trying to update nickname!", nnErr, "trying to update nickname"), nil
		}

	default:
		// invalid target
		return models.InvalidRequestErrorResponse(""), nil
	}

	return models.SuccessfulRequestResponse("Resource updated successfully!", false), nil
}
