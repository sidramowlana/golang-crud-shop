# Golang CRUD Operations
This is a basic example of CRUD (Create, Read, Update, Delete) operations in GoLang using the net/http package.

# Prerequisites
GoLang version 1.23.0

# Installation
No external dependencies are required for this project.

# Project Structure
The project structure is as follows:

main.go
models
product.go
go.mod
go.sum
main.go

# Endpoints
The following endpoints are available:

1. GET /: Home page
2. GET /all-products: Get all products
3. GET /product/{id}: Get a product by ID
4. POST /create-product: Create a new product
5. PUT /update-product/{id}: Update a product by ID
6. DELETE /delete-product/{id}: Delete a product by ID

# Usage
Run the application using go run main.go
Use a tool like curl or a REST client to test the endpoints

# Example Use Cases
1. Create a new product: curl -X POST -H "Content-Type: application/json" -d '{"name":"New Product","description":"This is a new product","price":10.99,"quantity":5}' http://localhost:8081/create-product
2. Get all products: curl http://localhost:8081/all-products
3. Get a product by ID: curl http://localhost:8081/product/1
4. Update a product by ID: curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Product","description":"This is an updated product","price":11.99,"quantity":10}' http://localhost:8081/update-product/1
5. Delete a product by ID: curl -X DELETE http://localhost:8081/delete-product/1

Note: This example uses a simple in-memory data store for demonstration purposes. In a real-world application, you would want to use a more robust data storage solution.