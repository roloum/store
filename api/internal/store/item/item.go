package item

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog/log"
)

const (
	//DynamoDBPrefixItem Prexix for product item key
	DynamoDBPrefixCategory = "CATEGORY#"

	//DynamoDBPrefixItem Prexix for product item key
	DynamoDBPrefixItem = "ITEM#"
)

var (
	//ErrStoreTableNameIsEmpty Error describes when DynamoDB table name is empty
	ErrStoreTableNameIsEmpty = "StoreTableNameIsEmpty"

	//ErrCouldNotLoadItems error returned if we failed to load the cart
	ErrCouldNotLoadItems = errors.New("CouldNotLoadItems")

	//ErrCategoryIDIsEmpty error returned if the categoryID is empty
	ErrCategoryIDIsEmpty = errors.New("CategoryIDIsEmpty")
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
//It uses a GSI to load the items based on categoryID
func (h *Handler) List(ctx context.Context, categoryID string) (*List, error) {

	if categoryID == "" {
		return nil, ErrCategoryIDIsEmpty
	}

	log.Debug().Msgf("Loading items for categoryID: %s", categoryID)

	result, err := h.svc.QueryWithContext(ctx, &dynamodb.QueryInput{
		IndexName: aws.String("gsi1pk"),
		KeyConditions: map[string]*dynamodb.Condition{
			"gsi1pk": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					//Load items from category 1
					{S: aws.String(getItemGSI1SK(categoryID))},
				},
			},
			"gsi1sk": {
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

	l := List{Items: []Item{}}

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

//getItemGSI1SK returns the categoryID formatted for the gsi1pk
func getItemGSI1SK(categoryID string) string {
	return fmt.Sprintf("%s%s", DynamoDBPrefixCategory, categoryID)
}
