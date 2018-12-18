package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sir-wiggles/chat/api/cassandra"
	"github.com/sir-wiggles/chat/api/postgres"
)

func init() {

	data, err := ioutil.ReadFile("./oauth.json")
	if err != nil {
		log.Fatal(err)
	}

	dataReader := bytes.NewReader(data)

	clientInfo := struct {
		ID     string `json:"client_id"`
		Secret string `json:"client_secret"`
	}{}

	err = json.NewDecoder(dataReader).Decode(&clientInfo)
	if err != nil {
		log.Fatal(err)
	}

	googleClientID = clientInfo.ID
	googleClientSecret = clientInfo.Secret

	log.Println(clientInfo)

}

var (
	port               = os.Getenv("PORT")
	host               = os.Getenv("HOST")
	postgresURL        = os.Getenv("POSTGRES_URL")
	cassandraURL       = os.Getenv("CASSANDRA_URL")
	cassandraURLs      []string
	corsAllowedHeaders = os.Getenv("CORS_ALLOWED_HEADERS")
	corsAllowedMethods = os.Getenv("CORS_ALLOWED_METHODS")
	corsAllowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")

	secretSigningKey = os.Getenv("JWT_SECRET_KEY")
	jwtIssuer        = os.Getenv("JWT_ISSUER")
	jwtExpiresAt     int64
	_jwtExpiresAt    = os.Getenv("JWT_EXPIRES_IN_MINUTES")

	broadcastChannelBufferSize  = 8
	registerChannelBufferSize   = 8
	unregisterChannelBufferSize = 8

	googleClientID     string
	googleClientSecret string

	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

func main() {

	flags()

	db, err := postgres.New(postgresURL)
	if err != nil {
		log.Fatalf("Postgres Connection Error: %s", err)
	}

	cass, err := cassandra.New(cassandraURLs)
	if err != nil {
		log.Fatal("Cassandra Connection Error: %s", err)
	}

	var (
		auth    = NewAuthenticationController(db)
		router  = mux.NewRouter()
		chat    = NewClientManager(cass)
		address = fmt.Sprintf("%s:%s", host, port)
	)
	router.NotFoundHandler = &NotFoundHandler{}

	authR := router.NewRoute().PathPrefix("/auth").Methods("POST").Subrouter()
	authR.Handle("/google", auth.SetHandler(auth.Google))

	apiR := router.NewRoute().PathPrefix("/api").Subrouter()
	apiR.Use(auth.Middleware)
	apiR.Handle("/ws", chat).Methods("GET").Queries("token", "{token}")
	apiR.HandleFunc("/health", health).Methods("GET")

	router.HandleFunc("/chat", index)

	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	var headersOk = handlers.AllowedHeaders(strings.Split(corsAllowedHeaders, ","))
	var methodsOk = handlers.AllowedMethods(strings.Split(corsAllowedMethods, ","))
	var originsOk = handlers.AllowedOrigins(strings.Split(corsAllowedOrigins, ","))

	var handler = handlers.LoggingHandler(os.Stdout, router)
	handler = handlers.CORS(headersOk, originsOk, methodsOk)(handler)

	var server = http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Printf("server listening on %s", address)
	log.Fatal(server.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	file, err := ioutil.ReadFile("./static/index.html")

	log.Println("reading flie error: ", err)
	fmt.Fprint(w, string(file))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type NotFoundHandler struct{}

func (h NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func flags() {
	var err error

	flag.StringVar(&port, "port", port, "port the server will listen on")
	flag.StringVar(&host, "host", host, "host to serve on")
	flag.StringVar(&postgresURL, "postgres", postgresURL, "postgres url")
	flag.StringVar(&cassandraURL, "cassandra", cassandraURL, "cassandra url")
	flag.StringVar(&corsAllowedHeaders, "corsAllowedHeaders", corsAllowedHeaders, "headers allowed for cors")
	flag.StringVar(&corsAllowedMethods, "corsAllowedMethods", corsAllowedMethods, "methods allowed for cors")
	flag.StringVar(&corsAllowedOrigins, "corsAllowedOrigins", corsAllowedOrigins, "origins allowed for cors")
	flag.StringVar(&secretSigningKey, "jwtSecretKey", secretSigningKey, "secret for jwt token signing")
	flag.StringVar(&jwtIssuer, "jwtIssuer", jwtIssuer, "issuer of the jwt token")

	jwtExpiresAt, err = strconv.ParseInt(_jwtExpiresAt, 10, 64)
	if err != nil {
		log.Fatalf("Invalid value for jwtExpiresAt: %s should be a number", _jwtExpiresAt)
	}
	flag.Int64Var(&jwtExpiresAt, "jwtExpiresAt", jwtExpiresAt, "expiration of the jwt token")

	cassandraURLs = strings.Split(cassandraURL, ",")

	flag.Parse()
}
