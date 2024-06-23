package v1

import (
	"github.com/gin-gonic/gin"

	_ "github.com/futugyou/infr-project/resource"
	_ "github.com/futugyou/infr-project/view_models"

	"github.com/futugyou/infr-project/controller"
)

func ConfigPlatformRoutes(v1 *gin.RouterGroup) {
	v1.GET("/platform", getAllPlatform)
	v1.GET("/platform/:id", getPlatform)
	v1.POST("/platform", createPlatform)
	v1.PUT("/platform/:id", updatePlatform)
	v1.DELETE("/platform/:id", deletePlatform)
	// platform project
	v1.POST("/platform/:id/project", createPlatformProject)
	v1.PUT("/platform/:id/project/:project_id", updatePlatformProject)
	v1.DELETE("/platform/:id/project/:project_id", deletePlatformProject)
	v1.PUT("/platform/:id/project/:project_id/hook", updatePlatformHook)
}

// @Summary update platform webhook
// @Description update platform webhook
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param project_id path string true "Platform Project ID"
// @Param request body viewmodels.UpdatePlatformWebhookRequest true "Request body"
// @Success 200 {object} platform.Platform
// @Router /platform/{id}/project/{project_id}/hook [put]
func updatePlatformHook(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.UpdatePlatformHook(c.Param("id"), c.Param("project_id"), c.Writer, c.Request)
}

// @Summary delete platform project
// @Description delete platform project
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param project_id path string true "Platform Project ID"
// @Success 200 {object} platform.Platform
// @Router /platform/{id}/project/{project_id} [delete]
func deletePlatformProject(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.DeletePlatformProject(c.Param("id"), c.Param("project_id"), c.Writer, c.Request)
}

// @Summary update platform project
// @Description update platform project
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param project_id path string true "Platform Project ID"
// @Param request body viewmodels.UpdatePlatformProjectRequest true "Request body"
// @Success 200 {object} platform.Platform
// @Router /platform/{id}/project/{project_id} [put]
func updatePlatformProject(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.CreatePlatformProject(c.Param("id"), c.Param("project_id"), c.Writer, c.Request)
}

// @Summary create platform webhook
// @Description create platform webhook
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param request body viewmodels.UpdatePlatformProjectRequest true "Request body"
// @Success 200 {object} platform.Platform
// @Router /platform/{id}/project [post]
func createPlatformProject(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.CreatePlatformProject(c.Param("id"), "", c.Writer, c.Request)
}

// @Summary delete platform
// @Description delete platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Success 200 {object} platform.Platform
// @Router /platform/{id} [delete]
func deletePlatform(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.DeletePlatform(c.Param("id"), c.Writer, c.Request)
}

// @Summary update platform
// @Description update platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param request body viewmodels.UpdatePlatformRequest true "Request body"
// @Success 200 {object} platform.Platform
// @Router /platform/{id} [put]
func updatePlatform(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.UpdatePlatform(c.Param("id"), c.Writer, c.Request)
}

// @Summary create platform
// @Description create platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param request body viewmodels.CreatePlatformRequest true "Request body"
// @Success 200 {object} platform.Platform
// @Router /platform [post]
func createPlatform(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.CreatePlatform(c.Writer, c.Request)
}

// @Summary get all platform
// @Description get all platform
// @Tags Platform
// @Accept json
// @Produce json
// @Success 200 {array} platform.Platform
// @Router /platform [get]
func getAllPlatform(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.GetAllPlatform(c.Writer, c.Request)
}

// @Summary get platform
// @Description get platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Success 200 {object} platform.Platform
// @Router /platform/{id} [get]
func getPlatform(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.GetPlatform(c.Param("id"), c.Writer, c.Request)
}
