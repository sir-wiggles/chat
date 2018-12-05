package main

import (
	"encoding/json"
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
	switch body.(type) {
	case string:
		body = HTTPResponseBody{body}
	}
	return &HTTPResponse{
		code: code,
		body: body,
	}
}

type HTTPResponseBody struct {
	Message interface{} `json:"message"`
}

// AsJSON wraps a ModifiedHTTPHandler formatting the response as a JSON object
func AsJSON(w http.ResponseWriter, r HTTPResponder) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code())
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "")
	encoder.Encode(r.Body())
}
