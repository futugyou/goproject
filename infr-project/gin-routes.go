package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/resource"
	"github.com/futugyou/infr-project/sdk"
	"github.com/futugyou/infr-project/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		v1.GET("/resource", resourceMarshal)
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

func resourceMarshal(c *gin.Context) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		log.Fatal(err)
	}

	eventStore := infra.NewMongoEventStore[resource.IResourceEvent](client, config)
	snapshotStore := infra.NewMongoSnapshotStore[*resource.Resource](client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		log.Fatal(err)
	}

	r := application.NewResourceService(eventStore, snapshotStore, unitOfWork)

	res, _ := r.CreateResource("ok", resource.Excalidraw, "no data")
	log.Println(1, res.DomainEvents())

	r.UpdateResourceDate(res.Id, "not ok")

	res, _ = r.CurrentResource(res.Id)

	d, _ := json.Marshal(res)
	log.Println(4, string(d))

	c.JSON(200, res)
}
