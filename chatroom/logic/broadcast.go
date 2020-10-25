package logic

import (
	"expvar"
	"fmt"
	"log"

	"github.com/goproject/chatroom/global"
)

func init() {
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}

func calcMessageQueueLen() interface{} {
	fmt.Println("===len===:", len(Broadcaster.messageChannel))
	return len(Broadcaster.messageChannel)
}

type broadcaster struct {
	users           map[string]*User
	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	checkUserChannel      chan string
	checkUserCanInChannel chan bool
	requestUsersChannel   chan struct{}
	usersChannel          chan []*User
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}

var Broadcaster = &broadcaster{
	users:           make(map[string]*User),
	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.NickName] = user
			//b.sendUserList()
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			delete(b.users, user.NickName)
			user.CloseMessageChannel()
			//b.sendUserList()
		case msg := <-b.messageChannel:
			if msg.To == "" {
				for _, user := range b.users {
					if user.UID == msg.User.UID {
						continue
					}
					user.MessageChannel <- msg
				}
				OfflineProcessor.Save(msg)
			} else {
				if user, ok := b.users[msg.To]; ok {
					user.MessageChannel <- msg
				} else {
					log.Println("user:", msg.To, " not exists.")
				}
			}
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) sendUserList() {
	userList := make([]*User, 0, len(b.users))
	for _, user := range b.users {
		userList = append(userList, user)
	}
	go func() {
		if len(b.messageChannel) < global.MessageQueueLen {
			b.messageChannel <- NewUserListMessage(userList)
		} else {
			log.Println("too mach messages")
		}
	}()
}

func NewUserListMessage(userList []*User) *Message {
	return &Message{
		Users: userList,
	}
}
