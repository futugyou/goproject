package main

import (
	_ "github.com/joho/godotenv/autoload"
)

//go:generate go install github.com/joho/godotenv/cmd/godotenv@latest
//go:generate godotenv -f ./.env go run ../tour/main.go mongo generate
func main() {
	r := NewGinRoute()

	r.Run() // listen and serve on 0.0.0.0:8080
}
