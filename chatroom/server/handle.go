package server

import (
	"net/http"

	"github.com/goproject/chatroom/logic"
)

func RegisterHandle() {
	go logic.Broadcaster.Start()
	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
}
