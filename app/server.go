package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gorilla/mux"
)

// Init function of the app. For now it only runs the server on the given port.
func Init(port string) {
	var (
		dbPath   = "./database/database.sqlite"
		database = newDB(dbPath)
	)

	// urls for warehouse management API
	maHandler := newMAdminHandler(database)
	http.Handle("/data/", authMiddleware(maHandler))

	http.Handle("/", authMiddleware(http.FileServer(http.Dir("static/"))))

	registerCleanUp(database)

	log.Println("Listening...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func registerCleanUp(db *sql.DB) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range c {
			// todo: gracefully shut down http server and close db
			db.Close()

			os.Exit(0)
		}
	}()
}

type madminHandler struct {
	router *mux.Router

	userManager *userManager

	warehouse Warehouse
	database  *sql.DB
}

func (m *madminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func newMAdminHandler(db *sql.DB) *madminHandler {
	maHandler := &madminHandler{}

	maHandler.database = db
	maHandler.userManager = newUserManager(maHandler.database)
	maHandler.warehouse = NewWarehouse(maHandler.database)

	maHandler.router = mux.NewRouter()

	maHandler.router.HandleFunc("/data/stock/{....-..-..-..-......}", maHandler.stockItemHandler).Methods("GET", "DELETE")
	maHandler.router.HandleFunc("/data/stock/", maHandler.stockHandler).Methods("GET", "POST")

	return maHandler
}

func (m *madminHandler) stockHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		m.listStockHandler(w, r)
	case "POST":
		m.addStockHandler(w, r)
	default:
		respondMethodNotAllowed(w, r)
	}
}

func (m *madminHandler) stockItemHandler(w http.ResponseWriter, r *http.Request) {
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
func (m *madminHandler) listStockHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp = &CollectionResponseDTO{"List of existing stock items", make([]string, 0, m.warehouse.Size())}

		itemURL string
	)

	for _, item := range m.warehouse.Stock() {
		itemURL = fmt.Sprintf("/data/stock/%s", item.Name())
		resp.Urls = append(resp.Urls, itemURL)
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
func (m *madminHandler) getStockItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL
		_, id = path.Split(query.String())

		item, ok = m.warehouse.ReadStock(id)
	)
	if !ok {
		w.Header().Add("ContentLength", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := &StockDTO{}

	resp.ID = item.ID()
	resp.Name = item.Name()
	resp.Type = item.Type()
	if item.IsExpirable() {
		resp.ExpirationDate = item.ExpirationDate().String()
	}
	resp.MinQuantity = item.MinQuantity().String()
	resp.DistributorID = item.DistributorID()

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
func (m *madminHandler) addStockHandler(w http.ResponseWriter, r *http.Request) {
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
	id := m.warehouse.CreateStock(stockItem)

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(id)); err != nil {
		log.Printf("Error while writing response: %s", err)
	}
}

// Handler for DELETE /stock/<id>
//
// Removes the item with <id> from the warehouse.
func (m *madminHandler) removeStockItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL
		_, id = path.Split(query.String())
	)
	m.warehouse.DeleteStock(id)
	w.WriteHeader(http.StatusNoContent)
}
