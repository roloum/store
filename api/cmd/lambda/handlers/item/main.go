package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	saws "github.com/roloum/store/api/internal/aws"
	"github.com/roloum/store/api/internal/config"
	"github.com/roloum/store/api/internal/store/item"
	"github.com/roloum/store/api/internal/web"
	"github.com/rs/zerolog/log"
)

const (
	//PathParamItemID parameter name for the item_id
	PathParamItemID = "item_id"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest,
	dynamoDB *dynamodb.DynamoDB, cfg config.Configuration) (
	events.APIGatewayProxyResponse, error) {

	//Instantiate item API Handler
	ih, err := item.New(dynamoDB, cfg.AWS.DynamoDB.Table.Store)
	if err != nil {
		return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	log.Debug().Msgf("Executing method %s for path: %s with body: %v",
		request.HTTPMethod, request.Path, request.Body)

	switch request.HTTPMethod {
	case http.MethodGet:

		return getItems(ctx, ih)

	}

	//APIGateway would not allow the function to get to this point
	//Since all the supported http methods are in the switch
	return web.GetResponse(ctx, struct{}{}, http.StatusMethodNotAllowed)

}

//getItems Returns the list of items
func getItems(ctx context.Context, ih *item.Handler) (
	events.APIGatewayProxyResponse, error) {

	var list *item.List

	list, err := ih.List(ctx)
	if err != nil {
		return web.GetResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	return web.GetResponse(ctx, list, http.StatusOK)
}

//initHandler is the function invoked by lambda that sets up the Configuration
//for the real Handler. This allows for the implementation of test cases
//for the Handler function
func initHandler(ctx context.Context, request events.APIGatewayProxyRequest) (
	events.APIGatewayProxyResponse, error) {

	//Config holds the configuration for the application
	var cfg config.Configuration
	err := config.Load(&cfg)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	sess, err := saws.GetSession(cfg.AWS.Region)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return Handler(ctx, request, saws.GetDynamoDB(sess), cfg)

}

func main() {
	lambda.Start(initHandler)
}
