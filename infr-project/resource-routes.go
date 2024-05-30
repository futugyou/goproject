package main

import (
	"context"
	"os"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/resource"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConfigResourceRoutes(v1 *gin.RouterGroup) {
	v1.GET("/resource/:id", getResource)
	v1.POST("/resource", createResource)
	v1.PUT("/resource", updateResource)
	v1.DELETE("/resource", deleteResource)
	v1.GET("/resource/:id/history", getResourceHistory)
}

func getResourceHistory(c *gin.Context) {
	r, err := createResourceService()

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	res, err := r.AllVersionResource(c.Param("id"))

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, res)
}

func deleteResource(c *gin.Context) {
	r, err := createResourceService()

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	err = r.DeleteResource(c.Param("id"))

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "")
}

func updateResource(c *gin.Context) {
	service, err := createResourceService()

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	aux := &struct {
		Id   string `json:"id"`
		Data string `json:"data"`
	}{}

	err = c.ShouldBind(aux)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	err = service.UpdateResourceDate(aux.Id, aux.Data)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "")
}

// @Summary create resource
// @Description create resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param request body application.CreateResourceRequest true "Request body"
// @Success 200
// @Router /resource [post]
func createResource(c *gin.Context) {
	service, err := createResourceService()

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	var aux application.CreateResourceRequest

	if err := c.ShouldBindJSON(&aux); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resourceType := resource.GetResourceType(aux.Type)
	res, err := service.CreateResource(aux.Name, resourceType, aux.Data)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, res)
}

func getResource(c *gin.Context) {
	r, err := createResourceService()

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	res, err := r.CurrentResource(c.Param("id"))

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, res)
}

func createResourceService() (*application.ResourceService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	eventStore := infra.NewMongoEventStore[resource.IResourceEvent](client, config, "resource_events")
	snapshotStore := infra.NewMongoSnapshotStore[*resource.Resource](client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	return application.NewResourceService(eventStore, snapshotStore, unitOfWork), nil
}
