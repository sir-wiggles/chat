package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Controller exposes the methods needed to mock
type Controller interface {
	QueryRow(query string, args ...interface{}) Scanner
}

// Scanner provides an interface around the Row scan function for easy testing
type Scanner interface {
	Scan(dest ...interface{}) error
}

// Postgres is the generic database connection
type Postgres struct {
	*sql.DB
	conn *sql.DB
}

// NewPostgres creates a new database connection that implements the Postgreser interface
func NewPostgres(url string) (*Postgres, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{conn: conn}, nil
}

// QueryRow wraps the default sql.QueryRow but with a custom Scanner interface for testing
func (db Postgres) QueryRow(query string, args ...interface{}) Scanner {
	return db.conn.QueryRow(query, args...)
}
