# Store
This is a simple implementation of a shopping cart using Go, Serverless, AWS Lambda and DynamoDB.

## Disclosure
This was the first time I ever used React. The code might not be what a React Developer would expect to see, since this was a coding exercise with limited time frame, but I'm tryin to improve it.

# Considerations
The following is a list of ideas that, due to initial requirements or time constraints, can not be implemented right now but should be taken into consideration.

## Architectural
- There is no authentication
- Use DAX or another caching mechanism for hot partitions in DynamoDB or queries that are frequently executed
- https://aws.amazon.com/blogs/database/choosing-the-right-dynamodb-partition-key/
- Docker container for frontend component (using npm start for now)
- localstack for serverless backend component (currently deploying to AWS)
- React configuration using webpack or another approach

## Functional
- Update quantity in shopping cart
- Store cart_id in cookie or local storage so you can revisit cart if the page is reloaded or the browser's window is closed
- There is no inventory. The user can add as many products as they want
- Unavailable products should not be added to shopping carts
- Rows for shopping carts that were created but are empty, can be expired(deleted)
- Int overflow of item's Quantity
- Float overflow of item's price
- Loader
- More test coverage

## Improvements
- Back button in the items list page does not work
- The very first time an item is deleted from the shopping cart, there is an error while rendering the page. The code still goes ahead and renders the page and you can continue adding and deleting products. The error is only seen in the developer console.

# Architecture
## Backend component
The backend API is implemented in Go, using:
 - AWS Lambda for the API functions
 - AWS APIGateway for accessing the functions (CORS is enabled in order to grant access for requests from dev environments)
 - DynamoDB for the database

 There are two components in the application:
 - api/internal/store/cart: contains the logic for all the functionality related to the shopping cart (create cart, add, update and delete item)
 - api/internal/store/item: for now, it contains the logic to return the list of items

 I am using the fat lambda approach, so there are two main binaries:
  - bin/cart: receives GET, POST, PATCH and DELETE requests
  - bin/item: receives GET requests

## Database design
I am using the single table design approach for DynamoDB, overloading the keys to store multiple entities.

There are 3 entities in the application:
 - Category
 - Item
 - Cart

There is a 1-N relationship between Category and Item.

There is a N-N relationship between Cart and Item.

The Item row has a GSI with CategoryID, that allow us to load items by Category. That way we can use the ItemID in the Item row as PK, so we can validate that only existing items are added to shopping carts.

## Frontend component
The frontend application is implemented using React. It requires npm to run.

# API Endpoints
There are 5 API endpoints:
- GET: /items/{categoryId}
Retrieves the list of items by category. Right now, there is only categoryId 1.

- GET: /cart/{cartId}
Retrieves the information of a shopping cart

- POST: /cart
Creates a shopping cart in the database and adds an item. Parameters:
  - "item_id"
  - "description"
  - "quantity"
  - "price"

- POST: /cart/{cartId}
Adds an item to an existing shopping cart. Parameters:
  - "item_id"
  - "description"
  - "quantity"
  - "price"
 
 If a cart_id is sent in the request, it will return an error

- PATCH: /cart/{cartId}/items/{itemId}
Updates the quantity of an item in the shopping cart. Parameters:
  - "quantity"

- DELETE: /cart/{cartId}/items/{itemId}
Deletes an item from the shopping cart

# Requirements
- go version go1.15.5
- NPM Version 15.10.0
- aws-cli/2.0.58 Python/3.7.4
- serverless framework version 2

## IAM Roles for deploying the backend component
There should be a user with necessary permissions to deploy and run the application in AWS
- IAM access
- AmazonS3 access
- AmazonAPIGatewayInvoke Access
- AmazonAPIGatewayPushToCloudWatchLogs
- AmazonDynamoDBFullAccess
- AWSCloudFormation Access
- AWSLambda Access

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
- store the endpoint server url since we're going to need it for the React application

## Installing react application
- cd web
- Update the server url in the following files, with the value from the last step in the previous section:
 - src/components/ItemsList.js:    const serverUrl = "https://changethisurl.com"
 - src/components/ItemsList.js:    const serverUrl = "https://changethisurl.com"
 - src/components/Cart.js:    const serverUrl = "https://changethisurl.com"
- npm install
- npm start

## Run application
- http://localhost:3000/

# Screenshots
## Empty shopping cart
![Empty shopping cart](https://github.com/roloum/store/blob/main/screenshots/shopping_cart_empty.jpeg?raw=true)

## Items list
![Items list](https://github.com/roloum/store/blob/main/screenshots/shopping_cart_items_list.jpeg?raw=true)

## Shopping cart
![Shopping cart](https://github.com/roloum/store/blob/main/screenshots/shopping_cart.jpeg?raw=true)
