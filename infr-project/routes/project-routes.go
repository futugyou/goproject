package routes

import (
	"github.com/gin-gonic/gin"

	_ "github.com/futugyou/infr-project/resource"
	_ "github.com/futugyou/infr-project/view_models"

	apiadapter "github.com/futugyou/infr-project/api_adapter"
)

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
// @Param request body []viewmodels.UpdateProjectDesignRequest true "Request body"
// @Success 200
// @Router /project/{id}/design [put]
func updateProjectDesign(c *gin.Context) {
	apiadapter.UpdateProjectDesign(c.Param("id"), c.Writer, c.Request)
}

// @Summary update project platform
// @Description update project platform
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body []viewmodels.UpdateProjectPlatformRequest true "Request body"
// @Success 200
// @Router /project/{id}/platform [put]
func updateProjectPlatform(c *gin.Context) {
	apiadapter.UpdateProjectPlatform(c.Param("id"), c.Writer, c.Request)
}

// @Summary update project
// @Description update project
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body viewmodels.UpdateProjectRequest true "Request body"
// @Success 200
// @Router /project/{id} [put]
func updateProject(c *gin.Context) {
	apiadapter.UpdateProject(c.Param("id"), c.Writer, c.Request)
}

// @Summary create project
// @Description create project
// @Tags Project
// @Accept json
// @Produce json
// @Param request body viewmodels.CreateProjectRequest true "Request body"
// @Success 200
// @Router /project [post]
func createProject(c *gin.Context) {
	apiadapter.CreateProject(c.Writer, c.Request)
}

// @Summary get all project
// @Description get all project
// @Tags Project
// @Accept json
// @Produce json
// @Success 200 {array} project.Project
// @Router /project [get]
func getAllProject(c *gin.Context) {
	apiadapter.GetAllProject(c.Writer, c.Request)
}

// @Summary get project
// @Description get project
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} project.Project
// @Router /project/{id} [get]
func getProject(c *gin.Context) {
	apiadapter.GetProject(c.Param("id"), c.Writer, c.Request)
}
