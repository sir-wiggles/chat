package main

import (
	"log"
	"net/http"
)

type clientID string
type channelID string

type ChannelManager struct {
	Clients  map[clientID]*Client
	Channels map[channelID]map[clientID]*Client
}

func (c ChannelManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error", err)
		return
	}

	var (
		name    = r.Context().Value(ContextName).(string)
		picture = r.Context().Value(ContextPicture).(string)
		gid     = r.Context().Value(ContextGID).(string)
		client  = NewClient(conn, gid, name, picture)
	)

	c.Clients[clientID(gid)] = client

}
