package madmin

import (
	"fmt"
	"net/http"
)

type DalekHandler struct {
	w Warehouse
}

func (d DalekHandler) ServeHTTP (rw ResponseWriter, req *Request) {
	mux := http.NewServeMux()

	// todo - write the RESTful stuff
	mux.HandleFunc("/", func (rw ResponseWriter, req *Request) {})
}