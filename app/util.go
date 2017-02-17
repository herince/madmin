package app

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
)

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func respondMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Error in request method. Method not allowed - %s", r.Method)
	return
}

func respondStatusUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("WWW-Authenticate", "Basic realm=\"mira administrator\"")
	w.WriteHeader(http.StatusUnauthorized)
	return
}
