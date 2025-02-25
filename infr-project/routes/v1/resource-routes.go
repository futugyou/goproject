package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/futugyou/infr-project/controller"
	_ "github.com/futugyou/infr-project/view_models"
)

func ConfigResourceRoutes(v1 *gin.RouterGroup) {
	v1.GET("/resource", getAllResource)
	v1.GET("/resource/:id", getResource)
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
// @Success 200 {array} viewmodels.ResourceView
// @Router /v1/resource/{id}/history [get]
func getResourceHistory(c *gin.Context) {
	ctrl := controller.NewResourceController()
	ctrl.GetResourceHistory(c.Param("id"), c.Writer, c.Request)
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
	ctrl := controller.NewResourceController()
	ctrl.DeleteResource(c.Param("id"), c.Writer, c.Request)
}

// @Summary update resource
// @Description update resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Param request body viewmodels.UpdateResourceRequest true "Request body"
// @Success 200 {string} string "ok"
// @Router /v1/resource/{id} [put]
func updateResource(c *gin.Context) {
	ctrl := controller.NewResourceController()
	ctrl.UpdateResource(c.Param("id"), c.Writer, c.Request)
}

// @Summary create resource
// @Description create resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param request body viewmodels.CreateResourceRequest true "Request body"
// @Success 200 {object} viewmodels.CreateResourceResponse
// @Router /v1/resource [post]
func createResource(c *gin.Context) {
	ctrl := controller.NewResourceController()
	ctrl.CreateResource(c.Writer, c.Request)
}

// @Summary get resource
// @Description get resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {object} viewmodels.ResourceView
// @Router /v1/resource/{id} [get]
func getResource(c *gin.Context) {
	ctrl := controller.NewResourceQueryController()
	ctrl.GetResource(c.Param("id"), c.Writer, c.Request)
}

// @Summary get all resources
// @Description get all resources
// @Tags Resource
// @Accept json
// @Produce json
// @Success 200 {array} viewmodels.ResourceView
// @Router /v1/resource [get]
func getAllResource(c *gin.Context) {
	ctrl := controller.NewResourceQueryController()
	ctrl.GetAllResource(c.Writer, c.Request)
}
