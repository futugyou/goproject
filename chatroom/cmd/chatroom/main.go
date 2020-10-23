package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	addr   = ":2022"
	banner = `
	-----------------------------------
	-
	-
	-
	-----------------------------------
	||||| doing project with golang, start on: %s
	`
)

// go get -v github.com/goproject/chatroom/cmd/chatroom
func main() {
	fmt.Printf(banner+"\n", addr)
	server.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))
}
