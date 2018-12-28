package main

import (
	"github.com/gorilla/websocket"
	"github.com/sir-wiggles/chat/api/structs"
)

//Client is created for every websocket connection to the server
type Client struct {
	id      string
	name    string
	picture string
	manager *ClientManager
	send    chan *structs.Message
	socket  *websocket.Conn
}

// NewClient returns a new client with the given manager and socket connection
func NewClient(socket *websocket.Conn, id, name, picture string) *Client {
	client := &Client{
		id:      id,
		name:    name,
		picture: picture,
		send:    make(chan *structs.Message),
		socket:  socket,
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
		//log.Printf("r %s %s\n", client.name, data)
		if err != nil {
			break
		}
		client.manager.broadcast <- structs.NewMessage(string(data))
	}
}

func (client *Client) write() {
	defer func() {
		client.socket.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			//log.Printf("w %s %s\n", client.id, message)
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.socket.WriteJSON(message)
		}
	}
}
