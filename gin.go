package main

import (
	"github.com/gin-gonic/gin"
)

func myHandler(c *gin.Context) {
	c.String(200, "Handler for multiple paths")
}

func main() {
	r := gin.Default()

	// Register the handler for multiple paths
	r.GET("/path1", myHandler)
	r.GET("/path2", myHandler)
	r.GET("/path3", myHandler)

	r.Run(":8080")
}
