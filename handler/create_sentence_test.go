package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/takumi616/generate-example/entity"
	"github.com/takumi616/generate-example/testhelper"
)

func TestCreateNewSentence(t *testing.T) {
	type expected struct {
		//Expected http status code
		statusCode int
		//Expected http response body
		responseBody string
		//Expected returned sentenceID
		sentenceID int
	}

	type testData struct {
		//Http request body
		requestBody string
		//Expected data
		expected expected
	}

	//Prepare test cases
	testCases := map[string]testData{}

	//Case ok
	testCases["ok"] = testData{
		requestBody: "../testhelper/golden/create/ok_req.json.golden",
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/create/ok_resp.json.golden",
			sentenceID:   6,
		},
	}

	//Case http bad request
	testCases["badreq"] = testData{
		requestBody: "../testhelper/golden/create/badreq_req.json.golden",
		expected: expected{
			statusCode:   http.StatusBadRequest,
			responseBody: "../testhelper/golden/create/badreq_resp.json.golden",
			sentenceID:   0,
		},
	}

	//Case http internal server error
	testCases["internalServErr"] = testData{
		requestBody: "../testhelper/golden/create/internalServErr_req.json.golden",
		expected: expected{
			statusCode:   http.StatusInternalServerError,
			responseBody: "../testhelper/golden/create/internalServErr_resp.json.golden",
			sentenceID:   0,
		},
	}

	for name, testData := range testCases {
		testData := testData

		//Execute as subtest
		t.Run(name, func(t *testing.T) {
			//Run in parallel
			t.Parallel()

			//Prepare http request and response writer
			r := httptest.NewRequest(http.MethodPost, "/sentences", bytes.NewReader(testhelper.LoadJsonGoldenFile(t, testData.requestBody)))
			w := httptest.NewRecorder()

			//Create service layer mock
			moq := &SentenceCreaterMock{}
			moq.CreateNewSentenceFunc = func(ctx context.Context, vocabularies []string, body string) (int, error) {
				//Case bad request and internalServErr do not include in this mock func
				//because validator and json decoder will catch request body error
				//before calling mock func

				sentence := &entity.Sentence{
					Body:         body,
					Vocabularies: vocabularies,
					Created:      time.Now().String(),
					Updated:      time.Now().String(),
				}

				//Check if fields of sentence entity are set correctly
				if sentence.Body != "" && len(sentence.Vocabularies) != 0 && sentence.Created != "" && sentence.Updated != "" {
					return testData.expected.sentenceID, nil
				}

				return 0, errors.New("Unexpected Error. Failed to set request body to Go struct field.")
			}

			//Send http request
			sut := CreateSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testData.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testData.expected.responseBody))
		})
	}
}
