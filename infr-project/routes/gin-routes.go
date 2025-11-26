package routes

import (
	"os"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/futugyou/extensions"

	platformv1 "github.com/futugyou/platformservice/routes/v1"
	projectv1 "github.com/futugyou/projectservice/routes/v1"
	resourcequeryv1 "github.com/futugyou/resourcequeryservice/routes/v1"
	resourcev1 "github.com/futugyou/resourceservice/routes/v1"
	vaultv1 "github.com/futugyou/vaultservice/routes/v1"

	docs "github.com/futugyou/infr-project/docs"
)

func AddDefaultTokenToRequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the current request context
		reqCtx := c.Request.Context()

		// Add a custom value to the context
		updatedCtx := extensions.WithJWT(reqCtx, os.Getenv("VAULT_API_KEY"))
		// Set the updated context back to the request
		c.Request = c.Request.WithContext(updatedCtx)
		// Proceed with the request
		c.Next()
	}
}

func NewGinRoute() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())
	if gin.Mode() == gin.DebugMode {
		router.Use(AddDefaultTokenToRequestContext())
	}

	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Project Display Swagger Doc"
	docs.SwaggerInfo.Version = "v2.0.0"

	v1api := router.Group("/api/v1")
	{
		// resource routes
		resourcequeryv1.ConfigResourceRoutes(v1api)
		resourcev1.ConfigResourceRoutes(v1api)
		// platform routes
		platformv1.ConfigPlatformRoutes(v1api)
		// project routes
		projectv1.ConfigProjectRoutes(v1api)
		// vault routes
		vaultv1.ConfigVaultRoutes(v1api)
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
