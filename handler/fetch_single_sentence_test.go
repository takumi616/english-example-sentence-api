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

func TestFetchSingleSentence(t *testing.T) {
	type expected struct {
		//Expected http status code
		statusCode int
		//Expected json response
		responseBody string
	}

	type testData struct {
		//Test sentence data
		sentence entity.Sentence
		//Expected result
		expected expected
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
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/fetchsingle/ok_resp.json.golden",
		},
	}

	//Data does not exist (sentenceID 5)
	testCases["error"] = testData{
		sentence: entity.Sentence{
			SentenceID:   5,
			Body:         "The application communicates with the database server to retrieve and store data.",
			Vocabularies: pq.StringArray{"application", "store", "server"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expected: expected{
			statusCode:   http.StatusInternalServerError,
			responseBody: "../testhelper/golden/fetchsingle/err_resp.json.golden",
		},
	}

	for n, testData := range testCases {
		testData := testData
		//Execute as parallel tests
		//Run runs function as a subtest of t called name n(first parameter of Run)
		//It runs function in a separate goroutine and blocks
		//until this function returns or calls t.Parallel to become a parallel test
		t.Run(n, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//Create test http request and response writer
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/sentences/6", nil)

			//SentenceFetcherMock　mocks SentenceFetcher interface
			//which is used to call service package method
			moq := &SentenceFetcherMock{}
			moq.FetchSingleSentenceFunc = func(ctx context.Context, id string) (entity.Sentence, error) {
				if testData.expected.statusCode == http.StatusOK {
					return testData.sentence, nil
				}
				return entity.Sentence{}, errors.New("sql: no rows in result set")
			}

			//Send http request
			sut := FetchSingleSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Compare http response body to expected result
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testData.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testData.expected.responseBody))
		})
	}
}
