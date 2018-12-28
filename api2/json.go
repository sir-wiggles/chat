package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// ResponseWriter wraps the incoming http.ResponseWriter so we can handle returning JSON structs
type ResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (w *ResponseWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

// JSON takes a struct and encodes it to JSON.  If the value is an error, then the error will
// be returned as an object with error as the key and value as the .Error() of the error.
//
// Status is the http status code to be returned with the response.
// If no status is given, then it will default to StatusOK for non error value types; otherwise,
// it will be StatusInternalServerError
//
// If more than one status code is given, only the first status will be used.
//
// An error will be returned if an error occurred; however, this is just for information as the
// header will be set to StatusInternalServerError upon an error
func (w *ResponseWriter) JSON(v interface{}, status ...int) error {

	var code int

	switch v.(type) {
	case validator.ValidationErrors:
		code = http.StatusBadRequest

		errs := make(map[string]string, 2)
		for _, err := range v.(validator.ValidationErrors) {
			errs[err.Field()] = err.Tag()
		}
		v = struct {
			Error map[string]string `json:"error"`
		}{errs}
	case error:

		code = http.StatusInternalServerError
		v = struct {
			Error string `json:"error"`
		}{v.(error).Error()}

	default:
		code = http.StatusOK
	}

	if len(status) > 1 {
		code = status[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w.buf).Encode(v); err != nil {
		log.Printf("ResponseWriter.JSON error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

// JSONMiddleWare wraps response writer allowing the ability to respond with a struct as JSON.
// Example:
//    err := w.(*ResponseWriter).JSON(struct)
func JSONMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		rw := &ResponseWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
		}

		h.ServeHTTP(rw, r)

		if _, err := io.Copy(w, rw.buf); err != nil {
			log.Printf("Failed to send out response: %v", err)
		}
	})
}
