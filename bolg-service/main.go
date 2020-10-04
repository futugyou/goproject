package main

import (
	"net/http"
	"time"

	"github.com/goproject/blog-service/internal/routers"
)

//go mod init github.com/goproject/blog-service
//go get -u github.com/gin-gonic/gin
func main() {
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })
	// r.Run()

	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":8090",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
