package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/domaincore/qstashdispatcherimpl"

	"github.com/futugyou/vaultservice/application"
	"github.com/futugyou/vaultservice/infrastructure"
	"github.com/futugyou/vaultservice/options"

	"github.com/futugyou/vaultservice/viewmodel"
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
// @Param request body viewmodel.CreateVaultsRequest true "Request body"
// @Success 200 {object} viewmodel.CreateVaultsResponse
// @Router /v1/vault/batch [post]
func createVaults(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aux := viewmodel.CreateVaultsRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.CreateVaults(ctx, aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary create vault
// @Description create vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param request body viewmodel.CreateVaultRequest true "Request body"
// @Success 200 {object} viewmodel.VaultView
// @Router /v1/vault [post]
func createVault(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aux := viewmodel.CreateVaultRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.CreateVault(ctx, aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
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
	ctx := c.Request.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.ShowVaultRawValue(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
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
// @Success 200 {array} viewmodel.VaultView
// @Router /v1/vault [get]
func getVault(c *gin.Context) {
	var req viewmodel.SearchVaultsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.SearchVaults(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary update vault
// @Description update vault
// @Tags Vault
// @Accept json
// @Produce json
// @Param id path string true "vault ID"
// @Param request body viewmodel.ChangeVaultRequest true "Request body"
// @Success 200 {object} viewmodel.VaultView
// @Router /v1/vault/{id} [put]
func updateVault(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aux := viewmodel.ChangeVaultRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.ChangeVault(ctx, c.Param("id"), aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
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
	ctx := c.Request.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.DeleteVault(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary import vault from provider
// @Description import vault from provider
// @Tags Vault
// @Accept json
// @Produce json
// @Param request body viewmodel.ImportVaultsRequest true "Request body"
// @Success 200 {object} viewmodel.ImportVaultsResponse
// @Router /v1/import_vault [post]
func importVault(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aux := viewmodel.ImportVaultsRequest{}
	if err := c.ShouldBind(&aux); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.ImportVaults(ctx, aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func createVaultService(ctx context.Context) (*application.VaultService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	config := mongoimpl.DBConfig{
		DBName: option.DBName,
	}

	if err != nil {
		return nil, err
	}

	queryRepo := infrastructure.NewVaultRepository(mongoclient, config)

	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	eventPublisher := qstashdispatcherimpl.NewQStashEventDispatcher(option.QstashToken, option.QstashDestination)
	return application.NewVaultService(unitOfWork, queryRepo, eventPublisher), nil
}
