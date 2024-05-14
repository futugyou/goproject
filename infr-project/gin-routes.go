package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/resource"
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
	eventStore := infra.NewMemoryEventStore[resource.IResourceEvent]()
	snapshotStore := infra.NewMemorySnapshotStore[*resource.Resource]()
	var res *resource.Resource = &resource.Resource{}
	sourcer := application.NewEventSourcer[resource.IResourceEvent, *resource.Resource](eventStore, snapshotStore, res)

	r := application.NewResourceService(sourcer)

	res, _ = r.CreateResource("ok", resource.Excalidraw, "no data")
	log.Println(1, *res, res.DomainEvents())

	r.UpdateResourceDate(res.Id, "not ok")

	res, _ = r.CurrentResource(res.Id)
	log.Println(2, *res, res.DomainEvents())

	es, _ := sourcer.Load(res.Id)
	log.Println(2, es)

	d, _ := json.Marshal(res)
	log.Println(4, string(d))

	c.JSON(200, res)
}
