package main

import (
	"context"
	"os"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConfigPlatformRoutes(v1 *gin.RouterGroup) {
	v1.POST("/platform", createPlatform)
}

func createPlatform(c *gin.Context) {
	service, err := createPlatformService()

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	aux := &struct {
		Name     string            `json:"name"`
		Url      string            `json:"url"`
		Rest     string            `json:"rest"`
		Property map[string]string `json:"property"`
	}{}

	err = c.ShouldBind(aux)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	res, err := service.CreateResource(aux.Name, aux.Url, aux.Rest, aux.Property)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, res)
}

func createPlatformService() (*application.PlatformService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	repo := infra.NewPlatformRepository(client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	return application.NewPlatformService(unitOfWork, repo), nil
}
