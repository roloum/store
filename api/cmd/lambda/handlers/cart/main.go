package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	saws "github.com/roloum/store/api/internal/aws"
	"github.com/roloum/store/api/internal/config"
	"github.com/roloum/store/api/internal/store/cart"
	"github.com/roloum/store/api/internal/web"
	"github.com/rs/zerolog/log"
)

const (
	//PathParamCartID parameter name for the cart_id
	PathParamCartID = "cart_id"

	//PathParamItemID parameter name for the cart_id
	PathParamItemID = "item_id"

	//ErrRequestBodyContainsCartID error returned when adding item to existing
	//cart and there is a cart_id in the body
	ErrRequestBodyContainsCartID = "RequestBodyContainsCartID"
)

var (
	empty struct{}
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, dynamoDB *dynamodb.DynamoDB,
	request events.APIGatewayProxyRequest,
	cfg config.Configuration) (events.APIGatewayProxyResponse, error) {

	//Missing parameters? http.StatusBadRequest

	//Instantiate cart API Handler
	ch, err := cart.New(dynamoDB, cfg.AWS.DynamoDB.Table.Store)
	if err != nil {
		return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	log.Debug().Msgf("Executing method %s for path: %s with body: %v",
		request.HTTPMethod, request.Path, request.Body)

	switch request.HTTPMethod {
	case http.MethodPost:
		//Adds a item to the shopping cart request.PathParameters["cart_id"].
		//If cart_id is not set, it creates the shopping cart first

		var newItem cart.NewItemInfo
		err := json.Unmarshal([]byte(request.Body), &newItem)
		if err != nil {
			log.Error().Msgf("Error unmarshalling JSON: %s", err.Error())
			return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
		}

		var shoppingCart *cart.Cart

		//If there isn't a cartID present in the request it creates the shopping cart
		cartID, ok := request.PathParameters[PathParamCartID]
		//cartID is not in the Path, new shopping cart
		if !ok {
			shoppingCart, err = ch.CreateAndAddItem(ctx, &newItem)
		} else {
			//If cart_id is set in the path and body, return error
			if newItem.CartID != "" {
				return web.GetResponse(ctx, ErrRequestBodyContainsCartID,
					http.StatusBadRequest)
			}
			//Use cart_id from path
			newItem.CartID = cartID
			shoppingCart, err = ch.AddItem(ctx, &newItem)
		}
		if err != nil {
			return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
		}

		return web.GetResponse(ctx, shoppingCart, http.StatusCreated)

	case http.MethodGet:
		//Retrieves cart for request.PathParameters["cartId"]

	case http.MethodPatch:
		//Udpdates the quantity for item request.PathParameters["item_id"]
		//in cartId request.PathParameters["cart_id"]

		//Unmarshal the request body
		var updateItem cart.UpdateItemInfo
		err := json.Unmarshal([]byte(request.Body), &updateItem)
		if err != nil {
			log.Error().Msgf("Error unmarshalling JSON: %s", err.Error())
			return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
		}

		//Add parameters to the updateItem struct
		updateItem.CartID = request.PathParameters[PathParamCartID]
		updateItem.ItemID = request.PathParameters[PathParamItemID]

		shoppingCart, err := ch.UpdateItem(ctx, &updateItem)
		if err != nil {
			return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
		}

		return web.GetResponse(ctx, shoppingCart, http.StatusCreated)

	case http.MethodDelete:
		//Deletes item request.PathParameters["itemId"]
		//from cartId request.PathParameters["cartId"]

	}

	//APIGateway would not allow the function to get to this point
	//Since all the supported http methods are in the switch
	return web.GetResponse(ctx, empty, http.StatusMethodNotAllowed)

}

func initHandler(ctx context.Context, request events.APIGatewayProxyRequest) (
	events.APIGatewayProxyResponse, error) {

	//Config holds the configuration for the application
	var cfg config.Configuration
	err := config.Load(&cfg)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Debug().Msg("initHandler function")

	sess, err := saws.GetSession(cfg.AWS.Region)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return Handler(ctx, saws.GetDynamoDB(sess), request, cfg)

}

func main() {
	lambda.Start(initHandler)
}
