package logic

import "time"

const (
	MsgTypeNormal = iota
	MsgTypeSystem
	MsgTypeError
	MsgTypeUserList
)

type Message struct {
	User    *User            `json:"user"`
	Type    int              `json:"type"`
	Content string           `json:"content"`
	MsgTime time.Time        `json:"msg_time"`
	To      string           `json:"to"`
	Ats     []string         `jspn:"ats"`
	Users   map[string]*User `json:"users"`
}

func NewMessage(u *User, content string) *Message {
	return &Message{
		User:    u,
		Content: content,
		MsgTime: time.Now(),
		Type:    MsgTypeNormal,
	}
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
		Type:    MsgTypeNormal,
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
