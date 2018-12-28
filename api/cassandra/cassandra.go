package cassandra

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/sir-wiggles/chat/api/structs"
)

type Controller interface {
	LogMessage(string, string, string) error
	GetUser(string, string, string) (*structs.User, error)
}

type Cassandra struct {
	*gocql.Session
}

func New(urls []string) (*Cassandra, error) {

	cluster := gocql.NewCluster(urls...)
	cluster.Keyspace = "chatter"

	session, err := cluster.CreateSession()
	return &Cassandra{session}, err
}

func (c *Cassandra) LogMessage(cid, oid, body string) error {
	query := `INSERT INTO
		messages (channel, id, owner, body)
	VALUES (?, ?, ?, ?)`
	return c.Query(query, cid, gocql.TimeUUID(), oid, body).Exec()
}

func (c *Cassandra) GetMessages(cid string, limit int) ([]*structs.Message, error) {
	var (
		query     = `SELECT totimestamp(id), owner, body FROM messages WHERE channel = ? LIMIT ?`
		messages  = []*structs.Message{}
		timestamp time.Time
		owner     string
		body      string
	)
	iter := c.Query(query, cid, limit).Iter()

	for iter.Scan(&timestamp, &owner, &body) {
		messages = append(messages, structs.NewMessage(cid, owner, body, timestamp))
	}
	return messages, iter.Close()
}

func (c *Cassandra) IsUserPermitted(cid, mid string) (bool, error) {
	var (
		query  = `SELECT COUNT(1) FROM channels WHERE id = ? AND members CONTAINS ?;`
		result int
		valid  bool
	)
	iter := c.Query(query, cid, mid).Iter()
	iter.Scan(&result)
	err := iter.Close()
	if result == 1 {
		valid = true
	}
	return valid, err
}

func (c *Cassandra) RemoveUserFromChannel(cid, oid, mid string) (bool, error) {
	var (
		query = `UPDATE channels SET members = members - {?} WHERE id = ? AND owner = ? IF EXISTS;`
		valid bool
	)

	iter := c.Query(query, mid, cid, oid).Iter()
	iter.Scan(&valid)
	err := iter.Close()
	return valid, err

}

func (c *Cassandra) AddUserFromChannel(cid, oid, mid string) (bool, error) {
	var (
		query = `UPDATE channels SET members = members + {?} WHERE id = ? AND owner = ? IF EXISTS;`
		valid bool
	)
	iter := c.Query(query, mid, cid, oid).Iter()
	iter.Scan(&valid)
	err := iter.Close()
	return valid, err

}

func (c *Cassandra) DeleteChannel(cid, oid string) (bool, error) {
	var (
		query = `DELETE FROM channels WHERE id = ? AND owner = ? IF EXISTS;`
		valid bool
	)
	iter := c.Query(query, cid, oid).Iter()
	iter.Scan(&valid)
	err := iter.Close()
	return valid, err
}

func (c *Cassandra) GetUsersInChannel(cid string) ([]*structs.User, error) {
	var (
		memberQuery = `SELECT members FROM channels WHERE id = ?`
		nameQuery   = `SELECT id, name, picture FROM users WHERE id IN ?`
		uuids       = make([]string, 0, 2)
		uuid        string
		users       = make([]*structs.User, 0, 2)
		name        string
		picture     string
	)

	iter := c.Query(memberQuery, cid).Iter()
	for iter.Scan(&uuid) {
		uuids = append(uuids, uuid)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	iter = c.Query(nameQuery, uuids).Iter()
	for iter.Scan(&uuid, &name, &picture) {
		users = append(users, structs.NewUser(uuid, name, picture))
	}
	return users, iter.Close()
}

func (c *Cassandra) GetUser(gid, name, picture string) (*structs.User, error) {

	var (
		id   string
		user *structs.User
	)

	iter := c.Query(`SELECT id, name, picture FROM users where gid = ?`, gid).Iter()
	defer iter.Close()

	if iter.NumRows() == 1 {
		iter.Scan(&id, &name, &picture)
		user = structs.NewUser(id, name, picture)
		return user, iter.Close()
	}

	err := c.Query(
		`INSERT INTO users (gid, id, name, picture) VALUES (?, uuid(), ?, ?) IF NOT EXISTS`,
		gid,
		name,
		picture,
	).Exec()

	return structs.NewUser(id, name, picture), err

}
