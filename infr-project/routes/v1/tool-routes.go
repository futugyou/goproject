package v1

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-secrets/stable/2023-06-13/client/secret_service"

	"github.com/futugyou/circleci"
	"github.com/futugyou/vercel"

	"github.com/futugyou/infr-project/command"
	"github.com/futugyou/infr-project/controller"
	"github.com/futugyou/infr-project/extensions"
	sdk "github.com/futugyou/infr-project/platform_sdk"
	"github.com/futugyou/infr-project/services"
)

var cqrsRoute *command.Router

func ConfigTestingRoutes(v1 *gin.RouterGroup, route *command.Router) {
	cqrsRoute = route
	v1.GET("/test/ping", pingEndpoint)
	v1.GET("/test/workflow", workflowEndpoint)
	v1.GET("/test/vercel", vercelProjectEndpoint)
	v1.GET("/test/circleci", circleciPipeline)
	v1.GET("/test/circleci/project", circleciProject)
	v1.GET("/test/circleci/projects", circleciProjectList)
	v1.GET("/test/vault", vaultSecret)
	v1.GET("/test/tf", terraformWS)
	v1.GET("/test/cqrs", cqrstest)
	v1.GET("/test/redis", redisget)
	v1.GET("/test/redishash", redisHash)
	v1.GET("/test/webhook", webhook)
}

// @Summary ping
// @Description ping
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  map[string]string
// @Router /v1/test/ping [get]
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
// @Router /v1/test/workflow [get]
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
// @Router /v1/test/vercel [get]
func vercelProjectEndpoint(c *gin.Context) {
	ctx := c.Request.Context()
	f := vercel.NewClient(os.Getenv("VERCEL_TOKEN"))
	result, _ := f.Projects.ListProject(ctx, vercel.ListProjectParameter{})
	c.JSON(200, result)
}

// @Summary circle CI
// @Description circle CI
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {string}  string
// @Router /v1/test/circleci [get]
func circleciPipeline(c *gin.Context) {
	ctx := c.Request.Context()
	f := circleci.NewClientV2(os.Getenv("CIRCLECI_TOKEN"))
	result, _ := f.Pipeline.Pipelines(ctx, os.Getenv("CIRCLECI_ORG_SLUG"))
	c.JSON(200, result)
}

// @Summary circle CI project
// @Description circle CI project
// @Tags Test
// @Accept json
// @Produce json
// @Param org_slug query string  true  "org_slug"
// @Param name query string  true  "name"
// @Success 200 {string}  string
// @Router /v1/test/circleci/project [get]
func circleciProject(c *gin.Context) {
	ctx := c.Request.Context()
	org_slug := c.Query("org_slug")
	name := c.Query("name")
	f := circleci.NewClientV2(os.Getenv("CIRCLECI_TOKEN"))
	result, _ := f.Project.GetProject(ctx, org_slug, name)
	c.JSON(200, result)
}

// @Summary circle CI project List
// @Description circle CI project List
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {string}  string
// @Router /v1/test/circleci/projects [get]
func circleciProjectList(c *gin.Context) {
	ctx := c.Request.Context()
	f := circleci.NewClientV1(os.Getenv("CIRCLECI_TOKEN"))
	result, _ := f.Project.ListProject(ctx)
	c.JSON(200, result)
}

// @Summary vault
// @Description vault
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  secret_service.OpenAppSecretOK
// @Router /v1/test/vault [get]
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
// @Router /v1/test/tf [get]
func terraformWS(c *gin.Context) {
	tfclient, _ := sdk.NewTerraformClient(os.Getenv("TFC_TOKEN"), os.Getenv("TFC_APIBASEURL"), os.Getenv("TFC_ORG"), os.Getenv("TFC_WORKSPACE"))
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

// @Summary cqrstest
// @Description cqrstest
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  map[string]string
// @Router /v1/test/cqrstest [get]
func cqrstest(c *gin.Context) {
	commandBus := cqrsRoute.CommandBus

	bookRoomCmd := &command.BookRoom{
		RoomId:    "id-2000-01-01",
		GuestName: "John",
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}

	if err := commandBus.Send(c.Request.Context(), bookRoomCmd); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}

// @Summary redis
// @Description redis
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  map[string]string
// @Router /v1/test/redis [get]
func redisget(c *gin.Context) {
	client, err := extensions.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		c.JSON(500, gin.H{
			"linkMsg": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	err = client.Set(ctx, "foo", "bar", 10*time.Second).Err()
	if err != nil {
		c.JSON(500, gin.H{
			"WriteMsg": err.Error(),
		})
		return
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		c.JSON(500, gin.H{
			"ReadMsg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"ResultMsg": val,
	})
}

// @Summary redis hash
// @Description redis hash
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  map[string]string
// @Router /v1/test/redishash [get]
func redisHash(c *gin.Context) {
	ctx := c.Request.Context()

	client, err := extensions.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		c.JSON(500, gin.H{
			"ParseURL": err.Error(),
		})
		return
	}

	hashFields := []string{
		"model", "Deimos",
		"brand", "Ergonom",
		"type", "Enduro bikes",
		"price", "4972",
	}

	res1, err := client.HSet(ctx, "bike:1", hashFields).Result()
	if err != nil {
		c.JSON(500, gin.H{
			"HSet": err.Error(),
		})
		return
	}

	// redis could may not support this method
	// res2, err := client.HExpire(ctx, "bike:1", 10*time.Second, []string{"model", "brand", "type", "price"}...).Result()
	// if err != nil {
	// 	c.JSON(500, gin.H{
	// 		"HExpire": err.Error(),
	// 	})
	// 	return
	// }

	var res4a BikeInfo
	err = client.HGetAll(ctx, "bike:1").Scan(&res4a)
	if err != nil {
		c.JSON(500, gin.H{
			"HGetAll": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"ResultMsg":   res4a,
		"ResultCount": res1,
		// "HExpire":     res2,
	})
}

type BikeInfo struct {
	Model string `redis:"model" json:"model"`
	Brand string `redis:"brand" json:"brand"`
	Type  string `redis:"type" json:"type"`
	Price int    `redis:"price" json:"price"`
}

// @Summary webhook
// @Description webhook
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object}  map[string]string
// @Router /v1/test/webhook [get]
func webhook(c *gin.Context) {
	ctrl := controller.NewWebhookController()
	ctrl.VerifyTesting(c.Writer, c.Request)
}
