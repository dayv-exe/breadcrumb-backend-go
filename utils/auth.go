package utils

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func GetAuthUserId(req *events.APIGatewayProxyRequest) string {
	claims, ok := req.RequestContext.Authorizer["claims"].(map[string]interface{})
	log.Println(req.RequestContext.Authorizer)
	if !ok {
		// missing claims
		log.Println("missing claims")
		return ""
	}

	log.Println(claims)

	sub, ok := claims["sub"].(string)
	if !ok {
		// missing sub
		log.Println("missing sub")
		return ""
	}

	return sub
}
