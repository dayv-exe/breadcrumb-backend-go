// delete users that have not verified their email within 24hours

package auth

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/helpers"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	currLoops = 0
)

type RemoveStaleAccountsDependencies struct {
	Client     *cognitoidentityprovider.Client
	UserPoolId string
}

func (deps *RemoveStaleAccountsDependencies) HandleRemoveStaleAccounts(ctx context.Context) error {
	var paginationToken *string
	count := 0
	cognitoHelper := helpers.UserCognitoHelper{
		UserPoolId:    deps.UserPoolId,
		CognitoClient: deps.Client,
		Ctx:           ctx,
	}

	for {
		input := &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(deps.UserPoolId),
			Filter:     aws.String("cognito:user_status=\"UNCONFIRMED\""),
			Limit:      aws.Int32(60),
		}

		if paginationToken != nil {
			input.PaginationToken = paginationToken
		}

		res, err := deps.Client.ListUsers(ctx, input)
		if err != nil {
			return fmt.Errorf("FAILED TO LIST USERS: %w", err)
		}
		log.Printf("Found %v accounts to remove", len(res.Users))

		for _, user := range res.Users {
			if user.UserCreateDate == nil {
				// skip if we dont know account age (unfortunately)
				continue
			}

			if time.Since(*user.UserCreateDate) < 24*time.Hour {
				//skip if account is less than 24 hours old
				log.Printf("%s is less than 24 hours old, %v", *user.Username, user.UserCreateDate)
				continue
			}

			// delete user at this point (if they are not verified)
			username := aws.ToString(user.Username)
			uErr := cognitoHelper.DeleteFromCognito(username, false)

			if uErr != nil {
				log.Printf("FAILED TO DELETE USER %s: %v", username, uErr)
			} else {
				log.Printf("Removing %s", username)
				count++
			}
		}

		// escape clause
		currLoops++
		if res.PaginationToken == nil {
			break
		}
		if currLoops >= constants.MAX_STALE_ACCOUNTS_LOOPS {
			log.Print("Maximum loops reached, ending function execution prematurely.")
			break
		}

		paginationToken = res.PaginationToken
	}

	return nil
}
