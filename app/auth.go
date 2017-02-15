package app

import (
	"time"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
)

// global string for the secret
var signingKey = []byte("secret")

func newTokenHandler() *tokenHandler {
	th := &tokenHandler{}

	th.router = mux.NewRouter()

	th.router.HandleFunc("/auth/get-token", th.getTokenHandler).Methods("GET")

	return th
}

type tokenHandler struct {
	router *mux.Router
}

func (t *tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.ServeHTTP(w, r)
}

func (th *tokenHandler) getTokenHandler(w http.ResponseWriter, r *http.Request) {
	// 	new token
    token := jwt.New(jwt.SigningMethodHS256)

	// a map to store the claims
    claims := token.Claims.(jwt.MapClaims)

	// set token claims 
    claims["admin"] = true
    claims["name"] = "asdf"
    claims["exp"] = time.Now().Add(time.Hour).Unix()

	// sign the token with the secret 
    tokenString, _ := token.SignedString(signingKey)

    w.Write([]byte(tokenString))
}
