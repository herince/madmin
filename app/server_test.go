package app

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func buildUrl(base, path string) string {
	return fmt.Sprintf("%s%s", base, path)
}

func TestAddStockPOSTRequest(t *testing.T) {
	var (
		madminHandler = NewMAdminHandler()
		s             = httptest.NewServer(madminHandler)
	)
	defer s.Close()

	/* todo: add test cases for different sizes of information */
	requests := []struct {
		path   string
		body   string
		status int
	}{
		{"/data/stock/", `{"name": "Aspirin", "type": 0}`, http.StatusCreated},
		{"/data/stock/", `{"name": "Happy Doge","type": 2, "expirationDate": "2030-01-01T00:00:00.00Z", "minQuantity": "12.5", "distributor": "Happy Doge - Yakimovo"}`, http.StatusCreated},
	}

	for _, req := range requests {
		postUrl := buildUrl(s.URL, req.path)

		resp, err := http.Post(postUrl, "application/json", bytes.NewReader([]byte(req.body)))
		if err != nil {
			t.Fatalf("Error seding POST request: %s", err)
		}

		if resp.StatusCode != req.status {
			t.Errorf("Expected %d but got %d for adding stock with body: %s",
				req.status, resp.StatusCode, req.body)
		}
	}
}
