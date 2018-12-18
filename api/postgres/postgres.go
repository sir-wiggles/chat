package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Controller exposes the methods needed to mock
type Controller interface {
	QueryRow(query string, args ...interface{}) Scanner
	GetOrCreateUser(user *UserModel) error
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

// New creates a new database connection that implements the Postgreser interface
func New(url string) (*Postgres, error) {
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

type UserModel struct {
	GID        string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	WasCreated bool   `json:"was_created"`
}

// GetOrCreateUser is a query that will create a user if one does not exist or retrieve an
// existing one. There are four args in this query: name, email picture and gid respectively.
func (db *Postgres) GetOrCreateUser(user *UserModel) error {

	var (
		name    = user.Name
		email   = user.Email
		picture = user.Picture
		gid     = user.GID
	)

	const query = `
		WITH row AS (
		INSERT INTO
			users (name, email, picture, gid)
		SELECT
			$1, $2, $3, $4
		WHERE
			NOT EXISTS (
				SELECT
					id
				FROM
					users
				WHERE
					gid = $4
			) RETURNING *
		)
		SELECT
			FALSE as existing
		FROM
			row
		UNION
		SELECT
			TRUE as existing
		FROM
			users
		WHERE
			gid = $4;`

	return db.QueryRow(query, name, email, picture, gid).
		Scan(&user.WasCreated)

}
