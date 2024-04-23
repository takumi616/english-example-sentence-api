package main_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
)

func TestSetUpRouting(t *testing.T) {
	//Set up health check endpoint
	router := chi.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	//Create http request and response writer
	r := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	//Send http request
	router.ServeHTTP(w, r)

	//Get http response body
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })
	actualJson, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read http response body: %v", err)
	}

	//Unmarshal json data into go struct
	var expectedBody, actualBody any
	expectedJson := []byte(`{"status": "ok"}`)
	if err = json.Unmarshal(actualJson, &actualBody); err != nil {
		t.Errorf("Failed to unmarshal json data: %v", err)
	}
	if err = json.Unmarshal(expectedJson, &expectedBody); err != nil {
		t.Errorf("Failed to unmarshal json data: %v", err)
	}

	//Compare diff
	if diff := cmp.Diff(expectedBody, actualBody); diff != "" {
		t.Errorf("Some differences are found: (-expected +actual)\n%s", diff)
	}
}
