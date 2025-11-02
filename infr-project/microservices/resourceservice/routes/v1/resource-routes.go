package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/domaincore/qstashdispatcherimpl"

	"github.com/futugyou/resourceservice/application"
	"github.com/futugyou/resourceservice/domain"
	"github.com/futugyou/resourceservice/options"

	"github.com/futugyou/resourceservice/viewmodel"
)

func ConfigResourceRoutes(v1 *gin.RouterGroup) {
	v1.POST("/resource", createResource)
	v1.PUT("/resource/:id", updateResource)
	v1.DELETE("/resource/:id", deleteResource)
	v1.GET("/resource/:id/history", getResourceHistory)
}

// @Summary get resource change history
// @Description get resource change history
// @Tags Resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {array} viewmodel.ResourceView
// @Router /v1/resource/{id}/history [get]
func getResourceHistory(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.AllVersionResource(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary delete resource
// @Description delete resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {string} string "ok"
// @Router /v1/resource/{id} [delete]
func deleteResource(c *gin.Context) {

	ctx := c.Request.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = service.DeleteResource(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// @Summary update resource
// @Description update resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Param request body viewmodel.UpdateResourceRequest true "Request body"
// @Success 200 {string} string "ok"
// @Router /v1/resource/{id} [put]
func updateResource(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aux := viewmodel.UpdateResourceRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = service.UpdateResource(ctx, c.Param("id"), aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// @Summary create resource
// @Description create resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param request body viewmodel.CreateResourceRequest true "Request body"
// @Success 200 {object} viewmodel.CreateResourceResponse
// @Router /v1/resource [post]
func createResource(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aux := viewmodel.CreateResourceRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.CreateResource(ctx, aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func createResourceService(ctx context.Context) (*application.ResourceService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	if err != nil {
		return nil, err
	}

	eventStoreConfig := mongoimpl.DBConfig{
		DBName:         option.DBName,
		CollectionName: "resource_events",
	}

	eventStore := mongoimpl.NewMongoEventStore(mongoclient, eventStoreConfig, domain.CreateEvent)

	snapshotStoreConfig := mongoimpl.DBConfig{
		DBName:         option.DBName,
		CollectionName: "resources",
	}
	snapshotStore := mongoimpl.NewMongoSnapshotStore[*domain.Resource](mongoclient, snapshotStoreConfig)

	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	eventPublisher := qstashdispatcherimpl.NewQStashEventDispatcher(option.QstashToken, option.QstashDestination)
	return application.NewResourceService(eventStore, snapshotStore, unitOfWork, eventPublisher), nil
}
