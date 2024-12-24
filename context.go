package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Define a route with dynamic context
	router.GET("/dynamic/:text", func(c *gin.Context) {
		// Retrieve the dynamic text from the URL parameter
		text := c.Param("text")

		// Respond with the dynamic text
		c.String(200, "You entered: %s", text)
	})

	// Start the server on port 8080
	router.Run(":8080")
}
