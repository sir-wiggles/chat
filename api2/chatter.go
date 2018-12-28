package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Chatter is the main chat application
type Chatter struct {
	Handler http.HandlerFunc
}

// Register initializes the given router with chatter related operations
func (c *Chatter) Register(router *mux.Router) {

	// This line will prevent middleware from being used if channter is registered first
	sub := router.NewRoute().PathPrefix("/test").Subrouter()

	sub.Path("/ws").Handler(c.setHandler(c.WebSocket)).Methods("GET")
	sub.Path("/status").Handler(c.setHandler(c.Status)).Methods("GET")
	sub.Path("/test").Handler(c.setHandler(c.Test)).Methods("POST", "GET")
}

func (c *Chatter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler(w, r)
}

func (c Chatter) setHandler(h http.HandlerFunc) http.Handler {
	n := c
	n.Handler = h
	return &n
}

// WebSocket handles upgrading the websocket connection and registering the client to chatter
func (c *Chatter) WebSocket(w http.ResponseWriter, r *http.Request) {

}

type testPayload struct {
	Members []string `json:"members" validate:"min=1,dive,uuid"`
}

// Test is to play around with ideas
func (c *Chatter) Test(w http.ResponseWriter, r *http.Request) {

	payload := &testPayload{}
	err := ValidateBody(payload, r.Body)
	if err != nil {
		log.Println("ValidationErrors", err)
		w.(*ResponseWriter).JSON(err)
		return
	}

}

// Status just returns a OK as long as the server is up
func (c *Chatter) Status(w http.ResponseWriter, r *http.Request) {
	var (
		rw = w.(*ResponseWriter)
	)
	rw.JSON("OK")
}
