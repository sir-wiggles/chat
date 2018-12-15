package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func RespondWithJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var err error
	var encoder = json.NewEncoder(w)

	switch body.(type) {
	case error:
		log.Printf("server error: %+v", body)
		err = encoder.Encode(response{body.(error).Error()})
	case string:
		err = encoder.Encode(response{body.(string)})
	default:
		err = encoder.Encode(body)
	}

	if err != nil {
		log.Printf("responding with json: %s", err)
	}

}
