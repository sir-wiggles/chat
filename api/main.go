package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port     int
	host     string
	endpoint string

	broadcastChannelBufferSize  = 8
	registerChannelBufferSize   = 8
	unregisterChannelBufferSize = 8
)

func main() {

	flag.IntVar(&port, "port", 5050, "port the server will listen on")
	flag.StringVar(&host, "host", "localhost", "host to serve on")
	flag.StringVar(&endpoint, "endpoint", "/ws", "the path of the web socket")
	flag.Parse()

	var (
		mux     = mux.NewRouter()
		chat    = NewClientManager()
		address = fmt.Sprintf("%s:%d", host, port)
	)

	mux.Handle(endpoint, chat).Methods("GET")
	mux.PathPrefix("/images/").
		Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	log.Printf("server listening on %s\n", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
