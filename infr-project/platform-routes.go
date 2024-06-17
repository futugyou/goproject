package main

import (
	"github.com/gin-gonic/gin"

	_ "github.com/futugyou/infr-project/resource"
	_ "github.com/futugyou/infr-project/view_models"

	apiadapter "github.com/futugyou/infr-project/api_adapter"
)

func ConfigPlatformRoutes(v1 *gin.RouterGroup) {
	v1.POST("/platform", createPlatform)
	v1.GET("/platform/:id", getPlatform)
	v1.GET("/platform", getAllPlatform)
	v1.PUT("/platform/:id/hook", updatePlatformHook)
	v1.PUT("/platform/:id", updatePlatform)
	v1.DELETE("/platform/:id", deletePlatform)
}

// @Summary delete platform
// @Description delete platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Success 200 {string} string "ok"
// @Router /platform/{id} [delete]
func deletePlatform(c *gin.Context) {
	apiadapter.DeletePlatform(c.Param("id"), c.Writer, c.Request)
}

// @Summary update platform
// @Description update platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param request body viewmodels.UpdatePlatformRequest true "Request body"
// @Success 200
// @Router /platform/{id} [put]
func updatePlatform(c *gin.Context) {
	apiadapter.UpdatePlatform(c.Param("id"), c.Writer, c.Request)
}

// @Summary update platform webhook
// @Description update platform webhook
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param project_id path string true "Platform Project ID"
// @Param request body viewmodels.UpdatePlatformWebhookRequest true "Request body"
// @Success 200
// @Router /platform/{id}/hook [put]
func updatePlatformHook(c *gin.Context) {
	apiadapter.UpdatePlatformHook(c.Param("id"), c.Param("project_id"), c.Writer, c.Request)
}

// @Summary create platform
// @Description create platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param request body viewmodels.CreatePlatformRequest true "Request body"
// @Success 200
// @Router /platform [post]
func createPlatform(c *gin.Context) {
	apiadapter.CreatePlatform(c.Writer, c.Request)
}

// @Summary get all platform
// @Description get all platform
// @Tags Platform
// @Accept json
// @Produce json
// @Success 200 {array}  platform.Platform
// @Router /platform [get]
func getAllPlatform(c *gin.Context) {
	apiadapter.GetAllPlatform(c.Writer, c.Request)
}

// @Summary get platform
// @Description get platform
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Success 200 {object}  platform.Platform
// @Router /platform/{id} [get]
func getPlatform(c *gin.Context) {
	apiadapter.GetPlatform(c.Param("id"), c.Writer, c.Request)
}
