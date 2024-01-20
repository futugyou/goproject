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
	if r.URL.Path == "/api/batchpush" {
		batchpush(w, r)
		return
	}
	if r.URL.Path == "/api/info" {
		info(w, r)
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

func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Token, Content-Length, Access-Control-Allow-Headers, Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

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

	attributes := "subscription_count,user_count"
	params := pusher.TriggerParams{Info: &attributes}

	channels, err := pusherClient.TriggerMultiWithParams([]string{"my-channel", "my-channel-2", "my-channel-3"}, "my-event", data, params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for k, v := range channels.Channels {
		fmt.Printf("channel name: %s", k)
		if v.SubscriptionCount != nil {
			fmt.Printf(", subscription_count: %d", *v.SubscriptionCount)
		}
		if v.UserCount != nil {
			fmt.Printf(", user_count: %d", *v.UserCount)
		}
		fmt.Println()
	}
}

func batchpush(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

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

	channel1Info := "subscription_count"
	channel2Info := "subscription_count,user_count"

	batch := []pusher.Event{
		{Channel: "my-channel-1", Name: "my-event-1", Data: data, Info: &channel1Info},
		{Channel: "presence-my-channel-2", Name: "my-event-2", Data: "hi my name is bob", Info: &channel2Info},
		{Channel: "my-channel-3", Name: "my-event-3", Data: "hi my name is alice"},
	}

	response, err := pusherClient.TriggerBatch(batch)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, attributes := range response.Batch {
		fmt.Printf("channel: %s, name: %s", batch[i].Channel, batch[i].Name)
		if attributes.SubscriptionCount != nil {
			fmt.Printf(", subscription_count: %d", *attributes.SubscriptionCount)
		}
		if attributes.UserCount != nil {
			fmt.Printf(", user_count: %d", *attributes.UserCount)
		}
		fmt.Println()
	}
}

func info(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	pusherClient := pusher.Client{
		AppID:   config.APP_ID,
		Key:     config.APP_KEY,
		Secret:  config.APP_SECRET,
		Cluster: config.APP_CLUSTER,
		Secure:  true,
	}

	prefixFilter := "humble-"
	attributes := "user_count,subscription_count"
	params := pusher.ChannelsParams{FilterByPrefix: &prefixFilter, Info: &attributes}
	channels, err := pusherClient.Channels(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for ch, _ := range channels.Channels {
		users, err := pusherClient.GetChannelUsers(ch)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, v := range users.List {
			fmt.Println(v.ID)
		}
	}
}
