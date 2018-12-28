package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

var (
	// ErrInvalidOwner should be returned when owner is missing or invalid
	ErrInvalidOwner = errors.New(`Invalid "Owner" field in ChannelInfo`)

	// ErrInvalidName should be returned when name is missing or invalid
	ErrInvalidName = errors.New(`Invalid "Name" field in ChannelInfo`)

	// ErrInvalidMembersLen should be returned when members has invalid len
	ErrInvalidMembersLen = errors.New(`Invalid "Members" field in ChannelInfo`)

	// ErrInvalidChannel should be returned when channel is missing or invalid
	ErrInvalidChannel = errors.New(`Invalid "ID" field in ChannelInfo`)

	// ErrInvalidChannelLen should be returned when channels has invalid len
	ErrInvalidChannelLen = errors.New(`Invalid "Channels" length, must be more than one`)
)

// DatabaseController is composed of the user and channel controller interfaces
type DatabaseController interface {
	ChannelController
	UserController
}

// UserController is the user related method actions
type UserController interface {
	AddUsersToChannel(*ChannelInfo) error
	CreateUser(*UserInfo) error
	DeleteUsersFromChannel(*ChannelInfo) error
	ListUsersInChannel(*ChannelInfo) error
}

// ChannelController is the channel related method actions
type ChannelController interface {
	CreateChannel(*ChannelInfo) error
	DeleteChannels(*ChannelInfo) error
	ListChannels(*ChannelInfo) error
}

// Cassandra is the connection to cassandra
type Cassandra struct {
	*gocql.Session
}

// NewCassandra returns a new connection to cassandra using keyspace chatter
func NewCassandra(urls []string) (*Cassandra, error) {
	cluster := gocql.NewCluster(urls...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()

	return &Cassandra{
		session,
	}, err
}

// ChannelInfo is the model of channels in cassandra
type ChannelInfo struct {

	// ID is generated when creating a channel. It will have the UUID form
	ID string `json:"id,omitempty"`

	// Owner is the user that owns the channel and has the UUID form
	Owner string `json:"owner,omitempty"`

	// Created is when the channel was created and is generated when calling CreateChannel
	Created time.Time `json:"created,omitempty"`

	// Members are the users in the channel.  When createing a channel the Owner will be added
	// to the list automatically
	Members []*UserInfo `json:"members,omitempty"`

	// Name is the human readable name of the channel
	Name string `json:"name,omitempty"`

	// TODO: do we want this?
	Private bool `json:"private,omitempty"`

	// Channels is the lists of channels of a given user
	Channels []*ChannelInfo `json:"channels,omitempty"`
}

// UserInfo hold information for a particular user
type UserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	GID     string `json:"-"`
}

// CreateChannel takes ChannelInfo as input with the following fields required: "Owner" which is
// a uuid in string form, "Name" is the human readable name of the channel.
//
// When creating a channel the ID will be created and attached to the "ID" field and the "Owner"
// will be appended with the "Members" list.
func (c *Cassandra) CreateChannel(i *ChannelInfo) error {

	if i.Owner == "" {
		return ErrInvalidOwner
	}

	if strings.Trim(i.Name, " ") == "" {
		return ErrInvalidName
	}

	i.Members = append(i.Members, &UserInfo{ID: i.Owner})
	i.Created = time.Now().UTC()

	cid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	members := make([]string, 0, len(i.Members))
	for _, member := range i.Members {
		members = append(members, member.ID)
	}

	err = c.Query(`INSERT INTO channels
		(id, owner, created, members, name, private)
	VALUES
		(?, ?, now(), ?, ?, ?)`,
		cid.String(), i.Owner, members, i.Name, i.Private).Exec()

	if err != nil {
		return err
	}

	i.ID = cid.String()

	return nil
}

