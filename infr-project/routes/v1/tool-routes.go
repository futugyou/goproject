package v1

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-secrets/stable/2023-06-13/client/secret_service"

	"github.com/futugyou/infr-project/sdk"
	"github.com/futugyou/infr-project/services"
)

func ConfigTestingRoutes(v1 *gin.RouterGroup) {
	v1.GET("/test/ping", pingEndpoint)
	v1.GET("/test/workflow", workflowEndpoint)
	v1.GET("/test/vercel", vercelProjectEndpoint)
	v1.GET("/test/circleci", circleciPipeline)
	v1.GET("/test/vault", vaultSecret)
	v1.GET("/test/tf", terraformWS)
}

// @Summary ping
// @Description ping
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  map[string]string
// @Router /test/ping [get]
func pingEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// @Summary workflow
// @Description workflow
// @Tags Test
// @Accept json
// @Produce json
// @Param owner query string  true  "github owner"
// @Param repo query string  true  "github repository"
// @Success 200 {object}  map[string]string
// @Router /test/workflow [get]
func workflowEndpoint(c *gin.Context) {
	owner := c.Query("owner")
	repo := c.Query("repo")

	f := services.NewWorkflowService(os.Getenv("GITHUB_TOKEN"))
	f.Workflows(owner, repo)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

// @Summary vercel
// @Description vercel
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {string}  string
// @Router /test/vercel [get]
func vercelProjectEndpoint(c *gin.Context) {
	f := sdk.NewVercelClient(os.Getenv("VERCEL_TOKEN"))
	result := f.GetProjects()
	c.JSON(200, result)
}

// @Summary circle CI
// @Description circle CI
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {string}  string
// @Router /test/circleci [get]
func circleciPipeline(c *gin.Context) {
	f := sdk.NewCircleciClient(os.Getenv("CIRCLECI_TOKEN"))
	result := f.Pipelines(os.Getenv("CIRCLECI_ORG_SLUG"))
	c.JSON(200, result)
}

// @Summary vault
// @Description vault
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  secret_service.OpenAppSecretOK
// @Router /test/vault [get]
func vaultSecret(c *gin.Context) {
	f := sdk.NewVaultClient()
	result, err := f.GetAppSecret("VERCEL_TOKEN")
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.JSON(200, result)
}

// @Summary terraform
// @Description terraform
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  map[string]string
// @Router /test/tf [get]
func terraformWS(c *gin.Context) {
	tfclient, _ := sdk.NewTerraformClient(os.Getenv("TFC_TOKEN"))
	ws, _ := tfclient.CheckWorkspace("test")
	_, err := tfclient.CreateConfigurationVersions(ws.ID, "./tmp")
	if err != nil {
		fmt.Println("cv", err.Error())
	}
	plan, err := tfclient.CreateRun(ws, true)
	if err != nil {
		fmt.Println("cr", err.Error())
	}
	err = tfclient.ApplyRun(plan.ID)
	if err != nil {
		fmt.Println("ar", err.Error())
	}
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
