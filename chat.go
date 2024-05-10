package main

import (
	"fmt"
)

type Chatroom struct {
	clients  map[*Client]bool
	messages chan string
}

func NewChatroom() Chatroom {
	ch := Chatroom{
		make(map[*Client]bool),
		make(chan string),
	}
	go ch.run()
	return ch
}

func (c Chatroom) String() string {
	return fmt.Sprintf("-Chatroom- %d clients: %v", len(c.clients), c.clients)
}

func (c *Chatroom) AddClient(client *Client) {
	c.clients[client] = true
}

func (c *Chatroom) RemoveClient(client *Client) {
	delete(c.clients, client)
}

func (c *Chatroom) run() {
	defer fmt.Println("chatroom closed")
	for {
		select {
		case message := <-c.messages:
			for client := range c.clients {
				select {
				case client.inboxChannel <- message:
					fmt.Println("sent message to client", client)
				default:
					fmt.Println("client is offline", client)
					c.RemoveClient(client)
					client.socket.Close()
				}
			}
		}
	}
}
