package logic

import (
	"time"

	"github.com/spf13/cast"
)

const (
	MsgTypeNormal = iota
	MsgTypeWelcome
	MsgTypeUserEnter
	MsgTypeUserLeave
	MsgTypeError
)

type Message struct {
	User           *User     `json:"user"`
	Type           int       `json:"type"`
	Content        string    `json:"content"`
	MsgTime        time.Time `json:"msg_time"`
	To             string    `json:"to"`
	Ats            []string  `json:"ats"`
	Users          []*User   `json:"users"`
	ClientSendTime time.Time `json:"send_time"`
}

func NewMessage(u *User, content, clientTime string) *Message {
	message := &Message{
		User:    u,
		Content: content,
		MsgTime: time.Now(),
		Type:    MsgTypeNormal,
	}
	if clientTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}
	return message
}

func NewNoticeMessage(message string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeNormal,
		Content: message,
		MsgTime: time.Now(),
	}
}

func NewWelcomeMessage(nickname string) *Message {
	return &Message{
		User:    &User{NickName: nickname},
		Type:    MsgTypeWelcome,
		Content: "welcome " + nickname,
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(message string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: message,
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserEnter,
		Content: user.NickName + " add chatroom",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.NickName + " leave chatroom",
		MsgTime: time.Now(),
	}
}
