package main

import (
	"log"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

//Client is created for every websocket connection to the server
type Client struct {
	id       string
	username string
	avatar   string
	manager  *ClientManager
	send     chan *Message
	socket   *websocket.Conn
}

// NewClient returns a new client with the given manager and socket connection
func NewClient(manager *ClientManager, socket *websocket.Conn, username, avatar string) *Client {
	id := uuid.NewV4()

	client := &Client{
		id:       id.String(),
		username: username,
		avatar:   avatar,
		manager:  manager,
		send:     make(chan *Message),
		socket:   socket,
	}

	go client.read()
	go client.write()

	return client
}

func (client *Client) read() {
	defer func() {
		client.manager.unregister <- client
		client.socket.Close()
	}()

	for {
		_, data, err := client.socket.ReadMessage()
		log.Printf("r %s %s\n", client.username, data)
		if err != nil {
			break
		}
		client.manager.broadcast <- NewMessage(client, string(data))
	}
}

func (client *Client) write() {
	defer func() {
		client.socket.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			log.Printf("w %s %s\n", client.id, message)
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.socket.WriteJSON(message)
		}
	}
}
