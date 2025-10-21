package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

type Config struct {
	BackendURL string
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		config := Config{
			BackendURL: os.Getenv("BACKEND_URL"),
		}
		c.HTML(200, "index.html", config)
	})

	router.Run(":3000")
}
