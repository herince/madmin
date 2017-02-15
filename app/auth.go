package app

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authString := r.Header.Get("Authorization")
		if len(authString) == 0 {
			writeStatusUnauthorized(w, r)
		} else {
			name, password, err := decodeAuthHeader(authString)
			if err != nil {
				writeStatusUnauthorized(w, r)
			}

			isValidUser := validateUser(name, password)

			if isValidUser {
				handler.ServeHTTP(w, r)
			} else {
				writeStatusUnauthorized(w, r)
			}
		}
	})
}

func decodeAuthHeader(authString string) (name string, password string, err error) {
	if !strings.HasPrefix(authString, "Basic ") {
		err = errors.New("Invalid authorization header: unknown authorization scheme.")
		return
	}

	encodedPart := authString[6:]

	credentials, err := base64.StdEncoding.DecodeString(encodedPart)
	if err != nil {
		return
	}

	userData := strings.SplitN(string(credentials), ":", 2)
	if len(userData) < 2 {
		err = errors.New("Invalid authorization header: bad user password format.")
		return
	}

	name = userData[0]
	password = userData[1]

	return
}

func validateUser(name, password string) bool {
	return true
}

func writeStatusUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("WWW-Authenticate", "Basic realm=\"mira administrator\"")
	w.WriteHeader(http.StatusUnauthorized)
}
