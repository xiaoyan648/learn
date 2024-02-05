package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	v1 := router.Group("/api")
	{
		v1.POST("/pipeline/build", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Success from %s", c.Request.Host),
			})
		})
	}
	// Define a route handler for the root path
	// router.GET("/api/pipeline", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Hello, World!",
	// 	})
	// })

	// Start the server
	panic(router.Run(":8080"))
}
