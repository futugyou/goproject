package routes

import (
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/futugyou/infr-project/command"
	docs "github.com/futugyou/infr-project/docs"

	v1 "github.com/futugyou/infr-project/routes/v1"
	v2 "github.com/futugyou/infr-project/routes/v2"
	projectv1 "github.com/futugyou/projectservice/routes/v1"
	vaultv1 "github.com/futugyou/vaultservice/routes/v1"
)

func NewGinRoute(cqrsRoute *command.Router) *gin.Engine {
	router := gin.Default()
	router.Use(Cors())
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Project Display Swagger Doc"
	docs.SwaggerInfo.Version = "v1.0.0"

	v2api := router.Group("/api/v2")
	{
		v2.ConfigPlatformRoutes(v2api, cqrsRoute)
	}

	v1api := router.Group("/api/v1")
	{
		// resource routes
		v1.ConfigResourceRoutes(v1api)
		// platform routes
		v1.ConfigPlatformRoutes(v1api)
		// project routes
		projectv1.ConfigProjectRoutes(v1api)
		// vault routes
		vaultv1.ConfigVaultRoutes(v1api)
		// sdk test routes
		v1.ConfigTestingRoutes(v1api, cqrsRoute)
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
