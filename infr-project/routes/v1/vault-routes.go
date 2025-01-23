package v1

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/futugyou/infr-project/controller"
	viewmodels "github.com/futugyou/infr-project/view_models"
)

func ConfigVaultRoutes(v1 *gin.RouterGroup) {
	v1.POST("/vault/batch", createVaults)
	v1.POST("/vault", createVault)
	v1.POST("/vault/:id/show", showVaultRawValue)
	v1.GET("/vault", getVault)
	v1.PUT("/vault/:id", updateVault)
	v1.DELETE("/vault/:id", deleteVault)
	v1.POST("/import_vault", importVault)
}

// @Summary batch create vault
// @Description batch create vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param request body viewmodels.CreateVaultsRequest true "Request body"
// @Success 200 {object} viewmodels.CreateVaultsResponse
// @Router /v1/vault/batch [post]
func createVaults(c *gin.Context) {
	ctrl := controller.NewVaultController()
	ctrl.CreateVaults(c.Writer, c.Request)
}

// @Summary batch create vault
// @Description batch create vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param request body viewmodels.CreateVaultRequest true "Request body"
// @Success 200 {object} viewmodels.VaultView
// @Router /v1/vault [post]
func createVault(c *gin.Context) {
	ctrl := controller.NewVaultController()
	ctrl.CreateVault(c.Writer, c.Request)
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
	ctrl := controller.NewVaultController()
	ctrl.ShowVaultRawValue(c.Writer, c.Request, c.Param("id"))
}

// @Summary get vault
// @Description get vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param key query string false "Key - Fuzzy Search"
// @Param storage_media query string false "Storage Media"
// @Param tags query []string false "Tags" collectionFormat(csv)
// @Param type_identity query string false "Type Identity"
// @Param vault_type query string false "Vault Type"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(100)
// @Success 200 {array} viewmodels.VaultView
// @Router /v1/vault [get]
func getVault(c *gin.Context) {
	key := c.Query("key")
	storageMedia := c.Query("storage_media")
	tags := strings.FieldsFunc(c.Query("tags"), func(r rune) bool {
		return r == ','
	})
	if len(tags) == 1 && tags[0] == "" {
		tags = nil
	}
	typeIdentity := c.Query("type_identity")
	vaultType := c.Query("vault_type")

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "100")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sizeInt = 100
	}
	request := viewmodels.SearchVaultsRequest{
		Key:          key,
		StorageMedia: storageMedia,
		VaultType:    vaultType,
		TypeIdentity: typeIdentity,
		Tags:         tags,
		Page:         pageInt,
		Size:         sizeInt,
	}
	ctrl := controller.NewVaultController()
	ctrl.SearchVaults(c.Writer, c.Request, request)
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
	ctrl := controller.NewVaultController()
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
	ctrl := controller.NewVaultController()
	ctrl.DeleteVault(c.Writer, c.Request, c.Param("id"))
}

// @Summary import vault from provider
// @Description import vault from provider
// @Tags Vault
// @Accept json
// @Produce json
// @Param request body viewmodels.ImportVaultsRequest true "Request body"
// @Success 200 {object} viewmodels.ImportVaultsResponse
// @Router /v1/import_vault [post]
func importVault(c *gin.Context) {
	ctrl := controller.NewVaultController()
	ctrl.ImportVaults(c.Writer, c.Request)
}
