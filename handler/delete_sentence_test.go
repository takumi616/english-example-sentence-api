package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/generate-example/entity"
	"github.com/takumi616/generate-example/testhelper"
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

	//Prepare two test cases
	testCases := map[string]testData{}
	//OK
	testCases["ok"] = testData{
		sentence: entity.Sentence{
			SentenceID:   6,
			Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			Vocabularies: pq.StringArray{"build", "deployment", "container"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expectedStatusCode: http.StatusOK,
		expectedResponse:   "../testhelper/golden/delete/ok_resp.json.golden",
	}

	//Data does not exist
	testCases["error"] = testData{
		sentence: entity.Sentence{
			SentenceID:   5,
			Body:         "The application communicates with the database server to retrieve and store data.",
			Vocabularies: pq.StringArray{"application", "store", "server"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
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
			r := httptest.NewRequest(http.MethodDelete, "/sentences/6", nil)

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
