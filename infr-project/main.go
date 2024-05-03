package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-gonic/gin"
)

//go:generate go install github.com/joho/godotenv/cmd/godotenv@latest
//go:generate godotenv -f ./.env go run ../tour/main.go mongo generate
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
