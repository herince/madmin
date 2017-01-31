package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func Init(port string) {
	madminHandler := NewMAdminHandler()
	http.Handle("/", madminHandler)

	log.Println("Listening...")
	http.ListenAndServe(port, nil)
}

type MAdminHandler struct {
	router *mux.Router

	wh *Warehouse
}

func (m *MAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func NewMAdminHandler() *MAdminHandler {
	madminHandler := &MAdminHandler{}

	madminHandler.wh = &Warehouse{make(map[string]Stock)}

	madminHandler.router = mux.NewRouter()
	madminHandler.router.HandleFunc("/stock/{....-..-..-..-......}", madminHandler.stockItemHandler).Methods("GET")
	madminHandler.router.HandleFunc("/stock/", madminHandler.stockRequestsHandler).Methods("GET", "POST")

	return madminHandler
}

/*
 * JSON format for the GET requests for "/<items>/"
 */
type CollectionResponse struct {
	Info string   `json:"info"`
	Urls []string `json:"urls"`
}

func (m *MAdminHandler) stockRequestsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			m.stockListHandler(w, r)
		case "POST":
			m.addStockHandler(w, r)
	}
}

/*
 * Handler for GET /stock/
 * 
 * Lists existing stock items.
 */
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
type StockDTO struct {
	Id string `json:"id"`

	Name string    `json:"name"`
	Type StockType `json:"type"`

	ExpirationDate string `json:"expirationDate"`
	MinQuantity    string `json:"minQuantity"`

	Distributor string `json:"distributor"`
}

/*
 * Handler for GET /stock/<id>
 * 
 * Returns JSON with data for the stock item with the given id (if such item exists in the warehouse)
 * or an emptry response with status code 204 (if there is no such item in the warehouse).
 */
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

	resp := &StockDTO{}

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

/*
 * Handler for POST /stock/
 * 
 * Adds an item to the warehouse.
 */
func (m *MAdminHandler) addStockHandler(w http.ResponseWriter, r *http.Request) {
	var (
		newItem = &StockDTO{}

		decoder = json.NewDecoder(r.Body)
		err = decoder.Decode(newItem)
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error in unmarchaling request body: %s", err)
		return
	}
	defer r.Body.Close()

	log.Println(newItem)

	stockItem, err := NewStock(newItem.Type, newItem.Name)
	m.wh.Add(stockItem)
	
	log.Println(m.wh)
	// todo - test it properly.
}
