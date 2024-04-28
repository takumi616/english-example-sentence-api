package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takumi616/generate-example/testhelper"
)

func TestCreateNewSentence_handler(t *testing.T) {
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

	testCases["Ok"] = testdata{
		requestBody:        "../testhelper/golden/create/ok_req.json.golden",
		rowsAffectedNumber: 1,
		errorMessage:       nil,
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/create/ok_resp.json.golden",
		},
	}

	//Bad request case
	testCases["Bad request"] = testdata{
		requestBody:        "../testhelper/golden/create/badreq_req.json.golden",
		rowsAffectedNumber: 0,
		errorMessage:       errors.New("Key: 'Vocabularies' Error:Field validation for 'Vocabularies' failed on the 'required' tag\nKey: 'Body' Error:Field validation for 'Body' failed on the 'required' tag"),
		expected: expected{
			statusCode:   http.StatusBadRequest,
			responseBody: "../testhelper/golden/create/badreq_resp.json.golden",
		},
	}

	//Case http internal server error
	testCases["Internal server error"] = testdata{
		requestBody:        "../testhelper/golden/create/internalServErr_req.json.golden",
		rowsAffectedNumber: 0,
		errorMessage:       errors.New("pq: value too long for type character varying(120)"),
		expected: expected{
			statusCode:   http.StatusInternalServerError,
			responseBody: "../testhelper/golden/create/internalServErr_resp.json.golden",
		},
	}

	for testcase, testdata := range testCases {
		testdata := testdata
		testcase := testcase

		//Execute as subtest
		t.Run(testcase, func(t *testing.T) {
			//Run in parallel
			t.Parallel()

			//Prepare http request and response writer
			r := httptest.NewRequest(http.MethodPost, "/sentences", bytes.NewReader(testhelper.LoadJsonGoldenFile(t, testdata.requestBody)))
			w := httptest.NewRecorder()

			//Create service layer mock
			moq := &SentenceCreaterMock{}
			moq.CreateNewSentenceFunc = func(ctx context.Context, vocabularies []string, body string) (int64, error) {
				return testdata.rowsAffectedNumber, testdata.errorMessage
			}

			//Send http request
			sut := CreateSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testdata.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testdata.expected.responseBody))
		})
	}
}
