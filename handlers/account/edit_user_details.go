package account

// TO EDIT USER DETAILS, PARSE THE ATTRIBUTE TO EDIT IN {ACTION}

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
	DdbClient   *dynamodb.Client
	UserTable   string
	SearchTable string
}

type editUserDetailsReq struct {
	Payload json.RawMessage `json:"payload"`
}

func (deps *EditUserDetailsDependency) HandleEditUserDetails(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get action and payload
	action := req.PathParameters["action"]
	var reqBody editUserDetailsReq
	if err := json.Unmarshal([]byte(req.Body), &reqBody); err != nil {
		// error while unmarshalling
		return models.InvalidRequestErrorResponse(""), nil
	}
	userid := utils.GetAuthUserId(req)

	// action based switch
	switch action {
	case "nickname":
		// change nickname
		return deps.handleChangeName(userid, reqBody, ctx, true)

	case "name":
		// change name
		return deps.handleChangeName(userid, reqBody, ctx, false)

	case "bio":
		// change bio
		userHelper := helpers.UserDynamoHelper{
			DbClient:      deps.DdbClient,
			UserTableName: deps.UserTable,
			Ctx:           ctx,
		}
		bioErr := userHelper.UpdateBio(userid, string(reqBody.Payload))
		if bioErr != nil {
			return models.ServerSideErrorResponse("", bioErr, "error while trying to update bio"), nil
		}

		return models.SuccessfulRequestResponse("", false), nil

	default:
		return models.InvalidRequestErrorResponse(""), nil
	}
}

func (deps *EditUserDetailsDependency) handleChangeName(userid string, reqBody editUserDetailsReq, ctx context.Context, updateNickname bool) (events.APIGatewayProxyResponse, error) {
	// handles updating both name and nickname
	userHelper := helpers.UserDynamoHelper{
		DbClient:      deps.DdbClient,
		UserTableName: deps.UserTable,
		Ctx:           ctx,
	}

	// gets the user from db
	user, uErr := userHelper.FindById(userid)
	if uErr != nil {
		return models.ServerSideErrorResponse("Failed to get user", uErr, "failed to get user and extract last name changed date"), nil
	}

	// to toggle between checking for last nickname change or last name change depending on weather we want to change name or nickname
	lastChangeDate := user.LastNameChange
	if updateNickname {
		lastChangeDate = user.LastNicknameChange
	}

	allowed, ncErr := utils.NameChangeAllowed(lastChangeDate)
	if ncErr != nil {
		return models.ServerSideErrorResponse("Failed to get last changed date", uErr, "failed to get last name or nickname changed date"), nil
	}

	if !allowed {
		return models.InvalidRequestErrorResponse("You can only change your name or nickname once in 30 days."), nil
	}

	err := userHelper.UpdateName(user, string(reqBody.Payload), updateNickname)

	if err != nil {
		return models.ServerSideErrorResponse("", err, "error while trying to update name or nickname"), nil
	}

	return models.SuccessfulRequestResponse("", false), nil
}
