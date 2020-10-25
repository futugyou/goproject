package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	isNew bool
	Token string `json:"token"`
	conn  *websocket.Conn
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
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		sendMsg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])
		sendMsg.Content = FilterSensitive(sendMsg.Content)
		// if strings.HasPrefix(sendMsg.Content, "@") {
		// 	sendMsg.To = strings.SplitN(sendMsg.Content, " ", 2)[0][1:]
		// }
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
		Broadcaster.Broadcast(sendMsg)
	}
}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

var globalUID uint32 = 0

func NewUser(conn *websocket.Conn, token string, nickname string, addr string) *User {
	user := &User{
		NickName:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 8),
		Token:          token,
		conn:           conn,
	}
	if user.Token != "" {
		uid, err := parseTokenAndValidate(token, nickname)
		if err == nil {
			user.UID = uid
		}
	}

	if user.UID == 0 {
		user.UID = int(atomic.AddUint32(&globalUID, 1))
		user.Token = genToken(user.UID, user.NickName)
		user.isNew = true
	}
	return user
}

func genToken(uid int, nickname string) string {
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	messageMAC := macSha256([]byte(message), []byte(secret))
	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func parseTokenAndValidate(token, nickname string) (int, error) {
	pos := strings.LastIndex(token, "uid")
	messageMac, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	uid := cast.ToInt(token[pos+3:])
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	ok := validateMac([]byte(message), messageMac, []byte(secret))
	if ok {
		return uid, nil
	}

	return 0, errors.New("token is illegal")
}

func validateMac(message, messageMac, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMac, expectedMAC)
}
