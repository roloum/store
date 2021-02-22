# Store
This is a simple implementation of a shopping cart using Go, Serverless, AWS Lambda and DynamoDB.

# Considerations
The following is a list of ideas that, due to initial requirements or time constraints, can not be implemented right now but should be taken into consideration.

## Architectural
- There is no authentication
- Since the API is public, it is going to be restricted to a certain number of calls

## Functional
- There is no inventory. The user can add as many products as they want
- Unavailable products should not be added to shopping carts
- Rows for shopping carts that were created but are empty, can be expired(deleted)
- Int overflow of item's Quantity
- Float overflow of item's price
- When trying to add an item that already exists in the shopping cart, increment the quantity by 1 instead

# API Endpoints
There are 5 API endpoints:
- GET: cart/{cartId}
Retrieves the information of a shopping cart

- POST: cart
Creates a shopping cart in the database and adds an item

- POST: cart/{cartId}
Adds an item to an existing shopping cart

- PATCH: cart/{cartId}/items/{itemId}
Updates the quantity of an item in the shopping cart

- DELETE: cart/{cartId}/items/{itemId}
Deletes an item from the shopping cart
