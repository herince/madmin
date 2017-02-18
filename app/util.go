package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"time"
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

func validDateFromString(dateString string) (date time.Time, err error) {
	if dateString == "" {
		return time.Unix(0, 0), errors.New("expected expiration dateString but not set")
	}
	layout := "2006-01-02T15:04:05.000Z"
	date, err = time.Parse(layout, dateString)
	if err != nil {
		return
	}
	return
}

func validQuantityFromString(quantityString string) (quantity decimal.Decimal, err error) {
	if quantityString == "" {
		return decimal.Zero, errors.New("no quantity set for stock")
	}
	quantity, err = decimal.NewFromString(quantityString)
	if err != nil {
		return
	}
	return
}
