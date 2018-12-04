package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	images       = make([]string, 0, 10)
	names        = make([]string, 0, 10)
	clientCount  int
	greetingText = "Welcom %s!"
)

func init() {
	listing, err := ioutil.ReadDir("./images/128x128/")
	if err != nil {
		log.Println("failed to list images", err)
	}

	for _, file := range listing {

		parts := strings.Split(file.Name(), ".")
		name := parts[0]

		if name == "system" {
			continue
		}

		name = strings.Replace(name, "_", " ", -1)
		name = strings.Title(name)
		names = append(names, name)
		images = append(images, fmt.Sprintf("http://localhost:5050/images/128x128/%s", file.Name()))
	}

	rand.Shuffle(len(images), func(i, j int) {
		images[i], images[j] = images[j], images[i]
		names[i], names[j] = names[j], names[i]
	})

}

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
			message := NewSystemMessage(fmt.Sprintf("%s has joined the conversation", client.username))
			manager.send(message, client)

		// Client leaving
		case client := <-manager.unregister:
			if _, ok := manager.connections[client]; ok {
				log.Printf("- %s\n", client.id)
				delete(manager.connections, client)
				message := NewSystemMessage(fmt.Sprintf("%s has left the conversation", client.username))
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

	index := clientCount % len(names)

	// TODO: Atomic ???
	clientCount++
	avatar := images[index]
	author := names[index]

	client := NewClient(&manager, conn, author, avatar)
	manager.register <- client

	conn.WriteJSON(NewInitializeMessage(client, fmt.Sprintf(greetingText, client.username)))
}
