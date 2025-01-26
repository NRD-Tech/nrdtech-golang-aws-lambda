package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

// Initialize Gin app
func init() {
	// Create a new Gin router
	r := gin.Default()

	// Define GET /items route
	r.GET("/items", func(c *gin.Context) {
		category := c.Query("category") // Parse query parameter
		response := fmt.Sprintf("Fetching items in category: %s", category)
		c.JSON(200, gin.H{"message": response})
	})

	// Define GET /items/:id route
	r.GET("/items/:id", func(c *gin.Context) {
		itemID := c.Param("id") // Extract path parameter
		response := fmt.Sprintf("Fetching item with ID: %s", itemID)
		c.JSON(200, gin.H{"message": response})
	})

	// Define POST /items route
	r.POST("/items", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON body"})
			return
		}

		// Extract fields from request body
		name, nameOk := requestBody["name"].(string)
		price, priceOk := requestBody["price"].(float64)

		if !nameOk || !priceOk {
			c.JSON(400, gin.H{"error": "Invalid data types in JSON body"})
			return
		}

		response := fmt.Sprintf("Created item: %s with price: %.2f", name, price)
		c.JSON(201, gin.H{"message": response})
	})

	// Wrap the Gin router with the Gin adapter
	ginLambda = ginadapter.New(r)
}

// API Gateway Handler
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Start the Lambda function
	lambda.Start(handler)
}
