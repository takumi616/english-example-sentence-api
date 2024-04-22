package testhelper

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func CompareHTTPResponse(t *testing.T, response *http.Response, expectedStatusCode int, expectedResponseBody []byte) {
	//Helper marks the calling function as a test helper function.
	//When printing file and line information, this function will be skipped.
	t.Helper()

	//Read response body
	t.Cleanup(func() { _ = response.Body.Close() })
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	//Compare http status code to expected
	if response.StatusCode != expectedStatusCode {
		t.Fatalf("Failed to get expected http status code. Expected:%d Result: %d", expectedStatusCode, response.StatusCode)
	}

	if len(responseBody) == 0 && len(expectedResponseBody) == 0 {
		//Not need to call compareResponseBody()
		//Because the length of response bodies are empty
		return
	}

	//Unmarshal http response body
	var expectedBody, resultBody any
	if err := json.Unmarshal(expectedResponseBody, &expectedBody); err != nil {
		t.Fatalf("Failed to unmarshal expectedResponseBody: %v", err)
	}
	if err := json.Unmarshal(responseBody, &resultBody); err != nil {
		t.Fatalf("Failed to unmarshal responseBody: %v", err)
	}

	//Compare diff
	if diff := cmp.Diff(expectedBody, resultBody); diff != "" {
		t.Errorf("Some differences are found: (-expected +result)\n%s", diff)
	}
}
