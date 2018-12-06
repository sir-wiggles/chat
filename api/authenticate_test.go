package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"github.com/sir-wiggles/chat/api/postgres"
	"golang.org/x/crypto/bcrypt"
)

func bcryptPasswordGen(t *testing.T, password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Error(err)
	}
	return string(hash)
}

var authenticateTT = []struct {
	PlainTextPassword string
	FormUsername      string
	FormPassword      string

	ExpectedStatusCode int
	ExpectedPattern    string
	SQLError           error
}{
	{"abc123", "foo", "abc123",
		http.StatusOK, `{"token":"([a-zA-Z0-9_+-].+){3}}`, nil},

	{"abc123", "foo", "123abc",
		http.StatusUnauthorized, `{"message":"invalid credentials"}`, nil},

	{"abc123", "foo", "abc123",
		http.StatusUnauthorized, `{"message":"invalid credentials"}`, sql.ErrNoRows},

	{"abc123", "foo", "abc123",
		http.StatusInternalServerError, `{"message":"sql:.*"}`, sql.ErrConnDone},
}

func TestAuthenticate(t *testing.T) {
	for _, tt := range authenticateTT {

		pwhash := bcryptPasswordGen(t, tt.PlainTextPassword)
		dbRow := []interface{}{pwhash}

		scanner := &postgres.MockScanner{
			ScanFn: func(dest ...interface{}) error {
				if tt.SQLError != nil {
					return tt.SQLError
				}
				for i, item := range dest {
					*item.(*string) = dbRow[i].(string)
				}
				return nil
			},
		}

		db := &postgres.MockPostgres{
			QueryRowFn: func(query string, args ...interface{}) postgres.Scanner {
				return scanner
			},
		}

		body, _ := json.Marshal(authenticateForm{tt.FormUsername, tt.FormPassword})

		r := httptest.NewRequest("POST", "http://chat.com/authenticate", bytes.NewReader(body))
		w := httptest.NewRecorder()

		c := &Authentication{db: db}
		c = c.SetHandler(c.authenticate)

		c.ServeHTTP(w, r)

		expectToEqual(t, w.Code, tt.ExpectedStatusCode)
		expectToMatch(t, w.Body.String(), tt.ExpectedPattern)
	}

}

var registerTT = []struct {
	FormUsername string
	FormEmail    string
	FormAvatar   string
	FormPassword string

	UserUUID           string
	ExpectedStatusCode int

	ExpectedPattern string
	SQLError        error
}{

	{"foo", "foo@chat.com", "http://icons.ru/foo.png", "abc123",
		"ce83dd2e-6b41-428f-b1fc-5ff3d4f1bacd", http.StatusCreated,
		"", nil},

	{"foo", "foo@chat.com", "http://icons.ru/foo.png", "abc123",
		"ce83dd2e-6b41-428f-b1fc-5ff3d4f1bacd", http.StatusInternalServerError,
		`{"message":"sql:.*"}`, sql.ErrNoRows},
}

func TestRegister(t *testing.T) {

	for _, tt := range registerTT {
		body, _ := json.Marshal(&registerForm{
			tt.FormUsername, tt.FormEmail, tt.FormAvatar, tt.FormPassword,
		})

		dbRow := []interface{}{tt.UserUUID, tt.FormUsername, tt.FormEmail, tt.FormAvatar}

		scanner := &postgres.MockScanner{
			ScanFn: func(dest ...interface{}) error {
				if tt.SQLError != nil {
					return tt.SQLError
				}
				for i, item := range dest {
					*item.(*string) = dbRow[i].(string)
				}
				return nil
			},
		}

		db := &postgres.MockPostgres{
			QueryRowFn: func(query string, args ...interface{}) postgres.Scanner {
				return scanner
			},
		}

		r := httptest.NewRequest("POST", "http://chat.com/register", bytes.NewReader(body))
		w := httptest.NewRecorder()

		c := &Authentication{db: db}
		c = c.SetHandler(c.register)

		c.ServeHTTP(w, r)

		expectToEqual(t, w.Code, tt.ExpectedStatusCode)

		if len(tt.ExpectedPattern) > 0 {
			expectToMatch(t, w.Body.String(), tt.ExpectedPattern)
		} else {
			want := &registerResponse{0, tt.UserUUID, tt.FormUsername, tt.FormEmail, tt.FormAvatar}
			have := &registerResponse{}
			json.Unmarshal(w.Body.Bytes(), have)
			expectToEqual(t, have, want)
		}

	}
}

func expectToEqual(t *testing.T, got, want interface{}) {
	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Expected %v to equal %v", got, want)
	}
}

func expectToMatch(t *testing.T, got, want string) {
	matches, err := regexp.MatchString(want, got)
	if err != nil {
		t.Errorf("regexp match error: %+v", err)
		return
	}

	if matches == false {
		t.Errorf("Expected %s to match %s", got, want)
	}
}
