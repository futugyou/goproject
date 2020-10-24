package logic

import "log"

type broadcaster struct {
	users           map[string]*User
	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	checkUserChannel      chan string
	checkUserCanInChannel chan bool
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

var Broadcaster = &broadcaster{
	users:           make(map[string]*User),
	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.NickName] = user
			b.sendUserList()
		case user := <-b.leavingChannel:
			delete(b.users, user.NickName)
			user.CloseMessageChannel()
			b.sendUserList()
		case msg := <-b.messageChannel:
			if msg.To == "" {
				for _, user := range b.users {
					if user.UID == msg.User.UID {
						continue
					}
					user.MessageChannel <- msg
				}
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
		}
	}
}

func (b *broadcaster) sendUserList() {

}
