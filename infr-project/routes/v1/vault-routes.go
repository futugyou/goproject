package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/futugyou/infr-project/controller"
	_ "github.com/futugyou/infr-project/view_models"
)

func ConfigVaultRoutes(v1 *gin.RouterGroup) {
	v1.POST("/vault", createVault)
	v1.POST("/vault/:id/show", showVaultRawValue)
	v1.GET("/vault", getVault)
	v1.PUT("/vault/:id", updateVault)
	v1.DELETE("/vault/:id", deleteVault)
}

// @Summary create vault
// @Description create vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param request body viewmodels.CreateVaultsRequest true "Request body"
// @Success 200 {object} viewmodels.CreateVaultsResponse
// @Router /v1/vault [post]
func createVault(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.CreateVaults(c.Writer, c.Request)
}

// @Summary show vault value
// @Description show vault value
// @Tags Vault
// @Accept json
// @Produce json
// @Param id path string true "vault ID"
// @Success 200 string string
// @Router /v1/vault/{id}/show [post]
func showVaultRawValue(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.ShowVaultRawValue(c.Writer, c.Request, c.Param("id"))
}

// @Summary get vault
// @Description get vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param request body viewmodels.SearchVaultsRequest true "Request body"
// @Success 200 {array} viewmodels.VaultView
// @Router /v1/vault [get]
func getVault(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.SearchVaults(c.Writer, c.Request, nil, nil)
}

// @Summary update vault
// @Description update vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param id path string true "vault ID"
// @Param request body viewmodels.ChangeVaultRequest true "Request body"
// @Success 200 {object} viewmodels.VaultView
// @Router /v1/vault/{id} [put]
func updateVault(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.ChangeVault(c.Writer, c.Request, c.Param("id"))
}

// @Summary delete vault
// @Description delete vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param id path string true "vault ID"
// @Success 200 boolean boolean
// @Router /v1/vault/{id} [delete]
func deleteVault(c *gin.Context) {
	ctrl := controller.NewController()
	ctrl.DeleteVault(c.Writer, c.Request, c.Param("id"))
}
