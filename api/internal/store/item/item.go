package item

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog/log"
)

const (
	//DynamoDBPrefixItem Prexix for product item key
	DynamoDBPrefixItem = "ITEM#"
)

var (
	//ErrStoreTableNameIsEmpty Error describes when DynamoDB table name is empty
	ErrStoreTableNameIsEmpty = "StoreTableNameIsEmpty"

	//ErrCouldNotLoadItems error returned if we failed to load the cart
	ErrCouldNotLoadItems = errors.New("CouldNotLoadItems")
)

//Handler struct is a handler for executing the actions related to the shopping cart
type Handler struct {
	svc       dynamodbiface.DynamoDBAPI
	tableName string
}

//New returns pointer to a struct of type Cart, that contains methods
//For each action that can be executed on this API
func New(svc dynamodbiface.DynamoDBAPI, tableName string) (*Handler, error) {
	if tableName == "" {
		log.Error().Msg("Table name is empty")
		return nil, errors.New(ErrStoreTableNameIsEmpty)
	}

	return &Handler{svc, tableName}, nil
}

//List returns the information of an item
//Eventuallly, this method will receive a category id
//For now, it loads all items from category 1
func (h *Handler) List(ctx context.Context) (*List, error) {

	result, err := h.svc.QueryWithContext(ctx, &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"pk": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					//Load items from category 1
					{S: aws.String("CATEGORY#1")},
				},
			},
			"sk": {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(DynamoDBPrefixItem)},
				},
			},
		},
		ProjectionExpression: aws.String("item_id,description,price"),
		TableName:            aws.String(h.tableName),
	})

	if err != nil {
		log.Error().Msgf("Error loading items: %s", err.Error())
		return nil, ErrCouldNotLoadItems
	}

	l := List{}

	if *result.Count == int64(0) {
		return &l, nil
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &l.Items)
	if err != nil {
		log.Error().Msgf("Error Unmarshaling items: %s", err.Error())
		return nil, ErrCouldNotLoadItems
	}

	return &l, nil
}
