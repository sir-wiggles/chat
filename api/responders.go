package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ModifiedHTTPHandler allows for a handler to return a HTTPResponder to be handled in the ServeHTTP
// method. This allows us to intercept and format the response to a JSON response without having to
// litter all handlers with JSON encoding.
type ModifiedHTTPHandler func(w http.ResponseWriter, r *http.Request) HTTPResponder

// HTTPResponder is the interface that is returned from user created handlers
type HTTPResponder interface {
	Code() int
	Body() interface{}
}

// HTTPResponse is the generic response with the http status code to return and body to for into
// a JSON response.
type HTTPResponse struct {
	code int
	body interface{}
}

func (r HTTPResponse) Error() string {
	return ""
}

// Code returns the http status code from the HTTPResponder
func (r HTTPResponse) Code() int {
	return r.code
}

// Body returns the desired body to send to the client
func (r HTTPResponse) Body() interface{} {
	return r.body
}

// NewHTTPResponse returns the generic HTTPResponse with the supplied http status code and body
func NewHTTPResponse(code int, body interface{}) HTTPResponder {
	return &HTTPResponse{
		code: code,
		body: body,
	}
}

// MalformedBodyError should be returned when a body is supplied that is malformed
type MalformedBodyError struct {
	HTTPResponse
	Message string `json:"message"`
}

// NewMalformedBodyError wil return a HTTPResponder with the status code of BadRequest and body
// specifying the error that happend
func NewMalformedBodyError(err error) HTTPResponder {
	return &HTTPResponse{
		code: http.StatusBadRequest,
		body: &MalformedBodyError{
			Message: fmt.Sprintf("Malformed request body: %s", err),
		},
	}
}

// SQLQueryRowError should be returned when QueryRow returns an error
type SQLQueryRowError struct {
	HTTPResponse
	Message string `json:"message"`
}

// NewSQLQueryRowError will return a HTTPResponder with the status code of InternalServerError and
// body specifying the error that happened
func NewSQLQueryRowError(err error) HTTPResponder {
	return &HTTPResponse{
		code: http.StatusInternalServerError,
		body: &SQLQueryRowError{
			Message: fmt.Sprintf("SQLQueryRowError: %s", err),
		},
	}
}

// AsJSON wraps a ModifiedHTTPHandler formatting the response as a JSON object
func AsJSON(w http.ResponseWriter, r HTTPResponder) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code())
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "")
	encoder.Encode(r.Body())
}
