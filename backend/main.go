package main

import (
	"backend/database"
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitDB(); err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := router.Group("/api")
	{
		api.GET("/products", handlers.GetProducts)
		api.GET("/cart", handlers.GetCart)
		api.POST("/cart", handlers.AddToCart)
		api.DELETE("/cart/:id", handlers.RemoveFromCart)
		api.POST("/checkout", handlers.Checkout)
	}

	router.Run(":8080")
}
