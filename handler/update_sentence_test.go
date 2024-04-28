package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takumi616/go-postgres-docker-restapi/testhelper"
)

func TestUpdateSentence_handler(t *testing.T) {
	type expected struct {
		//Expected http status code
		statusCode int
		//File url of expected json response body
		responseBody string
	}

	type testdata struct {
		//Http request body
		requestBody string
		//Returned rows affected number
		rowsAffectedNumber int64
		//Returned error message
		errorMessage error
		//expected result
		expected expected
	}

	//Prepare test cases
	testCases := map[string]testdata{}

	//Test case Ok
	testCases["Ok"] = testdata{
		requestBody:        "../testhelper/golden/update/ok_req.json.golden",
		rowsAffectedNumber: 1,
		errorMessage:       nil,
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/update/ok_resp.json.golden",
		},
	}

	//Test case Internal server error
	//Not exist a sentence specified by given sentenceID
	testCases["No rows"] = testdata{
		requestBody:        "../testhelper/golden/update/no_rows_req.json.golden",
		rowsAffectedNumber: 0,
		errorMessage:       errors.New("sql: no rows in result set"),
		expected: expected{
			statusCode:   http.StatusInternalServerError,
			responseBody: "../testhelper/golden/update/no_rows_resp.json.golden",
		},
	}

	//Test case Bad request
	testCases["Bad request"] = testdata{
		requestBody:        "../testhelper/golden/update/badreq_req.json.golden",
		rowsAffectedNumber: 0,
		errorMessage:       errors.New("Key: 'Body' Error:Field validation for 'Body' failed on the 'required' tag"),
		expected: expected{
			statusCode:   http.StatusBadRequest,
			responseBody: "../testhelper/golden/update/badreq_resp.json.golden",
		},
	}

	for testcase, testdata := range testCases {
		testdata := testdata
		testcase := testcase
		//Execute as parallel tests
		//Run runs function as a subtest of t called name n(first parameter of Run)
		//It runs function in a separate goroutine and blocks
		//until this function returns or calls t.Parallel to become a parallel test
		t.Run(testcase, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//Create test http request and response writer
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/sentences/1", bytes.NewReader(testhelper.LoadJsonGoldenFile(t, testdata.requestBody)))

			//SentenceUpdaterMock　mocks SentenceUpdater interface
			//which is used to call service package method
			moq := &SentenceUpdaterMock{}
			moq.UpdateSentenceFunc = func(ctx context.Context, id string, body string) (int64, error) {
				return testdata.rowsAffectedNumber, testdata.errorMessage
			}

			//Send http request
			sut := UpdateSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testdata.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testdata.expected.responseBody))
		})
	}
}
