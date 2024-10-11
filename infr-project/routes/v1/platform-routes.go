package v1

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/futugyou/extensions"

	viewmodels "github.com/futugyou/infr-project/view_models"

	"github.com/futugyou/infr-project/controller"
)

func ConfigPlatformRoutes(v1 *gin.RouterGroup) {
	v1.GET("/platform", searchPlatforms)
	v1.GET("/platform/:id", getPlatform)
	v1.POST("/platform", createPlatform)
	v1.PUT("/platform/:id", updatePlatform)
	v1.DELETE("/platform/:id", deletePlatform)
	// platform project
	v1.POST("/platform/:id/project", createPlatformProject)
	v1.PUT("/platform/:id/project/:project_id", updatePlatformProject)
	v1.DELETE("/platform/:id/project/:project_id", deletePlatformProject)
	// platform project webhook
	v1.PUT("/platform/:id/project/:project_id/hook", updatePlatformHook)
	v1.DELETE("/platform/:id/project/:project_id/hook/:hook_name", deletePlatformHook)
}

// @Summary delete platform webhook
// @Description delete platform webhook
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param project_id path string true "Platform Project ID"
// @Param hook_name path string true "Webhook Name"
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id}/project/{project_id}/hook/{hook_name} [delete]
func deletePlatformHook(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.DeletePlatformHook(c.Param("id"), c.Param("project_id"), c.Param("hook_name"), c.Writer, c.Request)
}

// @Summary update platform webhook
// @Description update platform webhook
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param project_id path string true "Platform Project ID"
// @Param request body viewmodels.UpdatePlatformWebhookRequest true "Request body"
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id}/project/{project_id}/hook [put]
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
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id}/project/{project_id} [delete]
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
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id}/project/{project_id} [put]
func updatePlatformProject(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.CreatePlatformProject(c.Param("id"), c.Param("project_id"), c.Writer, c.Request)
}

// @Summary create platform project
// @Description create platform project
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param request body viewmodels.UpdatePlatformProjectRequest true "Request body"
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id}/project [post]
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
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id} [delete]
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
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id} [put]
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
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform [post]
func createPlatform(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.CreatePlatform(c.Writer, c.Request)
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
// @Success 200 {array} viewmodels.PlatformView
// @Router /v1/platform [get]
func searchPlatforms(c *gin.Context) {
	ctrl := controller.NewController()

	name := c.Query("name")
	tags := strings.FieldsFunc(c.Query("tags"), func(r rune) bool {
		return r == ','
	})
	if len(tags) == 1 && tags[0] == "" {
		tags = nil
	}

	pageInt, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		pageInt = 1
	}

	sizeInt, err := strconv.Atoi(c.DefaultQuery("size", "100"))
	if err != nil {
		sizeInt = 100
	}

	request := viewmodels.SearchPlatformsRequest{
		Name:     name,
		Activate: extensions.StringToBoolPtr(c.Query("activate")),
		Tags:     tags,
		Page:     pageInt,
		Size:     sizeInt,
	}
	ctrl.SearchPlatforms(c.Writer, c.Request, request)
}

// @Summary get platform by id
// @Description get platform by id
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Success 200 {object} viewmodels.PlatformDetailView
// @Router /v1/platform/{id} [get]
func getPlatform(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.GetPlatform(c.Param("id"), c.Writer, c.Request)
}
