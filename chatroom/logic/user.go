package logic

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

var System = &User{}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			return err
		}
		sendMsg := NewMessage(u, receiveMsg["content"])
		sendMsg.Content = strings.TrimSpace(sendMsg.Content)
		if strings.HasPrefix(sendMsg.Content, "@") {
			sendMsg.To = strings.SplitN(sendMsg.Content, " ", 2)[0][1:]
		}
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
		Broadcaster.Broadcast(sendMsg)
	}
}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	return &User{}
}
