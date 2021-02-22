package aws

import (
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//GetSession returns an AWS session
func GetSession(region string) (*session.Session, error) {

	log.Debug().Msg("Retrieving AWS Session")

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
	dynamoSvc := dynamodb.New(sess)

	return dynamoSvc
}