// ListChannels lists all the channels for a given owner.
func (c *Cassandra) ListChannels(i *ChannelInfo) error {

	if i.Owner == "" {
		return ErrInvalidOwner
	}

	iter := c.Query(
		`SELECT id, created, name FROM channels WHERE owner = ?`,
		i.Owner,
	).Iter().Scanner()

	channels := make([]*ChannelInfo, 0, 2)
	for iter.Next() {
		ci := &ChannelInfo{}
		err := iter.Scan(&ci.ID, &ci.Created, &ci.Name)
		if err != nil {
			return err
		}
		channels = append(channels, ci)
	}

	i.Channels = channels

	return nil
}

// DeleteChannels will delete the channels specified in the "Channels" field array. The "Owner"
// field must be set and the "Owner" must own the channels they're trying to delete.
func (c *Cassandra) DeleteChannels(i *ChannelInfo) error {

	if i.Owner == "" {
		return ErrInvalidOwner
	}

	ids := make([]string, 0, len(i.Channels))
	for _, ch := range i.Channels {
		ids = append(ids, ch.ID)
	}

	if len(ids) == 0 {
		return ErrInvalidChannelLen
	}

	err := c.Query(
		`DELETE FROM channels WHERE id IN ? AND owner = ?`,
		ids, i.Owner,
	).Exec()

	return err
}

// AddUsersToChannel will add users to a channel given the channel id and the owner.
// the members should be an array of uuids representing the users you want to add.
func (c *Cassandra) AddUsersToChannel(i *ChannelInfo) error {

	if i.Owner == "" {
		return ErrInvalidOwner
	} else if len(i.Members) == 0 {
		return ErrInvalidMembersLen
	} else if i.ID == "" {
		return ErrInvalidChannel
	}

	members := make([]string, 0, len(i.Members))
	for _, member := range i.Members {
		members = append(members, member.ID)
	}

	err := c.Query(
		`UPDATE channels SET members = members + ?
		WHERE id = ? AND owner = ?`,
		members, i.ID, i.Owner,
	).Exec()

	return err
}

// DeleteUsersFromChannel will remove the specified users from a channel given the users uuids
func (c *Cassandra) DeleteUsersFromChannel(i *ChannelInfo) error {

	if i.Owner == "" {
		return ErrInvalidOwner
	} else if i.ID == "" {
		return ErrInvalidChannel
	}

	members := make([]string, 0, len(i.Members))
	for _, member := range i.Members {
		members = append(members, member.ID)
	}

	err := c.Query(`
		UPDATE channels SET members = members - ?
		WHERE id = ? AND owner = ?`,
		members, i.ID, i.Owner,
	).Exec()

	return err
}

// ListUsersInChannel will list all the users in a channel
func (c *Cassandra) ListUsersInChannel(i *ChannelInfo) error {

	members := make([]string, 0, 2)
	err := c.Query(`
		SELECT members FROM channels WHERE id = ?`,
		i.ID,
	).Scan(&members)

	fmt.Println("members", members)

	if err != nil {
		return err
	}

	scanner := c.Query(`
		SELECT id, gid, name, picture FROM users WHERE id IN ?`, members,
	).Iter().Scanner()

	users := make([]*UserInfo, 0, len(members))
	for scanner.Next() {
		user := &UserInfo{}
		err = scanner.Scan(&user.ID, &user.GID, &user.Name, &user.Picture)
		if err != nil {
			return err
		}
		users = append(users, user)
	}

	log.Println("users", users)
	i.Members = users

	return nil
}

// CreateUser adds a user to the database if the user does not already exist otherwise populate
// the user info with the existing fields.
func (c *Cassandra) CreateUser(i *UserInfo) error {

	err := c.Query(`
		SELECT id, name, picture FROM users WHERE gid = ?`,
		i.GID,
	).Scan(&i.ID, &i.Name, &i.Picture)

	if err == gocql.ErrNotFound {
	} else if err != nil {
		return err
	}

	err = c.Query(`
		INSERT INTO users (gid, id, name, picture) VALUES (?, ?, ?, ?) IF NOT EXISTS`,
		i.GID, i.ID, i.Name, i.Picture,
	).Exec()

	return err
}
