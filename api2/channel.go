package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Channel handles all channel related operations that the API can use
type Channel struct {
	Handler  http.HandlerFunc
	database DatabaseController
}

// Register initializes the given router with user related routes returning the sub router.
func (c *Channel) Register(router *mux.Router) {

	/*
	 *PUT    /channel                                        -- Create a channel
	 */
	sub := router.NewRoute().PathPrefix("/channel").Subrouter()
	sub.Path("/").Handler(c.setHandler(c.CreateChannel)).Methods("PUT")

	/*
	 *PUT    /channel/{channel_id}/users 					 -- Add one or more users to a channel
	 *DELETE /channel/{channel_id}/users                     -- Delete one or more users in a channel
	 *GET    /channel/{channel_id}/users					 -- Get all the users in a channel
	 *GET    /channel/{channel_id}/messages?limit=N&offset=M -- Get message in a channel
	 */
	sub = sub.PathPrefix(fmt.Sprintf("/{cid:%s}", UUIDPattern)).Subrouter()
	sub.Path("/users").Handler(c.setHandler(c.AddUsers)).Methods("PUT")
	sub.Path("/users").Handler(c.setHandler(c.ListUsers)).Methods("GET")
	sub.Path("/users").Handler(c.setHandler(c.DeleteUsers)).Methods("DELETE")
	sub.Path("/messages").Handler(c.setHandler(c.Messages)).Methods("GET")
}

func (c *Channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler(w, r)
}

// setHandler makes a copy of the controller and sets the handler to the given handler
func (c Channel) setHandler(h http.HandlerFunc) http.Handler {
	n := c
	n.Handler = h
	return &n
}

type addChannelPayload struct {
	Owner   string   `json:"owner"   validate:"required,uuid"`
	Name    string   `json:"name"    validate:"required"`
	Members []string `json:"members" validate:"required,gt=0,dive,uuid"`
}

// CreateChannel makes a new channel
func (c *Channel) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var (
		payload = &addChannelPayload{}
		rw      = w.(*ResponseWriter)
	)

	if err := ValidateBody(payload, r.Body); err != nil {
		rw.JSON(err)
		return
	}
	var members = convertUsers(payload.Members)

	var info = &ChannelInfo{
		Owner:   payload.Owner,
		Name:    payload.Name,
		Members: members,
	}

	if err := c.database.CreateChannel(info); err != nil {
		rw.JSON(err)
		return
	}
	rw.JSON(info, http.StatusCreated)
}

// Messages get the messages for a given channel
func (c *Channel) Messages(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

type addUsersPayload struct {
	ID      string   `json:"id"      validate:"required,uuid"`
	Owner   string   `json:"owner"   validate:"required,uuid"`
	Members []string `json:"members" validate:"required,gt=0,dive,uuid"`
}

// AddUsers will add a user to the channel members set
func (c *Channel) AddUsers(w http.ResponseWriter, r *http.Request) {
	var (
		payload = &addUsersPayload{}
		rw      = w.(*ResponseWriter)
	)

	if err := ValidateBody(payload, r.Body); err != nil {
		rw.JSON(err)
		return
	}

	var members = convertUsers(payload.Members)

	var info = &ChannelInfo{
		ID:      payload.ID,
		Owner:   payload.Owner,
		Members: members,
	}

	if err := c.database.AddUsersToChannel(info); err != nil {
		rw.JSON(err)
		return
	}
	rw.JSON("OK")
}

type listUsersPayload struct {
	ID string `json:"id" validate:"required,uuid"`
}

// ListUsers lists all users in a channel
func (c *Channel) ListUsers(w http.ResponseWriter, r *http.Request) {

	var (
		payload = &listUsersPayload{}
		rw      = w.(*ResponseWriter)
	)
	if err := ValidateBody(payload, r.Body); err != nil {
		rw.JSON(err)
		return
	}

	var info = &ChannelInfo{
		ID: payload.ID,
	}

	if err := c.database.ListUsersInChannel(info); err != nil {
		rw.JSON(err)
		return
	}
	rw.JSON(info)
}

type deleteUsersPayload struct {
	ID      string   `json:"id"      validate:"required,uuid"`
	Owner   string   `json:"owner"   validate:"required,uuid"`
	Members []string `json:"members" validate:"required,gt=0,dive,uuid"`
}

// DeleteUsers removes a user from a given channel.  Only the owner of the channel may remove a user
func (c *Channel) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	var (
		payload = &deleteUsersPayload{}
		rw      = w.(*ResponseWriter)
	)
	if err := ValidateBody(payload, r.Body); err != nil {
		rw.JSON(err)
		return
	}
	var members = convertUsers(payload.Members)
	var info = &ChannelInfo{
		ID:      payload.ID,
		Owner:   payload.Owner,
		Members: members,
	}

	if err := c.database.DeleteUsersFromChannel(info); err != nil {
		rw.JSON(err)
		return
	}
	rw.JSON("OK")
}

func convertUsers(users []string) []*UserInfo {
	var members = make([]*UserInfo, 0, len(users))
	for _, member := range users {
		members = append(members, &UserInfo{
			ID: member,
		})
	}
	return members
}

func convertChannels(channels []string) []*ChannelInfo {
	var chs = make([]*ChannelInfo, 0, len(channels))
	for _, channel := range channels {
		chs = append(chs, &ChannelInfo{
			ID: channel,
		})
	}
	return chs
}
