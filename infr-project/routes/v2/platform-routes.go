package v2

import (
	"github.com/gin-gonic/gin"

	"github.com/futugyou/infr-project/command"
	"github.com/futugyou/infr-project/controller"
)

var cqrsRoute *command.Router

func ConfigPlatformRoutes(v2 *gin.RouterGroup, route *command.Router) {
	cqrsRoute = route
	v2.POST("/platform", createPlatform)
}

// @Summary create platform v2
// @Description create platform v2
// @Tags Platform v2
// @Accept json
// @Produce json
// @Param request body command.CreatePlatformCommand true "Request body"
// @Success 200
// @Router /v2/platform [post]
func createPlatform(c *gin.Context) {
	ctrl := controller.NewPlatformController()
	ctrl.CreatePlatformV2(cqrsRoute, c.Writer, c.Request)
}

func ConfigPlatformRoutesWithOutCQRS(v2 *gin.RouterGroup) {
	v2.GET("/platform/:id/project/:project_id", getPlatformProjectV2)
}

// @Summary get platform project by id v2
// @Description get platform project by id v2
// @Tags Platform
// @Accept json
// @Produce json
// @Param id path string true "Platform ID or Name"
// @Param project_id path string true "Platform Project ID"
// @Success 200 {object} viewmodels.PlatformProjectV2
// @Router /v2/platform/{id}/project/{project_id} [get]
func getPlatformProjectV2(c *gin.Context) {
	ctrl := controller.NewPlatformController()
	ctrl.GetPlatformProjectV2(c.Param("id"), c.Param("project_id"), c.Writer, c.Request)
}
