package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/domaincore/qstashdispatcherimpl"

	"github.com/futugyou/platformservice/application"
	"github.com/futugyou/platformservice/infrastructure"
	"github.com/futugyou/platformservice/options"
	"github.com/futugyou/platformservice/viewmodel"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ConfigPlatformRoutes(v1 *gin.RouterGroup) {
	v1.GET("/platform", searchPlatforms)
	v1.GET("/platform/:id", getPlatform)
	v1.GET("/platform/:id/provider_projects", getProviderProjectList)
	v1.POST("/platform/:id/import", importProjectsFromProvider)

	v1.POST("/platform", createPlatform)
	v1.PUT("/platform/:id", updatePlatform)
	v1.DELETE("/platform/:id", deletePlatform)
	v1.POST("/platform/:id/recovery", recoveryPlatform)
	// platform project
	v1.POST("/platform/:id/project", createPlatformProject)
	v1.GET("/platform/:id/project/:project_id", getPlatformProject)
	v1.PUT("/platform/:id/project/:project_id", updatePlatformProject)
	v1.DELETE("/platform/:id/project/:project_id", deletePlatformProject)
}

// @Summary delete platform project
// @Description delete platform project
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Param project_id path string true "Platform Project ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id}/project/{project_id} [delete]
func deletePlatformProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	err = service.DeleteProject(ctx, c.Param("id"), c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

// @Summary update platform project
// @Description update platform project
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Param project_id path string true "Platform Project ID"
// @Param request body viewmodel.UpdatePlatformProjectRequest true "Request body"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id}/project/{project_id} [put]
func updatePlatformProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.UpdatePlatformProjectRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.UpsertProject(ctx, c.Param("id"), c.Param("project_id"), aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

// @Summary create platform project
// @Description create platform project
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Param request body viewmodel.UpdatePlatformProjectRequest true "Request body"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id}/project [post]
func createPlatformProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.UpdatePlatformProjectRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.UpsertProject(ctx, c.Param("id"), "", aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

// @Summary recovery platform
// @Description recovery platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id}/recovery [post]
func recoveryPlatform(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.RecoveryPlatform(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

// @Summary delete platform
// @Description delete platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id} [delete]
func deletePlatform(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.DeletePlatform(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

// @Summary update platform
// @Description update platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Param request body viewmodel.UpdatePlatformRequest true "Request body"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id} [put]
func updatePlatform(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.UpdatePlatformRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.UpdatePlatform(ctx, c.Param("id"), aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

// @Summary create platform
// @Description create platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param request body viewmodel.CreatePlatformRequest true "Request body"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform [post]
func createPlatform(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.CreatePlatformRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.CreatePlatform(ctx, aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

// @Summary search platforms
// @Description search platforms
// @Tags Platform
// @Accept json
// @Produce json
// @Param name query string false "name - Fuzzy Search"
// @Param activate query boolean false "activate"
// @Param tags query []string false "Tags" collectionFormat(csv)
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(100)
// @Success 200 {array} viewmodel.PlatformView
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform [get]
func searchPlatforms(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.SearchPlatformsRequest{}
	if err := c.ShouldBindQuery(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	datas, err := service.SearchPlatforms(ctx, aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, datas)
}

// @Summary get platform by id or name
// @Description get platform by id or name
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Success 200 {object} viewmodel.PlatformDetailView
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id} [get]
func getPlatform(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	datas, err := service.GetPlatform(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, datas)
}

// @Summary get platform provider project list by id or name
// @Description get platform provider project list by id or name
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Success 200 {array} viewmodel.PlatformProviderProject
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id}/provider_projects [get]
func getProviderProjectList(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	datas, err := service.GetProviderProjectList(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, datas)
}

// @Summary get platform project by id
// @Description get platform project by id
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Param project_id path string true "Platform Project ID"
// @Success 200 {object} viewmodel.PlatformProject
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id}/project/{project_id} [get]
func getPlatformProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	datas, err := service.GetPlatformProject(ctx, c.Param("id"), c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, datas)
}

// @Summary import platform provider projects by id or name
// @Description import platform provider projects by id or name
// @Tags Platform
// @Accept json
// @Produce json
// @Param request body viewmodel.ImportProjectsRequest true "Request body"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/platform/{id}/import [post]
func importProjectsFromProvider(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreatePlatformService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.ImportProjectsRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.ImportProjectsFromProvider(ctx, aux.PlatformID, aux.ProjectIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{})
}

func CreatePlatformService(ctx context.Context) (*application.PlatformService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	config := mongoimpl.DBConfig{
		DBName: option.DBName,
	}

	if err != nil {
		return nil, err
	}

	repo := infrastructure.NewPlatformRepository(mongoclient, config)
	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	eventPublisher := qstashdispatcherimpl.NewQStashEventDispatcher(option.QstashToken, option.QstashDestination)
	vaultService := infrastructure.NewVaultService(option)
	screenshot := infrastructure.NewScreenshot(option)
	eventhandler := infrastructure.NewEventHandler(mongoclient, config, eventPublisher)

	return application.NewPlatformService(unitOfWork, repo, eventhandler, vaultService, screenshot, option), nil
}
