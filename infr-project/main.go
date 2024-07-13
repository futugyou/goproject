package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/infr-project/command"
	"github.com/futugyou/infr-project/routes"
)

//go:generate go install github.com/joho/godotenv/cmd/godotenv@latest
//go:generate godotenv -f ./.env go run ../tour/main.go mongo generate
func main() {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return
	}

	ginServerStart()
}

func ginServerStart() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cqrsRoute, err := command.CreateCommandRouter(ctx)
	if err != nil {
		log.Println(err.Error())
	}
	router := routes.NewGinRoute(cqrsRoute)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
