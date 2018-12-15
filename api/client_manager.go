package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ClientManager manages all clients on the server
type ClientManager struct {
	connections map[*Client]bool
	broadcast   chan *Message
	register    chan *Client
	unregister  chan *Client
}

// NewClientManager creates a new ClientManager and starts the manager loop
func NewClientManager() *ClientManager {
	manager := &ClientManager{
		connections: make(map[*Client]bool),
		broadcast:   make(chan *Message, broadcastChannelBufferSize),
		register:    make(chan *Client, registerChannelBufferSize),
		unregister:  make(chan *Client, unregisterChannelBufferSize),
	}

	go manager.start()
	return manager
}

// Start will start the socket listening loop
func (manager ClientManager) start() {
	for {
		select {

		// Client joining
		case client := <-manager.register:
			log.Printf("+ %s\n", client.id)
			manager.connections[client] = true
			message := NewSystemMessage(fmt.Sprintf("%s has joined the conversation", client.name))
			manager.send(message, client)

		// Client leaving
		case client := <-manager.unregister:
			if _, ok := manager.connections[client]; ok {
				log.Printf("- %s\n", client.id)
				delete(manager.connections, client)
				message := NewSystemMessage(fmt.Sprintf("%s has left the conversation", client.name))
				manager.send(message, client)
			}

		// Broadcasting
		case message := <-manager.broadcast:
			for client := range manager.connections {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.connections, client)
				}
			}
		}
	}
}

func (manager ClientManager) send(message *Message, ignore *Client) {
	for client := range manager.connections {
		if client.id != ignore.id {
			log.Printf("  %s > %s | %s\n", ignore.id, client.id, message)
			client.send <- message
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (manager ClientManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error", err)
		return
	}

	name := r.Context().Value(ContextName)
	log.Println(name)
	picture := r.Context().Value(ContextPicture).(string)
	log.Println(picture)
	gid := r.Context().Value(ContextGID).(string)
	log.Println(gid)

	client := NewClient(&manager, conn, gid, name.(string), picture)
	manager.register <- client

	conn.WriteJSON(NewInitializeMessage(client, fmt.Sprintf("Welcome %s", client.name)))
}
