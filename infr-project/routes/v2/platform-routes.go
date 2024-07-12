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
	ctrl := controller.NewController()
	ctrl.CreatePlatformV2(cqrsRoute, c.Writer, c.Request)
}
