package aws

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//GetSession returns an AWS session
func GetSession(region string) (*session.Session, error) {

	log.Debug().Msgf("Retrieving AWS Session for region: %s", region)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return sess, nil

}

//GetDynamoDB returns an instance of the DynamoDB connection
func GetDynamoDB(sess *session.Session) *dynamodb.DynamoDB {

	//For local install, we have to specify the endpoint
	//Will try to find a cleaner solution for local deployments but this ilustrates
	//the changes for setting up the endpoint. Would use the config object to get
	//both the local stage and the endpoint.
	var dynamodbConfig *aws.Config
	if os.Getenv("STAGE") == "local" {
		dynamodbConfig.Endpoint = aws.String(os.Getenv("DYNAMODB_BACKEND"))
	}

	dynamoSvc := dynamodb.New(sess, dynamodbConfig)

	return dynamoSvc
}
