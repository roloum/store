package cart

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/google/uuid"
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
)

//Cart struct is a handler for executing the actions related to the shopping cart
type Cart struct {
	svc       dynamodbiface.DynamoDBAPI
	tableName string
}

//New returns pointer to a struct of type Cart, that contains methods
//For each action that can be executed on this API
func New(svc dynamodbiface.DynamoDBAPI, tableName string) (*Cart, error) {
	if tableName == "" {
		return nil, errors.New(ErrStoreTableNameIsEmpty)
	}

	log.Debug().Msg("Returning Cart object")

	return &Cart{svc, tableName}, nil
}

//Create Creates a shopping cart and returns a pointer to a struct of type Cart
func (cart *Cart) Create(ctx context.Context) (*Info, error) {

	//Generate a unique ID for the shopping cart
	cartID := uuid.New()
	log.Debug().Msgf("Generated UUID: %s", cartID.String())

	cartInfo := Info{CartID: cartID.String()}

	_, err := cart.svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"pk": {S: aws.String(cartInfo.getCartPK())},
			"sk": {S: aws.String(cartInfo.getCartPK())},
			"id": {S: aws.String(cartInfo.CartID)},
		},
		TableName:           aws.String(cart.tableName),
		ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
	})

	if err != nil {
		return nil, ErrCreateCart
	}

	log.Info().Msgf("Cart created with ID: %s", cartInfo.CartID)

	return &cartInfo, nil
}

//AddItem Adds new item to the shopping cart
func AddItem(ctx context.Context, svc dynamodbiface.DynamoDBAPI, ni *NewItem,
	tableName string) (*Cart, error) {

	return nil, nil

}

func (c Info) getCartPK() string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixCart, c.CartID)
}

func (n NewItem) getItemSK() string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixItem, n.ItemID)
}
