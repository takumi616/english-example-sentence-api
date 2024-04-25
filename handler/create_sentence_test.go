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
		//Expected http response body
		responseBody string
		//Expected returned rows affected number
		rowsAffectedNumber int64
		//Expected error message used in error cases
		errorMessage string
	}

	type testData struct {
		//Http request body
		requestBody string
		//Expected data
		expected expected
	}

	//Prepare test cases
	testCases := map[string]testData{}

	//Ok case
	testCases["Ok"] = testData{
		requestBody: "../testhelper/golden/create/ok_req.json.golden",
		expected: expected{
			statusCode:         http.StatusOK,
			responseBody:       "../testhelper/golden/create/ok_resp.json.golden",
			rowsAffectedNumber: 1,
			errorMessage:       "",
		},
	}

	//Bad request case
	testCases["Bad request"] = testData{
		requestBody: "../testhelper/golden/create/badreq_req.json.golden",
		expected: expected{
			statusCode:         http.StatusBadRequest,
			responseBody:       "../testhelper/golden/create/badreq_resp.json.golden",
			rowsAffectedNumber: 0,
			errorMessage:       "Key: 'Vocabularies' Error:Field validation for 'Vocabularies' failed on the 'required' tag\nKey: 'Body' Error:Field validation for 'Body' failed on the 'required' tag",
		},
	}

	//Case http internal server error
	testCases["Internal server error"] = testData{
		requestBody: "../testhelper/golden/create/internalServErr_req.json.golden",
		expected: expected{
			statusCode:         http.StatusInternalServerError,
			responseBody:       "../testhelper/golden/create/internalServErr_resp.json.golden",
			rowsAffectedNumber: 0,
			errorMessage:       "pq: value too long for type character varying(120)",
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
				switch testcase {
				//Expect request body is set to sentence entity correctly
				case "Ok":
					return testdata.expected.rowsAffectedNumber, nil
				case "Internal server error":
					return testdata.expected.rowsAffectedNumber, errors.New(testdata.expected.errorMessage)
				//Not find test case
				default:
					return 0, errors.New("Failed to find test case")
				}
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
