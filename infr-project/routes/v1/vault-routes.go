package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/futugyou/infr-project/controller"
	viewmodels "github.com/futugyou/infr-project/view_models"
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
// @Param key query string false "Key"
// @Param storage_media query string false "Storage Media"
// @Param tags query []string false "Tags" collectionFormat(csv)
// @Param type_identity query string false "Type Identity"
// @Param vault_type query string false "Vault Type"
// @Success 200 {array} viewmodels.VaultView
// @Router /v1/vault [get]
func getVault(c *gin.Context) {
	key := c.Query("key")
	storageMedia := c.Query("storage_media")
	tags := c.QueryArray("tags")
	typeIdentity := c.Query("type_identity")
	vaultType := c.Query("vault_type")

	request := viewmodels.SearchVaultsRequest{
		Key:          key,
		StorageMedia: storageMedia,
		Tags:         tags,
		TypeIdentity: typeIdentity,
		VaultType:    vaultType,
	}
	ctrl := controller.NewController()
	ctrl.SearchVaults(c.Writer, c.Request, request, nil, nil)
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
