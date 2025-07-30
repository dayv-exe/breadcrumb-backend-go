// delete users that have not verified their email within 24hours

package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	maxLoops  = 5
	currLoops = 0
)

type RemoveStaleAccountsDependencies struct {
	Client     *cognitoidentityprovider.Client
	UserPoolId string
}

func (rsaDeps RemoveStaleAccountsDependencies) HandleRemoveStaleAccounts(ctx context.Context) error {
	var paginationToken *string
	count := 0

	for {
		input := &cognitoidentityprovider.ListUsersInput{
			UserPoolId: &rsaDeps.UserPoolId,
			Filter:     aws.String("cognito:user_status=\"UNCONFIRMED\""),
			Limit:      aws.Int32(60),
		}

		if paginationToken != nil {
			input.PaginationToken = paginationToken
		}

		res, err := rsaDeps.Client.ListUsers(ctx, input)
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

			username := aws.ToString(user.Username)
			_, uErr := rsaDeps.Client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: &rsaDeps.UserPoolId,
				Username:   aws.String(username),
			})

			if uErr != nil {
				log.Printf("FAILED TO DELETE USER %s: %v", username, uErr)
			} else {
				log.Printf("Removing %s", username)
				count++
			}
		}

		// escape clause
		currLoops++
		if res.PaginationToken == nil || currLoops >= maxLoops {
			break
		}

		paginationToken = res.PaginationToken
	}

	return nil
}
