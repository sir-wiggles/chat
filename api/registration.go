package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Registration handles all forms of user registration.
type Registration struct {
	db      *sql.DB
	handler ModifiedHTTPHandler
}

type registrationForm struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

// NewRegistrationController returns a registration controller with reference to the supplied
// database connection. The user will need to call SetHandler to attach a handler
func NewRegistrationController(db *sql.DB) *Registration {
	return &Registration{
		db: db,
	}
}

// SetHandler will return a new Registration struct with the handler set to the ModifiedHTTPHandler
// with reference to the parent database connection
func (c Registration) SetHandler(handler ModifiedHTTPHandler) *Registration {
	return &Registration{
		db:      c.db,
		handler: handler,
	}
}

func (c Registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AsJSON(w, c.handler(w, r))
}

func (c Registration) register(w http.ResponseWriter, r *http.Request) HTTPResponder {
	var (
		form                      = registrationForm{}
		id                        int
		uuid, name, email, avatar string
		existing                  bool
		err                       error
	)

	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&form); err != nil {
		return NewMalformedBodyError(err)
	}

	args := []interface{}{form.Name, form.Email, form.Avatar}
	dest := []interface{}{&id, &uuid, &name, &email, &avatar, &existing}

	if err = c.db.QueryRow(GetOrCreateUser, args...).Scan(dest...); err != nil {
		return NewSQLQueryRowError(err)
	}

	code := http.StatusCreated
	if existing {
		code = http.StatusOK
	}

	return NewHTTPResponse(code, &UserModel{
		ID: id, UUID: uuid, Name: name, Email: email, Avatar: avatar,
	})
}

// UserModel represents a user row in the database
type UserModel struct {
	ID     int    `json:"-"`
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

// GetOrCreateUser is a query that will create a user if one does not exist or retrieve an
// existing one. There are three args in this query: name, email and avatar respectively.
const GetOrCreateUser = `
WITH row AS (
INSERT INTO
	users (name, email, avatar)
SELECT
	$1, $2, $3
WHERE
	NOT EXISTS (
		SELECT
			*
		FROM
			users
		WHERE
			name = $1
	) RETURNING *
)
SELECT
	id, uuid, name, email, avatar, FALSE as existing
FROM
	row
UNION
SELECT
	id, uuid, name, email, avatar, TRUE as existing
FROM
	users
WHERE
	name = $1;`
