package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sir-wiggles/chat/api/postgres"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	db      postgres.Controller
	handler ModifiedHTTPHandler
}

func NewAuthenticationController(db postgres.Controller) *Authentication {
	return &Authentication{
		db: db,
	}
}

func (c Authentication) SetHandler(handler ModifiedHTTPHandler) *Authentication {
	return &Authentication{
		db:      c.db,
		handler: handler,
	}
}

func (c Authentication) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AsJSON(w, c.handler(w, r))
}

type authenticateForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authenticateResponse struct {
	Token string `json:"token"`
}

type customJWTClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

func (c Authentication) authenticate(w http.ResponseWriter, r *http.Request) HTTPResponder {
	var err error

	var form = authenticateForm{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return NewHTTPResponse(http.StatusBadRequest, err)
	}

	var password, uuid string
	err = c.db.QueryRow(queryGetUserPassword, form.Username).Scan(&password, &uuid)
	if err == sql.ErrNoRows {
		return NewHTTPResponse(http.StatusUnauthorized, "invalid credentials")
	} else if err != nil {
		return NewHTTPResponse(http.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(form.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return NewHTTPResponse(http.StatusUnauthorized, "invalid credentials")
	} else if err != nil {
		return NewHTTPResponse(http.StatusInternalServerError, err.Error())
	}

	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, customJWTClaims{
		uuid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(jwtExpiresAt)).Unix(),
			Issuer:    jwtIssuer,
		},
	}).SignedString([]byte(secretSigningKey))
	if err != nil {
		return NewHTTPResponse(http.StatusInternalServerError, err.Error())
	}

	return NewHTTPResponse(http.StatusOK, &authenticateResponse{ss})
}

type registerForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID       int    `json:"-"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func (c Authentication) register(w http.ResponseWriter, r *http.Request) HTTPResponder {
	var (
		form                          = registerForm{}
		uuid, username, email, avatar string
		err                           error
	)

	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&form); err != nil {
		return NewHTTPResponse(http.StatusBadRequest, err.Error())
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return NewHTTPResponse(http.StatusInternalServerError, err.Error())
	}

	args := []interface{}{form.Username, form.Email, form.Avatar, passwordHash}
	dest := []interface{}{&uuid, &username, &email, &avatar}

	if err = c.db.QueryRow(queryCreateUser, args...).Scan(dest...); err != nil {
		return NewHTTPResponse(http.StatusInternalServerError, err.Error())
	}

	return NewHTTPResponse(http.StatusCreated, &registerResponse{
		UUID:     uuid,
		Username: username,
		Email:    email,
		Avatar:   avatar,
	})
}

func AuthenticateRequest(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		token, err := jwt.ParseWithClaims(tokenString, &customJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretSigningKey), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if claims, ok := token.Claims.(*customJWTClaims); ok && token.Valid {
			fmt.Printf("%v %v\n", claims.UUID, claims.StandardClaims.ExpiresAt)
		} else {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r)
	})
}

const queryGetUserPassword = `
SELECT
  password, uuid
FROM
  users
WHERE
  username = $1
`

const queryCreateUser = `
INSERT INTO
	users (username, email, avatar, password)
VALUES
	($1, $2, $3, $4)
RETURNING
	uuid, username, email, avatar;
`

// GetOrCreateUser is a query that will create a user if one does not exist or retrieve an
// existing one. There are three args in this query: name, email and avatar respectively.
const queryGetOrCreateUser = `
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
