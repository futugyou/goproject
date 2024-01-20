package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pusher/pusher-http-go/v5"
)

var config *PusherConfig

func main() {
	config = NewPusherConfig()
	mux := &MyMux{}
	http.ListenAndServe(":9090", mux)
}

type MyMux struct{}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/push" {
		sayhelloName(w, r)
		return
	}
	http.NotFound(w, r)
}

type Data struct {
	X0    interface{} `json:"x0"`
	X1    interface{} `json:"x1"`
	Y0    interface{} `json:"y0"`
	Y1    interface{} `json:"y1"`
	Color string      `json:"color"`
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Token, Content-Length, Access-Control-Allow-Headers, Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}

	pusherClient := pusher.Client{
		AppID:   config.APP_ID,
		Key:     config.APP_KEY,
		Secret:  config.APP_SECRET,
		Cluster: config.APP_CLUSTER,
		Secure:  true,
	}

	var d Data
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"x0":    d.X0,
		"x1":    d.X1,
		"y0":    d.Y0,
		"y1":    d.Y1,
		"color": d.Color,
	}

	err = pusherClient.Trigger("my-channel", "my-event", data)
	if err != nil {
		fmt.Println(err.Error())
	}
}
