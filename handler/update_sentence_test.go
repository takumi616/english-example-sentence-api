package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/generate-example/entity"
	"github.com/takumi616/generate-example/testhelper"
)

func TestUpdateSentence(t *testing.T) {

	type expected struct {
		//Expected http status code
		statusCode int
		//File url of expected json response body
		responseBody string
	}

	type testData struct {
		//Json request body
		requestBody string
		//Test sentence entity
		sentence entity.Sentence
		//expected result
		expected expected
	}

	//Prepare three test cases
	testCases := map[string]testData{}
	//OK
	testCases["ok"] = testData{
		requestBody: "../testhelper/golden/update/ok_req.json.golden",
		sentence: entity.Sentence{
			SentenceID:   6,
			Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			Vocabularies: pq.StringArray{"build", "deployment", "container"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/update/ok_resp.json.golden",
		},
	}

	//Bad request
	testCases["badRequest"] = testData{
		requestBody: "../testhelper/golden/update/badreq_req.json.golden",
		sentence: entity.Sentence{
			SentenceID:   6,
			Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			Vocabularies: pq.StringArray{"build", "deployment", "container"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expected: expected{
			statusCode:   http.StatusBadRequest,
			responseBody: "../testhelper/golden/update/badreq_resp.json.golden",
		},
	}

	//Data does not exist
	testCases["error"] = testData{
		requestBody: "../testhelper/golden/update/no_rows_req.json.golden",
		sentence: entity.Sentence{
			SentenceID:   5,
			Body:         "The application communicates with the database server to retrieve and store data.",
			Vocabularies: pq.StringArray{"application", "store", "server"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expected: expected{
			statusCode:   http.StatusInternalServerError,
			responseBody: "../testhelper/golden/update/no_rows_resp.json.golden",
		},
	}

	for name, testData := range testCases {
		testData := testData
		//Execute as parallel tests
		//Run runs function as a subtest of t called name n(first parameter of Run)
		//It runs function in a separate goroutine and blocks
		//until this function returns or calls t.Parallel to become a parallel test
		t.Run(name, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//Create test http request and response writer
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/sentences/6", bytes.NewReader(testhelper.LoadJsonGoldenFile(t, testData.requestBody)))

			//SentenceUpdaterMock　mocks SentenceUpdater interface
			//which is used to call service package method
			moq := &SentenceUpdaterMock{}
			moq.UpdateSentenceFunc = func(ctx context.Context, id string, body string) (int64, error) {
				if testData.expected.statusCode == http.StatusOK {
					testData.sentence.Body = body
					return testData.sentence.SentenceID, nil
				} else {
					return 0, errors.New("sql: no rows in result set")
				}
			}

			//Send http request
			sut := UpdateSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testData.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testData.expected.responseBody))
		})
	}
}
