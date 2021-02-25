# Store
This is a simple implementation of a shopping cart using Go, Serverless, AWS Lambda and DynamoDB.

# Considerations
The following is a list of ideas that, due to initial requirements or time constraints, can not be implemented right now but should be taken into consideration.

## Architectural
- There is no authentication
- Since the API is public, it is going to be restricted to a certain number of calls
- Caching for hot partitions in DynamoDB or queries that are frequently executed
- Docker container for frontend component (using npm start for now)
- localstack for serverless backend component (currently deploying to AWS)

## Functional
- Update quantity in shopping cart
- Store cart_id in cookie or local storage so you can revisit cart if browser's window is closed
- There is no inventory. The user can add as many products as they want
- Unavailable products should not be added to shopping carts
- Rows for shopping carts that were created but are empty, can be expired(deleted)
- Int overflow of item's Quantity
- Float overflow of item's price
- Loader
- More test coverage

# API Endpoints
There are 5 API endpoints:
- GET: /items
Retrieves the list of items by category. Right now, there is only one category.

- GET: /cart/{cartId}
Retrieves the information of a shopping cart

- POST: /cart
Creates a shopping cart in the database and adds an item.
Parameters:
 - "item_id"
 - "description"
 - "quantity"
 - "price"

- POST: /cart/{cartId}
Adds an item to an existing shopping cart
Parameters:
 - "item_id"
 - "description"
 - "quantity"
 - "price"
 If a cart_id is sent in the request, it will return an error

- PATCH: /cart/{cartId}/items/{itemId}
Updates the quantity of an item in the shopping cart
Parameters:
 - "quantity"

- DELETE: /cart/{cartId}/items/{itemId}
Deletes an item from the shopping cart

# Requirements
- go version go1.15.5
- NPM Version 15.10.0
- aws-cli/2.0.58 Python/3.7.4
- serverless framework version 2

# Installation

## Environment variables
 - STORE_AWS_DYNAMODB_TABLE_STORE: DynamoDB table name
 - STORE_AWS_REGION: AWS Region where the application is stored
 - STORE_LOG_PRETTY: Human-friendly log format [pretty]
 - STORE_LOG_LEVEL: Zerolog level [error,warn,info,debug,trace] default:info

## Environment variables for test cases
As of now, the test cases for the cart package are run against a mock of the DynamoDB client. If you want to use a real dynamodb connection, the environment configuration needs to be updated in the following file:
 - api/internal/test/environment.go

## AWS Profile

You can create a specific profile for deploying the application to AWS and save it on your credentials file:
- ~/.aws/credentials
[profile_name]
aws_access_key_id = [access_key]
aws_secret_access_key = [secret_key]

## Installing Backend component
- cd api
- make
- sls deploy --aws-profile [profile_name]
- aws dynamodb batch-write-item --request-items file://seed/itemsCatalog.json
- store the endpoints server name since we're going to need it for the React application

## Installing react application
