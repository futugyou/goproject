package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/projectservice/application"
	"github.com/futugyou/projectservice/infrastructure"
	"github.com/futugyou/projectservice/options"
	"github.com/futugyou/projectservice/viewmodel"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ConfigProjectRoutes(v1 *gin.RouterGroup) {
	v1.POST("/project", createProject)
	v1.GET("/project/:id", getProject)
	v1.GET("/project", getAllProject)
	v1.PUT("/project/:id", updateProject)
	v1.PUT("/project/:id/platform", updateProjectPlatform)
	v1.PUT("/project/:id/design", updateProjectDesign)
}

// @Summary update project design
// @Description update project design
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body []viewmodel.UpdateProjectDesignRequest true "Request body"
// @Success 200 {object} Response "Successfully"
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/project/{id}/design [put]
func updateProjectDesign(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := []viewmodel.UpdateProjectDesignRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.UpdateProjectDesign(ctx, c.Param("id"), aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, Response{})
}

// @Summary update project platform
// @Description update project platform
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body []viewmodel.UpdateProjectPlatformRequest true "Request body"
// @Success 200 {object} Response "Successfully"
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/project/{id}/platform [put]
func updateProjectPlatform(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := []viewmodel.UpdateProjectPlatformRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.UpdateProjectPlatform(ctx, c.Param("id"), aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{})
}

// @Summary update project
// @Description update project
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body viewmodel.UpdateProjectRequest true "Request body"
// @Success 200 {object} Response "Successfully"
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/project/{id} [put]
func updateProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.UpdateProjectRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	err = service.UpdateProject(ctx, c.Param("id"), aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{})
}

// @Summary create project
// @Description create project
// @Tags Project
// @Accept json
// @Produce json
// @Param request body viewmodel.CreateProjectRequest true "Request body"
// @Success 200 {object} viewmodel.CreateProjectResponse "Successfully"
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/project [post]
func createProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	aux := viewmodel.CreateProjectRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	data, err := service.CreateProject(ctx, aux)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary get all project
// @Description get all project
// @Tags Project
// @Accept json
// @Produce json
// @Success 200 {array} viewmodel.ProjectView
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/project [get]
func getAllProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}
	
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var page *int
	if p, err := strconv.Atoi(pageStr); err != nil && p > 0 {
		page = &p
	}

	var size *int
	if p, err := strconv.Atoi(sizeStr); err != nil && p > 0 {
		size = &p
	}

	datas, err := service.GetAllProject(ctx, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, datas)
}

// @Summary get project
// @Description get project
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} viewmodel.ProjectView
// @Failure 400 {object} Response "Incorrect request parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /v1/project/{id} [get]
func getProject(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	datas, err := service.GetProject(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, datas)
}

func createProjectService(ctx context.Context) (*application.ProjectService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	config := mongoimpl.DBConfig{
		DBName: option.DBName,
	}

	if err != nil {
		return nil, err
	}

	repo := infrastructure.NewProjectRepository(mongoclient, config)
	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	return application.NewProjectService(unitOfWork, repo), nil
}
