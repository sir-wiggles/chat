package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sir-wiggles/chat/api/cassandra"
	"github.com/sir-wiggles/chat/api/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Authentication struct {
	db      cassandra.Controller
	handler http.HandlerFunc
	google  *oauth2.Config
}

func NewAuthenticationController(db cassandra.Controller) *Authentication {
	googleConfig := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  "postmessage",
		Scopes: []string{
			"profile",
			"email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Authentication{
		db:     db,
		google: googleConfig,
	}
}

func (c Authentication) SetHandler(handler http.HandlerFunc) *Authentication {
	return &Authentication{
		db:      c.db,
		google:  c.google,
		handler: handler,
	}
}

func (c Authentication) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler(w, r)
}

type customJWTClaims struct {
	jwt.StandardClaims
	UserModel postgres.UserModel `json:"user"`
}

type GoogleAuthRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirectURI"`
}

type AuthResponse struct {
	Token     string              `json:"token"`
	UserModel *postgres.UserModel `json:"user"`
}

func (c Authentication) Google(w http.ResponseWriter, r *http.Request) {

	// Exchange code for token
	var gar = GoogleAuthRequest{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&gar); err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	token, err := c.google.Exchange(oauth2.NoContext, gar.Code)
	if err != nil {
		RespondWithJSON(w, http.StatusUnauthorized, err)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", googleUserInfoURL, token.AccessToken))
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err)
	}

	// Register the user in the system
	var gui = postgres.UserModel{}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&gui); err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	user, err := c.db.GetUser(gui.GID, gui.Name, gui.Picture)
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	// Create a token for auth
	claims := customJWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "chatter",
		},
		gui,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtTokenString, err := jwtToken.SignedString([]byte(secretSigningKey))
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err)
	}

	authResponse := AuthResponse{Token: jwtTokenString}
	RespondWithJSON(w, http.StatusOK, authResponse)

}

type ContextKey string

const (
	ContextName    ContextKey = "name"
	ContextPicture ContextKey = "picture"
	ContextGID     ContextKey = "gid"
)

func (c *Authentication) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.URL.Query().Get("token")
		if len(tokenString) == 0 {
			tokenString = r.Header.Get("Authorization")
			if len(tokenString) <= 7 {
				RespondWithJSON(w, http.StatusBadRequest, "invalid token")
				return
			}
			tokenString = tokenString[7:]
		}

		token, err := jwt.ParseWithClaims(tokenString, &customJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretSigningKey), nil
		})
		if err != nil {
			RespondWithJSON(w, http.StatusForbidden, err)
			return
		}

		if _, ok := token.Claims.(*customJWTClaims); !ok && !token.Valid {
			RespondWithJSON(w, http.StatusForbidden, "invalid claims")
			return
		}

		claims := token.Claims.(*customJWTClaims)
		fmt.Println("Claims", claims.UserModel)
		r = r.WithContext(context.WithValue(r.Context(), ContextName, claims.UserModel.Name))
		r = r.WithContext(context.WithValue(r.Context(), ContextPicture, claims.UserModel.Picture))
		r = r.WithContext(context.WithValue(r.Context(), ContextGID, claims.UserModel.GID))

		handler.ServeHTTP(w, r)
	})
}
