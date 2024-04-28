package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takumi616/go-postgres-docker-restapi/testhelper"
)

func TestDeleteSentence_handler(t *testing.T) {
	type expected struct {
		//Expected status code
		statusCode int
		//Expected http response body
		response string
	}

	type testdata struct {
		//Returned rows affected number
		rowsAffectedNumber int64
		//Returned error message
		errorMessage error
		//Expected data
		expected expected
	}

	//Prepare test cases
	testCases := map[string]testdata{}
	//Test case Ok
	testCases["Ok"] = testdata{
		rowsAffectedNumber: 1,
		errorMessage:       nil,
		expected: expected{
			statusCode: http.StatusOK,
			response:   "../testhelper/golden/delete/ok_resp.json.golden",
		},
	}

	//Test case Internal server error
	//Not exist a sentence specified by given sentenceID
	testCases["No rows"] = testdata{
		rowsAffectedNumber: 0,
		errorMessage:       errors.New("sql: no rows in result set"),
		expected: expected{
			statusCode: http.StatusInternalServerError,
			response:   "../testhelper/golden/delete/no_rows_resp.json.golden",
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
			r := httptest.NewRequest(http.MethodDelete, "/sentences/1", nil)

			//SentenceDeleterMock　mocks SentenceDeleter interface
			//which is used to call service package method
			moq := &SentenceDeleterMock{}
			moq.DeleteSentenceFunc = func(ctx context.Context, id string) (int64, error) {
				return testdata.rowsAffectedNumber, testdata.errorMessage
			}

			//Send http request
			sut := DeleteSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testdata.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testdata.expected.response))
		})
	}
}
