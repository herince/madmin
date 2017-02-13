// Package madmin/app implements the back-end logic for the vet pharmacy administrating system.
// It implements a simple RESTful API that manages somehow the vet pharmacy warehouse.
package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

// Initialization function of the app. For now it only runs the server on the given port.
func Init(port string, handler http.Handler) {
	madminHandler := NewMAdminHandler()
	http.Handle("/data/", madminHandler)

	http.Handle("/", http.FileServer(http.Dir("static/")))

	log.Println("Listening...")
	http.ListenAndServe(port, nil)
}

type MAdminHandler struct {
	router *mux.Router

	wh *warehouse
}

func (m *MAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func NewMAdminHandler() *MAdminHandler {
	madminHandler := &MAdminHandler{}

	madminHandler.wh = NewWarehouse()

	madminHandler.router = mux.NewRouter()
	madminHandler.router.HandleFunc("/data/stock/{....-..-..-..-......}", madminHandler.stockItemHandler).Methods("GET", "DELETE")
	madminHandler.router.HandleFunc("/data/stock/", madminHandler.stockHandler).Methods("GET", "POST")

	return madminHandler
}

func respondMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Error in request method. Method not allowed - %d", r.Method)
	return
}

func (m *MAdminHandler) stockHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		m.listStockHandler(w, r)
	case "POST":
		m.addStockHandler(w, r)
	default:
		respondMethodNotAllowed(w, r)
	}
}

func (m *MAdminHandler) stockItemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		m.getStockItemHandler(w, r)
	case "DELETE":
		m.removeStockItemHandler(w, r)
	default:
		respondMethodNotAllowed(w, r)
	}
}

// Handler for GET /stock/
// 
// Lists existing stock items.
func (m *MAdminHandler) listStockHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp = &CollectionResponseDTO{"List of existing stock items", make([]string, 0, m.wh.Size())}

		itemUrl string
	)

	for _, item := range m.wh.Stock() {
		itemUrl = fmt.Sprintf("/data/stock/%s", item.Name())
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

// Handler for GET /stock/<id>
// 
// Returns JSON with data for the stock item with the given id (if such item exists in the warehouse)
// or an emptry response with status code 204 (if there is no such item in the warehouse).
func (m *MAdminHandler) getStockItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL
		_, id = path.Split(query.String())

		item, ok = m.wh.Get(id)
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
	resp.Distributor = string(item.Distributor())

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

// Handler for POST /stock/
// 
// Adds an item to the warehouse.
func (m *MAdminHandler) addStockHandler(w http.ResponseWriter, r *http.Request) {
	var (
		newItem = &NewStockDTO{}

		decoder = json.NewDecoder(r.Body)
		err     = decoder.Decode(newItem)
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error in unmarshaling request body: %s", err)
		return
	}
	defer r.Body.Close()

	stockItem, err := NewStock(newItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error in creating stock item: %s", err)
		return
	}
	id := m.wh.Add(stockItem)

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(id)); err != nil {
		log.Printf("Error while writing response: %s", err)
	}
}

// Handler for DELETE /stock/<id>
// 
// Removes the item with <id> from the warehouse.
func (m *MAdminHandler) removeStockItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL
		_, id = path.Split(query.String())
	)
	m.wh.Remove(id)
	w.WriteHeader(http.StatusNoContent)
}
