package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MAdminHandler struct {
	mux *http.ServeMux

	wh Warehouse
}

func (m *MAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func NewMAdminHandler() *MAdminHandler {
	mah := &MAdminHandler{}

	mah.mux = http.NewServeMux()
	mah.mux.HandleFunc("/stock/", RestrictedMethodHandler("GET", mah.stockListHandler))

	return mah
}

/*
 * JSON format for the GET requests for "/<items>/"
 */
type CollectionResponse struct {
	Info string   `json:"info"`
	Urls []string `json:"urls"`
}

// lists existing stock items
// handler for GET /stock/
func (m *MAdminHandler) stockListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp = &CollectionResponse{"List of existing stock items", make([]string, 0, len(m.wh.Stock))}

		itemUrl string
	)

	for _, item := range m.wh.Stock {
		itemUrl = fmt.Sprintf("/stock/%s", item.Name())
		resp.Urls = append(resp.Urls, itemUrl)
	}

	// todo - use the CollectionResponse structure
	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in marshalling results: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBytes); err != nil {
		log.Printf("Error while writing response: %s", err)
	}
}

func RestrictedMethodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}
