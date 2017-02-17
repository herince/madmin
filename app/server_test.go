package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"
)

func buildURL(base, path string) string {
	return fmt.Sprintf("%s%s", base, path)
}

func TestAddStockPOSTRequest(t *testing.T) {
	var (
		dbPath        = "./test_database.sqlite"
		database      = newDB(dbPath)
		madminHandler = newMAdminHandler(database)
		s             = httptest.NewServer(madminHandler)
	)
	defer database.Close()
	defer os.Remove("./test_database.sqlite")
	defer s.Close()

	requests := []struct {
		path   string
		body   string
		status int
	}{
		{"/data/stock/", `{"name": "Aspirin", "type": 0, "expirationDate": "2030-01-01T00:00:00.000Z"}`, http.StatusCreated},
		{"/data/stock/", `{"name": "Happy Doge","type": 1, "expirationDate": "2030-01-01T00:00:00.000Z", "minQuantity": "12.5", "distributor": "Happy Doge - Yakimovo"}`, http.StatusCreated},
	}

	for _, req := range requests {
		postURL := buildURL(s.URL, req.path)

		resp, err := http.Post(postURL, "application/json", bytes.NewReader([]byte(req.body)))
		if err != nil {
			t.Fatalf("Error sending POST request: %s", err)
		}

		if resp.StatusCode != req.status {
			t.Errorf("Expected %d but got %d for adding stock with body: %s",
				req.status, resp.StatusCode, req.body)
		}

		idBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if len(idBytes) != 36 || err != nil {
			t.Errorf("Error in response body. Invalid ID size returned. Got %d with error %s", len(idBytes), err)
		}

		if result, err := regexp.Match("........-....-....-....-............", idBytes); result != true || err != nil {
			t.Errorf("Error in response body: invalid ID format returned. Got %t with error %s for matching question", result, err)
		}
	}
}

func TestValidRemoveStockDELETERequest(t *testing.T) {
	var (
		dbPath        = "./test_database.sqlite"
		database      = newDB(dbPath)
		madminHandler = newMAdminHandler(database)
		s             = httptest.NewServer(madminHandler)
	)
	defer database.Close()
	defer os.Remove("./test_database.sqlite")
	defer s.Close()

	addReqest := struct {
		path   string
		body   string
		status int
	}{
		"/data/stock/", `{"name": "Happy Doge","type": 1, "expirationDate": "2030-01-01T00:00:00.000Z", "minQuantity": "12.5", "distributor": "Happy Doge - Yakimovo"}`, http.StatusCreated,
	}

	postURL := buildURL(s.URL, addReqest.path)

	addResponse, err := http.Post(postURL, "application/json", bytes.NewReader([]byte(addReqest.body)))
	if err != nil {
		t.Fatalf("Error sending POST request: %s", err)
	}

	idBytes, err := ioutil.ReadAll(addResponse.Body)
	addResponse.Body.Close()
	if len(idBytes) != 36 || err != nil {
		t.Errorf("Error in adding new stock for the test. Invalid ID size returned. Got %d with error %s", len(idBytes), err)
	}

	var (
		client http.Client

		deleteRequestPath   = fmt.Sprintf("/data/stock/%s", idBytes)
		deleteRequestStatus = http.StatusNoContent

		deleteURLString = buildURL(s.URL, deleteRequestPath)
	)
	deleteRequestURL, err := url.Parse(deleteURLString)
	if err != nil {
		t.Fatalf("Error in building request URL. %s", err)
	}

	deleteRequest := &http.Request{Method: "DELETE", URL: deleteRequestURL}
	deleteResponse, err := client.Do(deleteRequest)
	if err != nil {
		t.Fatalf("Error sending DELETE request: %s", err)
	}

	if deleteResponse.StatusCode != deleteRequestStatus {
		t.Errorf("Expected %d but got %d for deleting stock with id: %s",
			deleteRequestStatus, deleteResponse.StatusCode, idBytes)
	}
}
