package main

import (
	"log"
	"os"

	"github.com/futugyou/infr-project/sdk"
	"github.com/futugyou/infr-project/services"
	"github.com/gin-gonic/gin"
)

func NewGinRoute() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/ping", pingEndpoint)
		v1.GET("/workflow", workflowEndpoint)
		v1.GET("/vercel", vercelProjectEndpoint)
		v1.GET("/circleci", circleciPipeline)
		v1.GET("/vault", vaultSecret)
	}
	return router
}

func pingEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func workflowEndpoint(c *gin.Context) {
	owner := c.Query("owner")
	repo := c.Query("repo")

	f := services.NewWorkflowService(os.Getenv("GITHUB_TOKEN"))
	f.Workflows(owner, repo)
}

func vercelProjectEndpoint(c *gin.Context) {
	f := sdk.NewVercelClient(os.Getenv("VERCEL_TOKEN"))
	result := f.GetProjects()
	c.JSON(200, result)
}

func circleciPipeline(c *gin.Context) {
	f := sdk.NewCircleciClient(os.Getenv("CIRCLECI_TOKEN"))
	result := f.Pipelines(os.Getenv("CIRCLECI_ORG_SLUG"))
	c.JSON(200, result)
}

func vaultSecret(c *gin.Context) {
	f := sdk.NewVaultClient()
	result, err := f.GetAppSecret("VERCEL_TOKEN")
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.JSON(200, result)
}
