package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

//MockDynamoDB Is a Mock DynamoDB client that satisfies the DynamoDBAPI
//It implements the necessary methods to mock the dynamoDB calls that
//The API will make
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI

	PutItemOutput *dynamodb.PutItemOutput
	OutputError   error
}

//PutItemWithContext mocks the PutItemWithContext method
func (m *MockDynamoDB) PutItemWithContext(aws.Context, *dynamodb.PutItemInput,
	...request.Option) (*dynamodb.PutItemOutput, error) {

	return m.PutItemOutput, m.OutputError
}
