package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// User handles all user related operations that the API can use
type User struct {
	Handler http.HandlerFunc
}

// Register initializes the given router with user related routes returning the sub router
// created to be used with other register methods
func (c *User) Register(router *mux.Router) *mux.Router {
	sub := router.NewRoute().PathPrefix("/users").Subrouter()

	// Add users to a channel
	sub.Handle("/", c.setHandler(c.Add)).Methods("PUT")

	// Get all users for a channel
	sub.Handle("/", c.setHandler(c.List)).Methods("GET")

	// Delete users from a channel
	sub.Handle("/", c.setHandler(c.Delete)).Methods("DELETE")

	return sub
}

func (c *User) setHandler(h http.HandlerFunc) *User {
	u := c
	u.Handler = h
	return u
}

func (c *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler(w, r)
}

type addUsersPayload struct {
	ID      string   `json:"id"      validate:"required,uuid"`
	Owner   string   `json:"owner"   validate:"required,uuid"`
	Members []string `json:"members" validate:"required,gt=0,dive,uuid"`
}

// Add will add a user to the channel members set
func (c *User) Add(w http.ResponseWriter, r *http.Request) {
	var payload = &addUsersPayload{}
	if err := ValidateBody(payload, r.Body); err != nil {
		w.(*ResponseWriter).JSON(err)
	}
}

type listUsersPayload struct {
	ID string `json:"id" validate:"required,uuid"`
}

// List gets all the users for a given channel.
func (c *User) List(w http.ResponseWriter, r *http.Request) {
	var payload = &listUsersPayload{}
	if err := ValidateBody(payload, r.Body); err != nil {
		w.(*ResponseWriter).JSON(err)
	}
}

type deleteUsersPayload struct {
	ID      string   `json:"id"      validate:"required,uuid"`
	Owner   string   `json:"owner"   validate:"required,uuid"`
	Members []string `json:"members" validate:"required,gt=0,dive,uuid"`
}

// Delete removes a user from a given channel.  Only the owner of the channel may remove a user
func (c *User) Delete(w http.ResponseWriter, r *http.Request) {
	var payload = &deleteUsersPayload{}
	if err := ValidateBody(payload, r.Body); err != nil {
		w.(*ResponseWriter).JSON(err)
	}
}
