package server

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"

	"github.com/goproject/chatroom/logic"

	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("websocket accept error: ", err)
		return
	}
	nickname := r.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("novlid nickname, length 4~20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("nickname already existed: ", nickname)
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("nickname already existed."))
		conn.Close(websocket.StatusUnsupportedData, "nickname already existed.")
		return
	}
	token := r.FormValue("token")
	userHasToken := logic.NewUser(conn, token, nickname, r.RemoteAddr)

	go userHasToken.SendMessage(r.Context())

	userHasToken.MessageChannel <- logic.NewWelcomeMessage(nickname)

	msg := logic.NewNoticeMessage(nickname + " add chatroom")
	logic.Broadcaster.Broadcast(msg)

	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""

	msg = logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	err = user.ReceiveMessage(r.Context())

	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error :", err)
		conn.Close(websocket.StatusInternalError, "read from client error")
	}

}
