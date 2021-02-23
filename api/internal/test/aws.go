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

	PutItemOutput            *dynamodb.PutItemOutput
	UpdateItemOutput         *dynamodb.UpdateItemOutput
	TransactWriteItemsOutput *dynamodb.TransactWriteItemsOutput
	QueryOutput              *dynamodb.QueryOutput
	OutputError              error
}

//PutItemWithContext mocks the PutItemWithContext method
func (m *MockDynamoDB) PutItemWithContext(aws.Context, *dynamodb.PutItemInput,
	...request.Option) (*dynamodb.PutItemOutput, error) {

	return m.PutItemOutput, m.OutputError
}

//TransactWriteItemsWithContext mocks the TransactWriteItemsWithContext method
func (m *MockDynamoDB) TransactWriteItemsWithContext(aws.Context,
	*dynamodb.TransactWriteItemsInput, ...request.Option) (
	*dynamodb.TransactWriteItemsOutput, error) {
	return m.TransactWriteItemsOutput, m.OutputError
}

//UpdateItemWithContext mocks the UpdateItemWithContext method
func (m *MockDynamoDB) UpdateItemWithContext(aws.Context,
	*dynamodb.UpdateItemInput, ...request.Option) (*dynamodb.UpdateItemOutput,
	error) {
	return m.UpdateItemOutput, m.OutputError
}

//QueryWithContext mocks the QueryWithContext method
func (m *MockDynamoDB) QueryWithContext(aws.Context, *dynamodb.QueryInput,
	...request.Option) (*dynamodb.QueryOutput, error) {
	return &dynamodb.QueryOutput{Count: aws.Int64(0)}, m.OutputError
}
