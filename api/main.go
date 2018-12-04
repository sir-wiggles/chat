package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	port               = os.Getenv("PORT")
	host               = os.Getenv("HOST")
	postgresURL        = os.Getenv("POSTGRES_URL")
	corsAllowedHeaders = os.Getenv("CORS_ALLOWED_HEADERS")
	corsAllowedMethods = os.Getenv("CORS_ALLOWED_METHODS")
	corsAllowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")
	endpoint           string

	broadcastChannelBufferSize  = 8
	registerChannelBufferSize   = 8
	unregisterChannelBufferSize = 8
)

func main() {

	flags()

	db, err := postgresConnect(postgresURL)
	if err != nil {
		log.Fatalf("Postgres Connection Error: %s", err)
	}

	var (
		reg     = NewRegistrationController(db)
		mux     = mux.NewRouter()
		chat    = NewClientManager()
		address = fmt.Sprintf("%s:%s", host, port)
	)

	mux.Handle(endpoint, chat).Methods("GET")

	mux.Handle("/register", reg.SetHandler(reg.register)).Methods("POST")

	mux.PathPrefix("/images/").
		Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	var headersOk = handlers.AllowedHeaders(strings.Split(corsAllowedHeaders, ","))
	var methodsOk = handlers.AllowedMethods(strings.Split(corsAllowedMethods, ","))
	var originsOk = handlers.AllowedOrigins(strings.Split(corsAllowedOrigins, ","))

	var handler = handlers.LoggingHandler(os.Stdout, mux)
	handler = handlers.CORS(headersOk, originsOk, methodsOk)(handler)

	var server = http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Printf("server listening on %s\n", address)
	log.Fatal(server.ListenAndServe())
}

func postgresConnect(url string) (*sql.DB, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Postgres ping error: %s", err)
	}
	return db, err

}

func flags() {
	flag.StringVar(&port, "port", port, "port the server will listen on")
	flag.StringVar(&host, "host", host, "host to serve on")
	flag.StringVar(&endpoint, "endpoint", "/ws", "the path of the web socket")
	flag.StringVar(&postgresURL, "postgres", postgresURL, "postgres url")
	flag.StringVar(&corsAllowedHeaders, "corsAllowedHeaders", corsAllowedHeaders, "headers allowed for cors")
	flag.StringVar(&corsAllowedMethods, "corsAllowedMethods", corsAllowedMethods, "methods allowed for cors")
	flag.StringVar(&corsAllowedOrigins, "corsAllowedOrigins", corsAllowedOrigins, "origins allowed for cors")
	flag.Parse()
}
