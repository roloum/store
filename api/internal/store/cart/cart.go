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
	//ErrCreateCartWithExistingCartID error returned when attempting to create
	//A shopping cart while sending a cartID in the newItem struct
	ErrCreateCartWithExistingCartID = errors.New("CreateCartWithExistingCartID")

	//ErrCreateCart error returned if we failed to create the cart in the database
	ErrCreateCart = errors.New("CouldNotCreateCart")

	//ErrCouldNotAddItem error returned if we failed to create the cart in the database
	ErrCouldNotAddItem = errors.New("CouldNotAddItem")

	//ErrCouldNotUpdateItem error returned if we failed to update item's quantity
	ErrCouldNotUpdateItem = errors.New("CouldNotUpdateItem")
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

//CreateAndAddItem Creates a shopping cart and adds the first item
//ni contains the information about the new item
func (h *Handler) CreateAndAddItem(ctx context.Context, ni *NewItemInfo) (*Cart, error) {

	if ni.CartID != "" {
		log.Error().Msgf("Error attempting to create a cart with an ID: %s", ni.CartID)
		return nil, ErrCreateCartWithExistingCartID
	}

	//Generate a unique ID for the shopping cart
	cartID := uuid.New().String()
	log.Debug().Msgf("Generated UUID: %s", cartID)

	ni.CartID = cartID

	//Validate the struct after the cartID has been added to the newItem
	if err := validate.Struct(ni); err != nil {
		return nil, getValidationError(err)
	}

	log.Debug().Msgf("Creating cart with ID: %s and adding item ID :%s",
		ni.CartID, ni.ItemID)

	_, err := h.svc.TransactWriteItemsWithContext(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					Item: map[string]*dynamodb.AttributeValue{
						"pk":     {S: aws.String(getCartPK(ni.CartID))},
						"sk":     {S: aws.String(getCartPK(ni.CartID))},
						"cartId": {S: aws.String(ni.CartID)},
					},
					TableName:           aws.String(h.tableName),
					ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
				},
			},
			{
				Put: &dynamodb.Put{
					Item:                getNewItemDynamoAttributes(ni),
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

	c := &Cart{CartID: ni.CartID}
	log.Info().Msgf("Cart created with ID: %s", c.CartID)

	return c, nil
}

//AddItem Adds new item to the shopping cart
//Receives the NewItemInfo with all the information about the new item
func (h *Handler) AddItem(ctx context.Context, ni *NewItemInfo) (*Cart, error) {

	if err := validate.Struct(ni); err != nil {
		log.Error().Msgf("Error validating struct: %s", err.Error())
		return nil, getValidationError(err)
	}

	log.Debug().Msgf("Adding item %s to cart %s", ni.Description, ni.CartID)

	_, err := h.svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:                getNewItemDynamoAttributes(ni),
		TableName:           aws.String(h.tableName),
		ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
	})

	if err != nil {
		log.Error().Msgf("Error adding item: %s", err.Error())
		return nil, ErrCouldNotAddItem
	}

	log.Info().Msgf("Item %s added to cart %s", ni.ItemID, ni.CartID)

	return &Cart{}, nil

}

//UpdateItem Updates the quantity for an item in the shopping cart
func (h *Handler) UpdateItem(ctx context.Context, ui *UpdateItemInfo) (*Cart, error) {

	if err := validate.Struct(ui); err != nil {
		log.Error().Msgf("Error validating struct: %s", err.Error())
		return nil, getValidationError(err)
	}

	log.Debug().Msgf("Updating quantity: %d for item %s in cart %s", ui.Quantity,
		ui.ItemID, ui.CartID)

	_, err := h.svc.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#Q": aws.String("quantity"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":q": {N: aws.String(strconv.Itoa(ui.Quantity))},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {S: aws.String(getCartPK(ui.CartID))},
			"sk": {S: aws.String(getItemSK(ui.ItemID))},
		},
		TableName:           aws.String(h.tableName),
		UpdateExpression:    aws.String("SET #Q = :q"),
		ConditionExpression: aws.String("attribute_exists(pk) and attribute_exists(sk)"),
	})

	if err != nil {
		log.Error().Msgf("Error adding item: %s", err.Error())
		return nil, ErrCouldNotUpdateItem
	}

	log.Info().Msgf("Quantity set to %d for item %s in cart %s", ui.Quantity,
		ui.ItemID, ui.CartID)

	return &Cart{}, nil
}

//getNewItemDynamoAttributes Returns an array of attribute values for the
//item is going to be inserted by the Put method
func getNewItemDynamoAttributes(ni *NewItemInfo) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"pk":          {S: aws.String(getCartPK(ni.CartID))},
		"sk":          {S: aws.String(getItemSK(ni.ItemID))},
		"cartId":      {S: aws.String(ni.CartID)},
		"itemId":      {S: aws.String(ni.ItemID)},
		"description": {S: aws.String(ni.Description)},
		"price":       {N: aws.String(fmt.Sprintf("%f", ni.Price))},
		"quantity":    {N: aws.String(strconv.Itoa(ni.Quantity))},
	}
}

//getCartPK returns the shopping cartID formatted for the primary key column
//in the database
func getCartPK(cartID string) string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixCart, cartID)
}

//getItemSK returns the itemID formatted for the sort key column in the database
func getItemSK(itemID string) string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixItem, itemID)
}
