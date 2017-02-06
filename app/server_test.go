package app

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
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

	requests := []struct {
		path   string
		body   string
		status int
	}{
		{"/data/stock/", `{"name": "Aspirin", "type": 0, "expirationDate": "2030-01-01T00:00:00.000Z"}`, http.StatusCreated},
		{"/data/stock/", `{"name": "Happy Doge","type": 1, "expirationDate": "2030-01-01T00:00:00.000Z", "minQuantity": "12.5", "distributor": "Happy Doge - Yakimovo"}`, http.StatusCreated},
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

		idBytes := make([]byte, 0, 20)
		if n, err := resp.Body.Read(idBytes); n != 20 || err != nil {
			t.Errorf("Error in response body: invalid ID size returned", n, err)
		}

		if result, err := regexp.Match("....-..-..-..-......", idBytes); result != true || err != nil {
			t.Errorf("Error in response body: invalid ID format returned", result, err)
		}
	}
}
