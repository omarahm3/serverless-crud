package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/omarahm3/sls/pkg/user"
)

type (
	Request  = events.APIGatewayProxyRequest
	Response = events.APIGatewayProxyResponse
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) (*Response, error) {
	// TODO this should []string slice
	emails := req.QueryStringParameters["email"]

	if len(emails) > 0 {
		return getSingleUser(emails, tableName, dynamoClient)
	}

	return getMultipleUsers(tableName, dynamoClient)
}

func CreateUser(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) (*Response, error) {
	result, err := user.CreateUser(req, tableName, dynamoClient)
	if err != nil {
		return errorResponse(err)
	}

	return response(http.StatusCreated, result)
}

func UpdateUser(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) (*Response, error) {
	result, err := user.UpdateUser(req, tableName, dynamoClient)
	if err != nil {
		return errorResponse(err)
	}

	return response(http.StatusCreated, result)
}

func DeleteUser(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) (*Response, error) {
	err := user.DeleteUser(req, tableName, dynamoClient)
	if err != nil {
		return errorResponse(err)
	}

	return response(http.StatusCreated, nil)
}

func UnhandledMethod() (*Response, error) {
	return response(http.StatusMethodNotAllowed, "method is not allowed")
}

func getSingleUser(email, tableName string, dynamoClient *dynamodb.DynamoDB) (*Response, error) {
	// TODO FetchUser must accept string email
	u, err := user.FetchUser(email, tableName, dynamoClient)
	if err != nil {
		return response(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	return response(http.StatusOK, u)
}

func getMultipleUsers(tableName string, dynamoClient *dynamodb.DynamoDB) (*Response, error) {
	users, err := user.FetchUsers(tableName, dynamoClient)
	if err != nil {
		return response(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	return response(http.StatusOK, users)
}

func errorResponse(err error) (*Response, error) {
	return response(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
}
