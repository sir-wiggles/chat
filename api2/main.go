package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	host     string
	port     string
	keyspace string
)

// UUIDPattern used to match UUID patterns in urls
const UUIDPattern = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"

func init() {
	getHost()
	getPort()
	getKeyspace()
}

func main() {

	var (
		address = fmt.Sprintf("%s:%s", host, port)
		router  = mux.NewRouter().StrictSlash(true)
		api     = router.NewRoute().PathPrefix("/api").Subrouter()
		handler http.Handler

		chatter = &Chatter{}
		channel = &Channel{}
		user    = &User{}
	)

	chatter.Register(api)
	channel.Register(api)
	user.Register(api)

	api.Use(JSONMiddleWare)
	handler = handlers.LoggingHandler(os.Stdout, router)

	srv := http.Server{
		Handler:      handler,
		Addr:         address,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	log.Printf("server address %s", address)
	srv.ListenAndServe()

}

func getHost() string {
	host = os.Getenv("HOST")
	host = strings.Trim(host, " ")
	if len(host) == 0 {
		host = "localhost"
	}
	return host
}

func getPort() string {
	port = os.Getenv("PORT")
	port = strings.Trim(port, " ")
	if len(port) == 0 {
		port = "5050"
	}
	return port
}

func getKeyspace() string {
	keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	keyspace = strings.Trim(port, " ")
	if len(keyspace) == 0 {
		keyspace = "chatter"
	}
	return keyspace
}
