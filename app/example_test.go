package app_test

//
//import (
//	"database/sql"
//	"net/http"
//	"fmt"
//	"bytes"
//	"io"
//	"io/ioutil"
//	"os"
//	"time"
//	"encoding/json"
//
//	_ "github.com/mattn/go-sqlite3" // needed for for working with sqlite3
//
//	"github.com/herince/madmin/app"
//)

//func recreateDB(dbPath string) {
//	os.Remove(dbPath)
//	db, err := sql.Open("sqlite3", dbPath)
//	if err != nil {
//		panic(err)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		panic(err)
//	}
//
//	db.Close()
//
//	return
//}
//
//func removeExampleDB(dbPath string) {
//	os.Remove(dbPath)
//}
//
//func buildUrl(port, path string) string {
//	return fmt.Sprintf("%s%s%s", "http://localhost", port, path)
//}
//
//func jsonReader(json string) io.Reader {
//	return bytes.NewReader([]byte(json))
//}
//
//func Example_first() {
//	var (
//		dbPath = "example.sqlite"
//		port = ":1234"
//
//		bodyTypeJSON = "application/json"
//	)
//	recreateDB(dbPath)
//
//
//	go func() {
//		app.Init(port, dbPath)
//		// ??? stop server at some point ???
//
//		removeExampleDB(dbPath)
//	}()
//
//	duration, err := time.ParseDuration("2s")
//	time.Sleep(duration)
//
//	jsonBody := `{
//		"name": "Aspirin",
//		"type": 0,
//		"expirationDate": "2030-01-01T00:00:00.000Z",
//		"quantity" : "2.0",
//		"minQuantity": "3.0"
//	}`
//	resp, err := http.Post(buildUrl(port, "/data/stock/"), bodyTypeJSON, jsonReader(jsonBody))
//
//	if err != nil {
//		panic(err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		panic(fmt.Sprintf("could not add stock, code: %d", resp.StatusCode))
//	}
//
//	idBytes, err := ioutil.ReadAll(resp.Body)
//	resp.Body.Close()
//
//	fmt.Printf("ID of new stock: %s\n", string(idBytes))
//}
//
//
//func Example_second() {
//	var (
//		dbPath = "example.sqlite"
//		port = ":1234"
//
//		bodyTypeJSON = "application/json"
//	)
//	recreateDB(dbPath)
//
//
//	go func() {
//		app.Init(port, dbPath)
//		// ??? stop server at some point ???
//
//		removeExampleDB(dbPath)
//	}()
//
//	duration, err := time.ParseDuration("2s")
//	time.Sleep(duration)
//
//	jsons := []string {
//		/* already expired */
//		`{
//			"name": "Medicine1",
//			"type": 0,
//			"expirationDate": "2017-02-17T00:00:00.000Z",
//			"quantity" : "2.0"
//		}`,
//		/* expires in the next 7 days */
//		`{
//			"name": "Medicine2",
//			"type": 0,
//			"expirationDate": "2017-02-20T00:00:00.000Z",
//			"quantity" : "1.0"
//		}`,
//		/* not going to expire in the next 7 days */
//		`{
//			"name": "Medicine3",
//			"type": 0,
//			"expirationDate": "2017-04-20T00:00:00.000Z",
//			"quantity" : "1.0"
//		}`,
//	}
//
//	for _, jsonBody := range jsons {
//		resp, err := http.Post(buildUrl(port, "/data/stock/"), bodyTypeJSON, jsonReader(jsonBody))
//		if err != nil {
//			panic(err)
//		}
//
//		if resp.StatusCode != http.StatusCreated {
//			panic("could not add stock")
//		}
//	}
//
//	resp, err := http.Get(buildUrl(port, "/data/stock/expiring/"))
//	if err != nil {
//		panic(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		panic("could not get stock")
//	}
//
//	expiringItemsCollectionBytes, err := ioutil.ReadAll(resp.Body)
//	resp.Body.Close()
//
//	var expiringDTOs *app.CollectionResponseDTO
//	err = json.Unmarshal(expiringItemsCollectionBytes, expiringDTOs)
//	if err != nil {
//		panic(err)
//	}
//
//	var(
//		expiringStockItem *app.StockDTO
//		expiring = make(map[string]app.StockDTO)
//	)
//	for _, url := range expiringDTOs.URLs {
//		resp, err = http.Get(buildUrl(port, url))
//		if err != nil {
//			panic(err)
//		}
//
//		if resp.StatusCode != http.StatusOK {
//			panic("could not get stock")
//		}
//
//		expiringItemBytes, err := ioutil.ReadAll(resp.Body)
//		resp.Body.Close()
//
//		err = json.Unmarshal(expiringItemBytes, expiringStockItem)
//		if err != nil {
//			panic(err)
//		}
//
//		expiring[expiringStockItem.ID] = *expiringStockItem
//	}
//
//	for _, stockItem := range expiring {
//		fmt.Printf(`
//			Item %s expires on %s`,
//			stockItem.Name,
//			stockItem.ExpirationDate)
//	}
//	// Output:
//	// Item Medicine1 expires on 2017-02-17T00:00:00.000Z
//	// Item Medicine2 expires on 2017-02-20T00:00:00.000Z
//	// Item Medicine3 expires on 2017-04-20T00:00:00.000Z
//}
