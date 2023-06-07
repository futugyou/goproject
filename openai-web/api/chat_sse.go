package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	// "github.com/futugyousuzu/go-openai-web/oauth"
	"github.com/futugyousuzu/go-openai-web/services"
)

func Chatsse(w http.ResponseWriter, r *http.Request) {
	// oauth.AuthForVercel(w, r)

	var buf []byte
	buf, _ = io.ReadAll(r.Body)
	chatService := services.ChatService{}
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
