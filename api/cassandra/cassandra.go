package cassandra

import (
	"log"

	"github.com/gocql/gocql"
)

type Controller interface {
	Log(string, string)
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

func (c *Cassandra) Log(from, body string) {
	query := `
	INSERT INTO
		chat_messages (from_user, to_user, time, body)
	VALUES (?, 'channel', ?, ?)`
	err := c.Query(query, from, gocql.TimeUUID(), body).Exec()
	if err != nil {
		log.Println(err)
	}
}
