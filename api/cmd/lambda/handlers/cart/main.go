package main

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	saws "github.com/roloum/store/api/internal/aws"
	"github.com/roloum/store/api/internal/config"
	"github.com/roloum/store/api/internal/store/cart"
	"github.com/roloum/store/api/internal/web"
)

const (
	//MsgItemAdded message returned item is added to cart
	MsgItemAdded = "ItemAdded"
)

//Configuration Struct will be populated from environment variables
//Using github.com/kelseyhightower/envconfig
type (
	configuration struct {
		AWS struct {
			DynamoDB struct {
				Table struct {
					Store string `required:"true"`
				}
			}
			Region string `required:"true"`
		}
	}
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, dynamoDB *dynamodb.DynamoDB,
	request events.APIGatewayProxyRequest,
	cfg configuration) (events.APIGatewayProxyResponse, error) {

	//Missing parameters? http.StatusBadRequest

	//Instantiate cart API Handler
	c, err := cart.New(dynamoDB, cfg.AWS.DynamoDB.Table.Store)
	if err != nil {
		return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	switch request.HTTPMethod {
	case http.MethodPost:
		//Adds a item to the shopping cart request.PathParameters["cart_id"].
		//If cart_id is not set, it creates the shopping cart first

		cartID := request.QueryStringParameters["cart_id"]

		var shoppingCart *cart.Info

		//If there isn't a cartID present in the request it creates the shopping cart
		if cartID == "" {
			shoppingCart, err = c.Create(ctx)
			if err != nil {
				return web.GetResponse(ctx, err.Error(), http.StatusCreated)
			}
		}

		log.Info().Msg(MsgItemAdded)
		return web.GetResponse(ctx, shoppingCart, http.StatusCreated)

	case http.MethodGet:
		//Retrieves cart for request.PathParameters["cartId"]

	case http.MethodPatch:
		//Udpdates the quantity for item request.PathParameters["itemId"]
		//in cartId request.PathParameters["cartId"]

	case http.MethodDelete:
		//Deletes item request.PathParameters["itemId"]
		//from cartId request.PathParameters["cartId"]

	default:
		return web.GetResponse(ctx, "due to be removed", http.StatusForbidden)
	}

	return web.GetResponse(ctx, "due to be removed", http.StatusOK)
}

func initHandler(ctx context.Context, request events.APIGatewayProxyRequest) (
	events.APIGatewayProxyResponse, error) {

	//Config holds the configuration for the application
	var cfg configuration
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
