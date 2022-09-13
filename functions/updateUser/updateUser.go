package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/omarahm3/sls/pkg/handlers"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	handler, err := handlers.Prepare(req)
	if err != nil {
		return handlers.JSONResponse(http.StatusInternalServerError, err.Error())
	}

	return handler.UpdateUser()
}

func main() {
	lambda.Start(handler)
}
