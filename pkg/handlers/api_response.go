package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func response(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	strBody, _ := json.Marshal(body)

	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(strBody),
	}, nil
}
