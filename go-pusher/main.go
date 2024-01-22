package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pushnotifications "github.com/pusher/push-notifications-go"
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
	if r.URL.Path == "/api/hook" {
		hook(w, r)
		return
	}
	if r.URL.Path == "/api/auth" {
		auth(w, r)
		return
	}
	if r.URL.Path == "/api/to" {
		touser(w, r)
		return
	}
	if r.URL.Path == "/api/beams/push" {
		beamspush(w, r)
		return
	}
	if r.URL.Path == "/api/beams/pushuser" {
		beamspushuser(w, r)
		return
	}
	if r.URL.Path == "/api/beams/auth" {
		beamsauth(w, r)
		return
	}
	if r.URL.Path == "/api/beams/del" {
		beamsuserdel(w, r)
		return
	}
	http.NotFound(w, r)
}

type Data struct {
	X0       interface{} `json:"x0"`
	X1       interface{} `json:"x1"`
	Y0       interface{} `json:"y0"`
	Y1       interface{} `json:"y1"`
	Color    string      `json:"color"`
	SocketID string      `json:"socketID"`
}

type User struct {
	Userid string `json:"userid"`
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
	if d.SocketID != "" {
		params.SocketID = &d.SocketID
	}

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

	for ch := range channels.Channels {
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

func hook(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	body, _ := io.ReadAll(r.Body)
	pusherClient := pusher.Client{
		AppID:   config.APP_ID,
		Key:     config.APP_KEY,
		Secret:  config.APP_SECRET,
		Cluster: config.APP_CLUSTER,
		Secure:  true,
	}
	webhook, err := pusherClient.Webhook(r.Header, body)

	if err == nil {
		for _, event := range webhook.Events {
			switch event.Name {
			case "channel_occupied":
				fmt.Println("Channel occupied: " + event.Channel)
			case "channel_vacated":
				fmt.Println("Channel vacated: " + event.Channel)
			}
		}
	}

	fmt.Fprintf(w, "ok")
}

func auth(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	pusherClient := pusher.Client{
		AppID:   config.APP_ID,
		Key:     config.APP_KEY,
		Secret:  config.APP_SECRET,
		Cluster: config.APP_CLUSTER,
		Secure:  true,
	}
	params, _ := io.ReadAll(r.Body)
	presenceData := pusher.MemberData{
		UserID: "1",
		UserInfo: map[string]string{
			"twitter": "pusher",
		},
	}

	// This authenticates every user. Don't do this in production!
	// private channel
	// response, err := pusherClient.AuthorizePrivateChannel(params)
	// presence channel
	response, err := pusherClient.AuthorizePresenceChannel(params, presenceData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Fprint(w, string(response))
}

func touser(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	pusherClient := pusher.Client{
		AppID:   config.APP_ID,
		Key:     config.APP_KEY,
		Secret:  config.APP_SECRET,
		Cluster: config.APP_CLUSTER,
		Secure:  true,
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{"hello": "world"}
	err = pusherClient.SendToUser(u.Userid, "my-event", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func beamspush(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	beamsClient, err := pushnotifications.New(config.INSTANCE_ID, config.SECRET_KEY)
	if err != nil {
		fmt.Println("Could not create Beams Client:", err.Error())
		return
	}

	publishRequest := map[string]interface{}{
		"apns": map[string]interface{}{
			"aps": map[string]interface{}{
				"alert": map[string]interface{}{
					"title": "Hello",
					"body":  "Hello, world",
				},
			},
		},
		"fcm": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": "Hello",
				"body":  "Hello, world",
			},
		},
		"web": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": "Hello",
				"body":  "Hello, world",
			},
		},
	}

	pubId, err := beamsClient.PublishToInterests([]string{"hello"}, publishRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Publish Id:", pubId)
}

func beamspushuser(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	beamsClient, err := pushnotifications.New(config.INSTANCE_ID, config.SECRET_KEY)
	if err != nil {
		fmt.Println("Could not create Beams Client:", err.Error())
		return
	}

	publishRequest := map[string]interface{}{
		"apns": map[string]interface{}{
			"aps": map[string]interface{}{
				"alert": map[string]interface{}{
					"title": "Hello",
					"body":  "Hello, world",
				},
			},
		},
		"fcm": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": "Hello",
				"body":  "Hello, world",
			},
		},
		"web": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": "Hello",
				"body":  "Hello, world",
			},
		},
	}

	pubId, err := beamsClient.PublishToUsers([]string{"user-001", "user-002"}, publishRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Publish Id:", pubId)
}

func beamsauth(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	// do some user check

	beamsClient, err := pushnotifications.New(config.INSTANCE_ID, config.SECRET_KEY)
	if err != nil {
		fmt.Println("Could not create Beams Client:", err.Error())
		return
	}

	userid := "fake-user-id"
	beamsToken, err := beamsClient.GenerateToken(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	beamsTokenJson, err := json.Marshal(beamsToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(beamsTokenJson)
}

func beamsuserdel(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	// do some user id check

	beamsClient, err := pushnotifications.New(config.INSTANCE_ID, config.SECRET_KEY)
	if err != nil {
		fmt.Println("Could not create Beams Client:", err.Error())
		return
	}

	userid := "fake-user-id"
	err = beamsClient.DeleteUser(userid)
	if err != nil {
		fmt.Println("Could not delete user:", err.Error())
	}
}
