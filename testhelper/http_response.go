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
	actualResponseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to get http response body: %v", err)
	}

	//Compare http status code to expected
	if response.StatusCode != expectedStatusCode {
		t.Fatalf("Failed to get expected http status code. Expected:%d Actual: %d", expectedStatusCode, response.StatusCode)
	}

	if len(actualResponseBody) == 0 && len(expectedResponseBody) == 0 {
		//Not need to call compareResponseBody()
		//Because the length of response bodies are empty
		return
	}

	//Unmarshal json data into go struct
	var expectedBody, actualBody any
	if err := json.Unmarshal(expectedResponseBody, &expectedBody); err != nil {
		t.Fatalf("Failed to unmarshal http response body: %v", err)
	}
	if err := json.Unmarshal(actualResponseBody, &actualBody); err != nil {
		t.Fatalf("Failed to unmarshal http response body: %v", err)
	}

	//Compare diff
	if diff := cmp.Diff(expectedBody, actualBody); diff != "" {
		t.Errorf("Some differences are found: (-expected +actual)\n%s", diff)
	}
}
