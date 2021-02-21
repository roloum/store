package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	saws "github.com/roloum/store/api/internal/aws"
	"github.com/roloum/store/api/internal/config"
)

const (
	//MsgProductAdded message returned product is added to cart
	MsgProductAdded = "ProductAdded"
)

type (

	// activateResponse
	activateResponse struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
	}

	// Response is of type APIGatewayProxyResponse since we're leveraging the
	// AWS Lambda Proxy Request functionality (default behavior)
	//
	// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
	Response events.APIGatewayProxyResponse

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
	cfg configuration) (Response, error) {

	//Missing parameters? http.StatusBadRequest

	var response int
	var message string

	switch request.HTTPMethod {
	case http.MethodPost:
		//Adds a product to the shopping cart request.PathParameters["cartId"].
		//If it's the first product, it creates the shopping cart first
		response = http.StatusCreated
		message = MsgProductAdded

	case http.MethodGet:
		//Retrieves cart for request.PathParameters["cartId"]
		response = http.StatusOK

	case http.MethodPatch:
		//Udpdates the quantity for product request.PathParameters["productId"]
		//in cartId request.PathParameters["cartId"]
		response = http.StatusNoContent

	case http.MethodDelete:
		//Deletes product request.PathParameters["productId"]
		//from cartId request.PathParameters["cartId"]
		response = http.StatusNoContent

	default:
		return getResponse(http.StatusMethodNotAllowed, "")
	}

	log.Info().Msg(message)

	return getResponse(response, message)
}

// getResponse builds an API Gateway Response
func getResponse(statusCode int, message string) (Response, error) {

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp := &activateResponse{
		StatusCode: statusCode,
		Message:    message,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	log.Debug().Msgf("status_code: %d, message: %s", resp.StatusCode, resp.Message)

	return Response{Headers: headers, Body: string(js),
		StatusCode: resp.StatusCode}, nil
}

func initHandler(ctx context.Context, request events.APIGatewayProxyRequest) (
	Response, error) {

	//Config holds the configuration for the application
	var cfg configuration
	err := config.Load(&cfg)
	if err != nil {
		return Response{}, err
	}

	log.Debug().Msg("initHandler function")

	sess, err := saws.GetSession(cfg.AWS.Region)
	if err != nil {
		return Response{}, err
	}

	return Handler(ctx, saws.GetDynamoDB(sess), request, cfg)

}

func main() {
	lambda.Start(initHandler)
}
