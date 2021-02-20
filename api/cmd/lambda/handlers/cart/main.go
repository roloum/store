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

	msg := "Successfully deployed function body"
	log.Info().Msg(msg)

	return getResponse(http.StatusOK, msg)
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
