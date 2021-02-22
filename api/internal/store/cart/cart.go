package cart

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	//DynamoDBPrefixCart Prexix for the shopping key
	DynamoDBPrefixCart = "CART"

	//DynamoDBPrefixItem Prexix for product item key
	DynamoDBPrefixItem = "ITEM"
)

var (
	//ErrCreateCart error returned if we failed to create the cart in the database
	ErrCreateCart = errors.New("CouldNotCreateCart")

	//ErrAddItem error returned if we failed to add item to the cart
	ErrAddItem = errors.New("CouldNotAddItem")
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

//CreateAndAddItem Creates a shopping cart and returns a pointer to a struct of type Cart
func (h *Handler) CreateAndAddItem(ctx context.Context, ni NewItem) (*Cart, error) {

	//Generate a unique ID for the shopping cart
	cartID := uuid.New().String()
	log.Debug().Msgf("Generated UUID: %s", cartID)

	c := Cart{CartID: cartID}

	_, err := h.svc.TransactWriteItemsWithContext(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					Item: map[string]*dynamodb.AttributeValue{
						"pk":     {S: aws.String(getCartPK(c.CartID))},
						"sk":     {S: aws.String(getCartPK(c.CartID))},
						"cartId": {S: aws.String(c.CartID)},
					},
					TableName:           aws.String(h.tableName),
					ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
				},
			},
			{
				Put: &dynamodb.Put{
					Item: map[string]*dynamodb.AttributeValue{
						"pk":          {S: aws.String(getCartPK(c.CartID))},
						"sk":          {S: aws.String(getItemSK(ni.ItemID))},
						"cartId":      {S: aws.String(c.CartID)},
						"itemId":      {S: aws.String(ni.ItemID)},
						"description": {S: aws.String(ni.Description)},
						"price":       {N: aws.String(fmt.Sprintf("%f", ni.Price))},
						"quantity":    {N: aws.String(strconv.Itoa(ni.Quantity))},
					},
					TableName:           aws.String(h.tableName),
					ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
				},
			},
		},
	})

	if err != nil {
		log.Error().Msgf("Error creating cart: %s", err.Error())
		return nil, ErrCreateCart
	}

	log.Info().Msgf("Cart created with ID: %s", c.CartID)

	return &c, nil
}

//AddItem Adds new item to the shopping cart
func (h *Handler) AddItem(ctx context.Context, ni NewItem) (*Cart, error) {

	if err := validate.Struct(ni); err != nil {
		log.Error().Msgf("Error validating struct: %s", err.Error())
		return nil, getValidationError(err)
	}

	log.Info().Msgf("Adding item %s to cart %s", ni.Description, ni.CartID)

	_, err := h.svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"pk":          {S: aws.String(getCartPK(ni.CartID))},
			"sk":          {S: aws.String(getItemSK(ni.ItemID))},
			"cartId":      {S: aws.String(ni.CartID)},
			"itemId":      {S: aws.String(ni.ItemID)},
			"description": {S: aws.String(ni.Description)},
			"price":       {N: aws.String(fmt.Sprintf("%f", ni.Price))},
			"quantity":    {N: aws.String(strconv.Itoa(ni.Quantity))},
		},
		TableName:           aws.String(h.tableName),
		ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
	})

	if err != nil {
		log.Error().Msgf("Error adding item: %s", err.Error())
		return nil, ErrCreateCart
	}

	log.Info().Msgf("Item %s added to cart %s", ni.ItemID, ni.CartID)

	return &Cart{}, nil

}

func getCartPK(cartID string) string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixCart, cartID)
}

func getItemSK(itemID string) string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixItem, itemID)
}
