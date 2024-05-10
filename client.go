package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clientId = 0

type Client struct {
	hub           *Chatroom
	screenName    string
	id            int
	inboxChannel  chan string
	outboxChannel chan string
	socket        *websocket.Conn
}

func (c Client) String() string {
	return fmt.Sprintf("-Client- %d: %s", c.id, c.screenName)
}

func NewClient(screenName string, hub *Chatroom, w http.ResponseWriter, r *http.Request) {
	clientId++

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error upgrading to websocket:", err)
		return
	}

	c := Client{
		hub:           hub,
		screenName:    screenName,
		id:            clientId,
		inboxChannel:  make(chan string),
		outboxChannel: make(chan string),
		socket:        socket,
	}

	hub.AddClient(&c)

	// defer func() {
	// 	fmt.Println("closing client:", c.screenName)
	// 	close(c.outboxChannel)
	// 	close(c.inboxChannel)
	// 	c.socket.Close()
	// 	c.hub.RemoveClient(&c)
	// 	fmt.Println("client closed:", c.screenName)
	// }()

	go c.sender()
	go c.receiver()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Client) sender() {
	for {
		select {
		case message, ok := <-c.inboxChannel:
			if !ok {
				fmt.Println(c.screenName, "inbox channel closed")
				return
			}
			fmt.Println(c.screenName, "received message:", message)
			c.socket.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

func (c *Client) receiver() {
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			fmt.Println("error reading message:", err)
			break
		}
		fmt.Println(c.screenName, "sent message:", string(message))
		c.hub.messages <- c.screenName + ": " + string(message)
	}
}
