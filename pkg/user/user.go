package user

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/omarahm3/sls/pkg/validators"
)

type User struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Request = events.APIGatewayProxyRequest

var (
	ErrorFailedToFetchUser     = fmt.Errorf("failed to fetch user/s")
	ErrorFailedToUnmarshalUser = fmt.Errorf("failed to unmarshal user")
	ErrorFailedToMarshalUser   = fmt.Errorf("failed to marshal user")
	ErrorInvalidUserData       = fmt.Errorf("invalid user data")
	ErrorUserAlreadyExists     = fmt.Errorf("user already exists")
	ErrorCouldNotCreateUser    = fmt.Errorf("could not create user")
	ErrorUserDoesNotExist      = fmt.Errorf("user does not exist")
	ErrorCouldNotUpdateUser    = fmt.Errorf("could not update user")
	ErrorCouldNotDeleteUser    = fmt.Errorf("could not delete user")
	ErrorEmailIsNotValid       = fmt.Errorf("email is not valid")
)

func FetchUser(email, tableName string, dynamoClient *dynamodb.DynamoDB) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		return nil, err
	}

	u := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, &u)
	if err != nil {
		return nil, ErrorFailedToUnmarshalUser
	}

	if u.Email == "" {
		return nil, nil
	}

	return u, nil
}

func FetchUsers(tableName string, dynamoClient *dynamodb.DynamoDB) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.Scan(input)
	if err != nil {
		return nil, ErrorFailedToFetchUser
	}

	users := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, ErrorFailedToUnmarshalUser
	}

	return users, nil
}

func CreateUser(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) (*User, error) {
	var u User

	err := json.Unmarshal([]byte(req.Body), &u)
	if err != nil {
		return nil, ErrorInvalidUserData
	}

	if !validators.IsValidEmail(u.Email) {
		return nil, ErrorEmailIsNotValid
	}

	current, _ := FetchUser(u.Email, tableName, dynamoClient)
	if current != nil && len(current.Email) != 0 {
		return nil, ErrorUserAlreadyExists
	}

	val, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, ErrorUserAlreadyExists
	}

	_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
		Item:      val,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, ErrorCouldNotCreateUser
	}

	return &u, nil
}

func UpdateUser(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) (*User, error) {
	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, ErrorInvalidUserData
	}

	current, _ := FetchUser(u.Email, tableName, dynamoClient)
	if current != nil && len(current.Email) == 0 {
		return nil, ErrorUserDoesNotExist
	}

	value, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, ErrorFailedToMarshalUser
	}

	_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
		Item:      value,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, ErrorCouldNotUpdateUser
	}

	return &u, nil
}

func DeleteUser(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		return ErrorCouldNotDeleteUser
	}

	return nil
}
