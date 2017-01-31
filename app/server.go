package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type MAdminHandler struct {
	router *mux.Router

	wh Warehouse
}

func (m *MAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func NewMAdminHandler() *MAdminHandler {
	madminHandler := &MAdminHandler{}

	madminHandler.router = mux.NewRouter()
	madminHandler.router.HandleFunc("/stock/{....-..-..-..-......}", madminHandler.stockItemHandler).Methods("GET")
	madminHandler.router.HandleFunc("/stock/", madminHandler.stockListHandler).Methods("GET")

	return madminHandler
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

/*
 * JSON format for the GET requests for "/<items>/<id>"
 */
type StockItemResponse struct {
	Id string `json:"id"`

	Name string    `json:"name"`
	Type StockType `json:"type"`

	ExpirationDate string `json:"expirationDate"`
	MinQuantity    string `json:"minQuantity"`

	Distributor string `json:"distributor"`
}

func (m *MAdminHandler) stockItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL
		_, id = path.Split(query.String())

		item, ok = m.wh.Stock[id]
	)
	if !ok {
		w.Header().Add("ContentLength", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := &StockItemResponse{}

	resp.Id = item.Id()
	resp.Name = item.Name()
	resp.Type = item.Type()
	if item.IsExpirable() {
		resp.ExpirationDate = item.ExpirationDate().String()
	}
	resp.MinQuantity = item.MinQuantity().String()
	resp.Distributor = item.Distributor()

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
