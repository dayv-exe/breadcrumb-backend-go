package utils

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func GetAuthUserId(req *events.APIGatewayProxyRequest) string {
	log.Println("getting user sub from request")
	jwtMap, ok := req.RequestContext.Authorizer["jwt"].(map[string]interface{})
	if !ok {
		log.Println("missing jwt map")
		return ""
	}

	claims, ok := jwtMap["claims"].(map[string]interface{})
	if !ok {
		log.Println("missing claims")
		return ""
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		log.Println("missing sub")
		return ""
	}

	log.Println("gotten user sub from request")
	return sub
}

func IsAuthenticatedUser(req *events.APIGatewayProxyRequest, userId string) bool {
	sub := GetAuthUserId(req)
	if sub == "" {
		return false
	}
	return sub == userId
}
