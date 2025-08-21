package utils

import "github.com/aws/aws-lambda-go/events"

func GetAuthUserId(req *events.APIGatewayProxyRequest) string {
	claims, ok := req.RequestContext.Authorizer["claims"].(map[string]interface{})
	if !ok {
		// missing claims
		return ""
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		// missing sub
		return ""
	}

	return sub
}
