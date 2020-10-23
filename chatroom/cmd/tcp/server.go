package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	users := make(map[*User]struct{})
	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
			close(user.MessageChannel)
		case msg := <-messageChannel:
			for user := range users {
				if user.ID != msg.OwnerID {
					user.MessageChannel <- msg.Content
				}
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		ID:             time.Now().Second(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	go sendMessage(conn, user.MessageChannel)

	message := Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has left",
	}

	user.MessageChannel <- "welcome, " + strconv.Itoa(user.ID)

	message.Content = "user: `" + strconv.Itoa(user.ID) + "`  has enter"
	messageChannel <- message
	enteringChannel <- user

	input := bufio.NewScanner(conn)
	for input.Scan() {
		message.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- message
		userActive <- struct{}{}
	}
	if err := input.Err(); err != nil {
		log.Println("read error: ", err)
	}
	leavingChannel <- user
	message.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
	messageChannel <- message
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	messageChannel  = make(chan Message, 8)
)

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}
type Message struct {
	OwnerID int
	Content string
}
