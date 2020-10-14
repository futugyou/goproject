package routers

import (
	"net/http"
	"time"

	"github.com/goproject/blog-service/pkg/limiter"

	"github.com/gin-gonic/gin"
	_ "github.com/goproject/blog-service/docs"
	"github.com/goproject/blog-service/global"
	"github.com/goproject/blog-service/internal/middleware"
	api "github.com/goproject/blog-service/internal/routers/api"
	v1 "github.com/goproject/blog-service/internal/routers/api/v1"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())
	r.Use(middleware.Tracing())
	r.Use(middleware.AccessLog())
	r.Use(middleware.Recovery())
	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(60 * time.Second))
	r.Use(middleware.Translations())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/auth", api.GetAuth)
	tag := v1.NewTage()
	article := v1.NewArticle()
	//upload := api.NewUpload()
	r.POST("/upload/file", api.UploadFile)
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	apiv1 := r.Group("/api/v1")
	{
		apiv1.Use(middleware.JWT())
		apiv1.POST("/tags", tag.Create)
		apiv1.GET("/tags", tag.List)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("tags/:id/state", tag.Update)

		apiv1.POST("/articles", article.Create)
		apiv1.GET("/articles", article.List)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.DELETE("/articles/:id", article.Create)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("articles/:id/state", article.Update)
	}

	return r
}
