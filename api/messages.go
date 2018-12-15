package main

import "time"

type MessageType string

const (
	system     MessageType = "system"
	message    MessageType = "message"
	initialize MessageType = "initialize"
)

// Author holds the relevent information of the client
type Author struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// Message is the message that will be encoded/decoded when writing/reading over the socket
type Message struct {
	Author *Author     `json:"author,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Time   time.Time   `json:"time,omitempty"`
	Type   MessageType `json:"type,omitempty"`
}

// NewMessage creates a new message with the time field set to now
func NewMessage(client *Client, text string) *Message {
	return &Message{
		Author: &Author{
			ID:     client.id,
			Name:   client.name,
			Avatar: client.picture,
		},
		Text: []string{text},
		Time: time.Now(),
		Type: message,
	}
}

// NewSystemMessage creates a new system message
func NewSystemMessage(text string) *Message {
	return &Message{
		Author: &Author{
			Name:   string(system),
			Avatar: "http://localhost:5050/images/128x128/system.png",
		},
		Text: []string{text},
		Time: time.Now(),
		Type: system,
	}
}

// NewInitializeMessage creates a new system message
func NewInitializeMessage(client *Client, text string) *Message {
	return &Message{
		Author: &Author{
			ID:     client.id,
			Name:   client.name,
			Avatar: "http://localhost:5050/images/128x128/system.png",
		},
		Text: []string{text},
		Time: time.Now(),
		Type: initialize,
	}
}
