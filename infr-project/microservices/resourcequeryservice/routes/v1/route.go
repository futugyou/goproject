package v1

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/resourcequeryservice/application"
	"github.com/futugyou/resourcequeryservice/infrastructure"
	"github.com/futugyou/resourcequeryservice/options"
	_ "github.com/futugyou/resourcequeryservice/viewmodel"
)

func ConfigResourceRoutes(v1 *gin.RouterGroup) {
	v1.GET("/resource", getAllResource)
	v1.GET("/resource/:id", getResource)
}

// @Summary get resource
// @Description get resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {object} viewmodel.ResourceView
// @Router /v1/resource/{id} [get]
func getResource(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreateResourceQueryService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := service.GetResource(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary get all resources
// @Description get all resources
// @Tags Resource
// @Accept json
// @Produce json
// @Success 200 {array} viewmodel.ResourceView
// @Router /v1/resource [get]
func getAllResource(c *gin.Context) {
	ctx := c.Request.Context()
	service, err := CreateResourceQueryService(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	datas, err := service.GetAllResources(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, datas)
}

var (
	redisClient *redis.Client
	once        sync.Once
)

func CreateRedisClient(url string) (*redis.Client, error) {
	var err error
	once.Do(func() {
		opt, err := redis.ParseURL(url)
		if err != nil {
			return
		}

		opt.MaxRetries = 3
		opt.DialTimeout = 10 * time.Second
		opt.ReadTimeout = -1
		opt.WriteTimeout = -1
		opt.DB = 0

		redisClient = redis.NewClient(opt)
	})

	return redisClient, err
}

func CreateResourceQueryService(ctx context.Context) (*application.ResourceQueryService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.QueryMongoDBURL)
	config := mongoimpl.DBConfig{
		DBName: option.QueryDBName,
	}

	if err != nil {
		return nil, err
	}

	client, err := CreateRedisClient(option.RedisURL)
	if err != nil {
		return nil, err
	}

	queryRepo := infrastructure.NewResourceQueryRepository(mongoclient, config)

	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	return application.NewResourceQueryService(queryRepo, client, unitOfWork), nil
}
