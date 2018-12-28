package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// User handles all user related operations that the API can use
type User struct {
	Handler  http.HandlerFunc
	database DatabaseController
}

// Register initializes the given router with user related routes returning the sub router
// created to be used with other register methods
func (c *User) Register(router *mux.Router) {
	/*
	 *DELETE /user/{user_id}/{channel_id}  -- Delete a channel
	 *GET    /user/{user_id}/channels      -- Get all channels for the user
	 */
	sub := router.NewRoute().PathPrefix(fmt.Sprintf("/user/{uid:%s}", UUIDPattern)).Subrouter()

	sub.Path(fmt.Sprintf(`/{cid:%s}`, UUIDPattern)).Handler(c.setHandler(c.DeleteChannel))
	sub.Path("/channels").Handler(c.setHandler(c.ListChannels))

}

func (c User) setHandler(h http.HandlerFunc) http.Handler {
	n := c
	n.Handler = h
	return &n
}

func (c *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler(w, r)
}

type deleteChannelPayload struct {
	Owner    string   `json:"owner" validate:"required,uuid"`
	Channels []string `json:"channels" validate:"required,dive,uuid"`
}

// DeleteChannel will remove a channel.  Can only be removed by the owner
func (c *User) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	var (
		payload = &deleteChannelPayload{}
		rw      = w.(*ResponseWriter)
	)
	if err := ValidateBody(payload, r.Body); err != nil {
		rw.JSON(err)
		return
	}

	var channels = convertChannels(payload.Channels)

	var info = &ChannelInfo{
		Owner:    payload.Owner,
		Channels: channels,
	}

	if err := c.database.DeleteChannels(info); err != nil {
		rw.JSON(err)
		return
	}
	rw.JSON("OK")
}

type listChannelPayload struct {
	Owner string `json:"owner" validate:"required,uuid"`
}

// ListChannels will get all the channels for a user
func (c *User) ListChannels(w http.ResponseWriter, r *http.Request) {
	var (
		payload = &listChannelPayload{}
		rw      = w.(*ResponseWriter)
	)
	if err := ValidateBody(payload, r.Body); err != nil {
		rw.JSON(err)
		return
	}

	var info = &ChannelInfo{
		Owner: payload.Owner,
	}

	if err := c.database.ListChannels(info); err != nil {
		rw.JSON(err)
		return
	}
	rw.JSON(info)

}
