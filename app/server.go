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
	"time"
)

type madminServer struct {
	*http.Server
}

func NewMAdminServer(port, dbPath string) *madminServer {
	database := newDB(dbPath)

	maUserManager := NewUserManager(database)
	maHandler := NewMAdminHandler(database)

	r := mux.NewRouter()
	r.Handle("/data/{path:.*}", authMiddleware(maHandler, maUserManager))
	r.Handle("/{path:.*}", authMiddleware(http.FileServer(http.Dir("static/")), maUserManager))

	registerCleanUp(database)

	return &madminServer{
		&http.Server{Addr: port, Handler: r},
	}
}

// todo: move functionalities to the (m madminServer) Shutdown() method
func registerCleanUp(db *sql.DB) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			// todo: gracefully shut down http server and close db
			db.Close()

			os.Exit(0)
		}
	}()
}

type madminHandler struct {
	router *mux.Router

	userManager UserManager

	warehouse Warehouse
	database  *sql.DB
}

func (m madminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func NewMAdminHandler(db *sql.DB) *madminHandler {
	maHandler := &madminHandler{}

	maHandler.database = db
	maHandler.userManager = NewUserManager(maHandler.database)
	maHandler.warehouse = NewWarehouse(maHandler.database)

	maHandler.router = mux.NewRouter()

	maHandler.router.HandleFunc("/data/stock/{id:....-..-..-..-......}", maHandler.stockItemHandler).Methods("GET", "DELETE", "PUT")
	maHandler.router.HandleFunc("/data/stock/", maHandler.stockHandler).Methods("GET", "POST")
	maHandler.router.HandleFunc("/data/stock/insufficient/", maHandler.insufficientStockHandler).Methods("GET")
	maHandler.router.HandleFunc("/data/stock/expiring/", maHandler.expiringStockHandler).Methods("GET")

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
	case "PUt":
		m.updateStockItemHandler(w, r)
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
		itemURL = fmt.Sprintf("/data/stock/%s", item.ID())
		resp.URLs = append(resp.URLs, itemURL)
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
		w.WriteHeader(http.StatusBadRequest)
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
	m.warehouse.CreateStock(stockItem)

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(stockItem.ID())); err != nil {
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

// Handler for PUT /stock/<id>
//
// UPDATE the item with <id> in the warehouse
// or return error if there is no such item.
func (m *madminHandler) updateStockItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL
		_, id = path.Split(query.String())
	)

	var (
		updateDto = &StockDTO{}

		decoder = json.NewDecoder(r.Body)
		err     = decoder.Decode(updateDto)
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error in unmarshaling request body: %s", err)
		return
	}
	defer r.Body.Close()

	stockItem, ok := m.warehouse.ReadStock(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}

	err = stockItem.Update(*updateDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	m.warehouse.UpdateStock(stockItem)

	w.WriteHeader(http.StatusAccepted)
}

// Handler for GET /stock/insufficient/
//
// Lists insufficient stock items
func (m *madminHandler) insufficientStockHandler(w http.ResponseWriter, r *http.Request) {

	stockItems := m.warehouse.Stock()
	insufficientStockItems := make([]string, 0, len(stockItems))

	for _, stock := range stockItems {
		if stock.Quantity().Cmp(stock.MinQuantity()) == -1 {

			itemURL := fmt.Sprintf("/data/stock/%s", stock.ID())
			insufficientStockItems = append(insufficientStockItems, itemURL)
		}
	}

	resp := &CollectionResponseDTO{"List of insufficient stock items", insufficientStockItems}

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

// TODO: add list of insufficient stock items for concrete distributor

// Handler for GET /stock/expiring/
//
// Lists expiring stock items
func (m *madminHandler) expiringStockHandler(w http.ResponseWriter, r *http.Request) {

	var (
		stockItems         = m.warehouse.Stock()
		expiringStockItems = make([]string, 0, len(stockItems))
		dateLimit          = time.Now().AddDate(0, 0, 7)
	)

	for _, stock := range stockItems {
		if stock.IsExpirable() {
			if !stock.ExpirationDate().Before(dateLimit) {
				itemURL := fmt.Sprintf("/data/stock/%s", stock.ID())
				expiringStockItems = append(expiringStockItems, itemURL)
			}
		}
	}

	resp := &CollectionResponseDTO{"List of existing stock items", expiringStockItems}

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
