package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/english-example-sentence-api/entity"
	"github.com/takumi616/english-example-sentence-api/testhelper"
)

func TestDeleteSentence(t *testing.T) {
	type testData struct {
		//Test sentence entity
		sentence entity.Sentence
		//Expected http status code
		expectedStatusCode int
		//Expected json response data
		expectedResponse string
	}

	//Prepare two test case (ok and error)
	testCases := map[string]testData{}
	//Case ok
	testCases["ok"] = testData{
		sentence: entity.Sentence{
			SentenceID:   1,
			Body:         "This is a test sentence.",
			Vocabularies: pq.StringArray{"This", "is", "test"},
			Created:      "2024-03-28T14:15:21.574757Z",
			Updated:      "2024-03-28T14:15:21.574758Z",
		},
		expectedStatusCode: http.StatusOK,
		expectedResponse:   "../testhelper/golden/delete/ok_resp.json.golden",
	}
	//Case error
	testCases["error"] = testData{
		sentence:           entity.Sentence{},
		expectedStatusCode: http.StatusInternalServerError,
		expectedResponse:   "../testhelper/golden/delete/err_resp.json.golden",
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
			r := httptest.NewRequest(http.MethodDelete, "/sentences/1", nil)

			//SentenceDeleterMock　mocks SentenceDeleter interface
			//which is used to call service package method
			moq := &SentenceDeleterMock{}
			moq.DeleteSentenceFunc = func(ctx context.Context, id string) (int, error) {
				if testData.expectedStatusCode == http.StatusOK {
					return testData.sentence.SentenceID, nil
				}
				return 0, errors.New("sql: no rows in result set")
			}

			//Send http request
			sut := DeleteSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CheckOutHTTPResponse(t, resp, testData.expectedStatusCode, testhelper.LoadJsonGoldenFile(t, testData.expectedResponse))
		})
	}
}
