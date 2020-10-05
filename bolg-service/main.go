package main

import (
	"log"
	"net/http"
	"time"

	"github.com/goproject/blog-service/pkg/logger"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/goproject/blog-service/internal/model"

	"github.com/gin-gonic/gin"

	"github.com/goproject/blog-service/global"

	"github.com/goproject/blog-service/pkg/setting"

	"github.com/goproject/blog-service/internal/routers"
)

//go mod init github.com/goproject/blog-service
//go get -u github.com/gin-gonic/gin

// @title blog_service
// @verson 0.1
// @description go project
func main() {
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })
	// r.Run()

	gin.SetMode(global.ServerSetting.RunModel)

	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init setting error : %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init db error : %v", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("init logger error : %v", err)
	}

	//global.Logger.Infof("%s  go-project-demo/%s", "test", "blog-server")
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	//log.Print(global.ServerSetting)
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(
		&lumberjack.Logger{
			Filename:  global.AppSetting.LogSavaPath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
			MaxSize:   600,
			MaxAge:    10,
			LocalTime: true,
		},
		"",
		log.LstdFlags,
	).WithCaller(2)
	return nil
}
