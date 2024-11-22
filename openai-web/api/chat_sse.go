package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/futugyou/extensions"
	openai "github.com/futugyousuzu/go-openai"
	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Chatsse(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var buf []byte
	buf, _ = io.ReadAll(r.Body)
	chatService := services.NewChatService(createOpenAICLient())
	var request services.CreateChatRequest
	json.Unmarshal(buf, &request)

	result := chatService.CreateChatSSE(request)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set(`Content-Type`, `text/event-stream;charset-utf-8`)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for response := range result {
		if response.Messages == nil || len(response.Messages) == 0 {
			continue
		}

		message := ""
		for i := 0; i < len(response.Messages); i++ {
			content := response.Messages[i].Content
			if len(content) > 0 {
				message += content
			}
		}

		message = url.QueryEscape(message)
		data := fmt.Sprintf("data: %s\n\n", message)
		w.Write([]byte(data))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	w.Write([]byte("data: [DONE]\n\n"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func createOpenAICLient() *openai.OpenaiClient {
	openaikey := os.Getenv("openaikey")
	openaiurl := os.Getenv("openaiurl")
	client := openai.NewClient(openaikey)
	if len(openaiurl) > 0 {
		client.SetBaseUrl(openaiurl)
	}

	return client
}

func createRedisICLient() *redis.Client {
	client, _ := services.RedisClient(os.Getenv("REDIS_URL"))
	return client
}

func createMongoDbCLient() *mongo.Client {
	uri := os.Getenv("mongodb_url")
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	return client
}
