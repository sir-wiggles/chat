package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Channel handles all channel related operations that the API can use
type Channel struct {
	Handler  http.HandlerFunc
	database DatabaseController
}

// Register initializes the given router with user related routes returning the sub router.
func (c *Channel) Register(router *mux.Router) *mux.Router {
	sub := router.NewRoute().PathPrefix("/channel").Subrouter()

	// Create a channel
	sub.Handle("/", c.setHandler(c.Add)).Methods("PUT")

	// Gets all channels
	sub.Handle("/", c.setHandler(c.List)).Methods("GET")

	// Delete channels
	sub.Handle("/", c.setHandler(c.Delete)).Methods("DELETE")

	return sub
}

func (c *Channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler(w, r)
}

func (c *Channel) setHandler(h http.HandlerFunc) *Channel {
	u := c
	u.Handler = h
	return u
}

type addChannelPayload struct {
	Owner   string   `json:"owner"   validate:"required,uuid"`
	Name    string   `json:"name"    validate:"required"`
	Members []string `json:"members" validate:"required,gt=0,dive,uuid"`
}

// Add makes a new channel
func (c *Channel) Add(w http.ResponseWriter, r *http.Request) {
	var payload = &addChannelPayload{}
	if err := ValidateBody(payload, r.Body); err != nil {
		w.(*ResponseWriter).JSON(err)
	}

}

type listChannelPayload struct {
	Owner string `json:"owner" validate:"required,uuid"`
}

// List will get all the channels for a user
func (c *Channel) List(w http.ResponseWriter, r *http.Request) {
	var payload = &listChannelPayload{}
	if err := ValidateBody(payload, r.Body); err != nil {
		w.(*ResponseWriter).JSON(err)
	}
}

type deleteChannelPayload struct {
	ID    string `json:"id"    validate:"required,uuid"`
	Owner string `json:"owner" validate:"required,uuid"`
}

// Delete will remove a channel.  Can only be removed by the owner
func (c *Channel) Delete(w http.ResponseWriter, r *http.Request) {
	var payload = &deleteChannelPayload{}
	if err := ValidateBody(payload, r.Body); err != nil {
		w.(*ResponseWriter).JSON(err)
	}
}
