package testhelper

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func CheckOutHTTPResponse(t *testing.T, got *http.Response, wantStatus int, wantBody []byte) {
	//Helper marks the calling function as a test helper function.
	//When printing file and line information, this function will be skipped.
	t.Helper()

	//Read response body
	t.Cleanup(func() { _ = got.Body.Close() })
	gotBody, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}

	//Check http status code
	if got.StatusCode != wantStatus {
		t.Fatalf("wantStatus: %d,  gotStatus: %d", wantStatus, got.StatusCode)
	}

	if len(gotBody) == 0 && len(wantBody) == 0 {
		//Not need to call AssertJSON()
		//Because the length of want body and got body are empty
		return
	}

	//Check out json http response body
	CheckOutJSON(t, wantBody, gotBody)
}

func CheckOutJSON(t *testing.T, wantBody, gotBody []byte) {
	//Helper marks the calling function as a test helper function.
	//When printing file and line information, this function will be skipped.
	t.Helper()

	//Convert json data into golang's data type
	var wantResponse, gotResponse any
	if err := json.Unmarshal(wantBody, &wantResponse); err != nil {
		t.Fatalf("Failed to unmarshal want %q: %v", wantBody, err)
	}
	if err := json.Unmarshal(gotBody, &gotResponse); err != nil {
		t.Fatalf("Failed to unmarshal got %q: %v", gotBody, err)
	}

	//Compare diff
	if diff := cmp.Diff(gotResponse, wantResponse); diff != "" {
		t.Errorf("Some differences are found: (-got +want)\n%s", diff)
	}
}
