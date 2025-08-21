package apis

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GetProfileDetailsDependencies struct {
	UserPoolId    string
	CognitoClient *cognitoidentityprovider.Client
	TableName     string
	DdbClient     *dynamodb.Client
}
