package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/goproject/blog-service/pkg/tracer"

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

var (
	port    string
	runMode string
	config  string
)

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
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s. listenandserve err: %v", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server exiting")
}

func init() {
	setupFlag()
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
	err = setupTracing()
	if err != nil {
		log.Fatalf("init tracing error : %v", err)
	}

	//global.Logger.Infof("%s  go-project-demo/%s", "test", "blog-server")
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "port")
	flag.StringVar(&runMode, "mode", "", "mode")
	flag.StringVar(&config, "config", "configs/", "config path")
	return nil
}

func setupTracing() error {
	tacer, _, err := tracer.NewJaegerTracer("blog_service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = tacer
	return nil
}

func setupSetting() error {
	setting, err := setting.NewSetting(strings.Split(config, ",")...)
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
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunModel = runMode
	}
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
