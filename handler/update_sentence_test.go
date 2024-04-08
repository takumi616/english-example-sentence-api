package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/english-example-sentence-api/entity"
	"github.com/takumi616/english-example-sentence-api/testhelper"
)

func TestUpdateSentence(t *testing.T) {
	type testData struct {
		//Json request body
		requestBody string
		//Test sentence entity
		sentence entity.Sentence
		//Expected http status code
		expectedStatusCode int
		//Expected json response
		expectedResponse string
	}

	//Prepare three test cases
	testCases := map[string]testData{}
	//OK
	testCases["ok"] = testData{
		requestBody: "../testhelper/golden/update/ok_req.json.golden",
		sentence: entity.Sentence{
			SentenceID:   9,
			Body:         "Updated sentence.",
			Vocabularies: pq.StringArray{"docker", "is", "used"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expectedStatusCode: http.StatusOK,
		expectedResponse:   "../testhelper/golden/update/ok_resp.json.golden",
	}
	//Bad request
	testCases["badRequest"] = testData{
		requestBody:        "../testhelper/golden/update/badreq_req.json.golden",
		sentence:           entity.Sentence{},
		expectedStatusCode: http.StatusBadRequest,
		expectedResponse:   "../testhelper/golden/update/badreq_resp.json.golden",
	}
	//Data not found
	testCases["empty"] = testData{
		requestBody:        "../testhelper/golden/update/empty_req.json.golden",
		sentence:           entity.Sentence{},
		expectedStatusCode: http.StatusInternalServerError,
		expectedResponse:   "../testhelper/golden/update/empty_resp.json.golden",
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
			r := httptest.NewRequest(http.MethodPut, "/sentences/9", bytes.NewReader(testhelper.LoadJsonGoldenFile(t, testData.requestBody)))

			//SentenceUpdaterMock　mocks SentenceUpdater interface
			//which is used to call service package method
			moq := &SentenceUpdaterMock{}
			moq.UpdateSentenceFunc = func(ctx context.Context, id string, body string) (entity.Sentence, error) {
				switch testData.expectedStatusCode {
				case http.StatusOK:
					return testData.sentence, nil
				case http.StatusBadRequest:
					return testData.sentence, errors.New("Key: 'Body' Error:Field validation for 'Body' failed on the 'required' tag")
				case http.StatusInternalServerError:
					return testData.sentence, errors.New("sql: no rows in result set")
				default:
					return testData.sentence, nil
				}
			}

			//Send http request
			sut := UpdateSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CheckOutHTTPResponse(t, resp, testData.expectedStatusCode, testhelper.LoadJsonGoldenFile(t, testData.expectedResponse))
		})
	}
}
