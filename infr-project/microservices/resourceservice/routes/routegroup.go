package routes

import (
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/futugyou/resourceservice/docs"

	v1 "github.com/futugyou/resourceservice/routes/v1"
)

func NewGinRoute() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Resource Swagger Doc"
	docs.SwaggerInfo.Version = "v2.0.0"

	v1api := router.Group("/api/v1")
	{
		// resource routes
		v1.ConfigResourceRoutes(v1api)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Token, Content-Length, Access-Control-Allow-Headers, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}
