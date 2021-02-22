package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

//GetResponse Returns a struct of type events.APIGatewayProxyResponse
//It receives an struct of any type, along with the status code
//Sets the headers as application/json, marshals the struct and then
//Build the APIGatewayProxyResponse struct
func GetResponse(ctx context.Context, data interface{},
	statusCode int) (events.APIGatewayProxyResponse, error) {

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	js, err := json.Marshal(data)
	if err != nil {
		log.Debug().Msgf("Error marshalling response: %v", data)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	log.Debug().Msgf("status_code: %d, message: %v", statusCode, data)

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(js),
		StatusCode: statusCode,
	}, nil

}
