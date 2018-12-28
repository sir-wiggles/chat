package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Chatter is the main chat application
type Chatter struct {
	Handler http.HandlerFunc
}

// Register initializes the given router with chatter related operations
func (c *Chatter) Register(router *mux.Router) *mux.Router {
	sub := router.NewRoute().PathPrefix("/").Subrouter()

	sub.Handle("/ws", c.setHandler(c.WebSocket)).Methods("GET")

	sub.Handle("/status", c.setHandler(c.Status)).Methods("GET")

	sub.Handle("/test", c.setHandler(c.Test)).Methods("POST", "GET")

	return sub
}

func (c *Chatter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler(w, r)
}

func (c *Chatter) setHandler(h http.HandlerFunc) *Chatter {
	u := c
	u.Handler = h
	return u
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
