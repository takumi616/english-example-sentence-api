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

func TestFetchSingleSentence_handler(t *testing.T) {
	type expected struct {
		//Expected http status code
		statusCode int
		//Expected json response
		responseBody string
	}

	type testData struct {
		//Test sentence data
		sentence entity.Sentence
		//Returned error message
		errorMessage error
		//Expected result
		expected expected
	}

	//Prepare test cases
	testCases := map[string]testData{}
	//OK test case
	testCases["Ok"] = testData{
		sentence: entity.Sentence{
			SentenceID:   1,
			Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			Vocabularies: pq.StringArray{"build", "deployment", "container"},
			Created:      "2024-04-25 05:25:21.177099049 +0000 UTC m=+161.586254951",
			Updated:      "2024-04-25 05:25:21.177125382 +0000 UTC m=+161.586281242",
		},
		errorMessage: nil,
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/fetchsingle/ok_resp.json.golden",
		},
	}

	//Internal server error test case
	//Not exist a sentence specified by given sentenceID
	testCases["No rows"] = testData{
		sentence:     entity.Sentence{},
		errorMessage: errors.New("sql: no rows in result set"),
		expected: expected{
			statusCode:   http.StatusInternalServerError,
			responseBody: "../testhelper/golden/fetchsingle/no_rows_resp.json.golden",
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
			r := httptest.NewRequest(http.MethodGet, "/sentences/1", nil)

			//SentenceFetcherMock　mocks SentenceFetcher interface
			//which is used to call service package method
			moq := &SentenceFetcherMock{}
			moq.FetchSingleSentenceFunc = func(ctx context.Context, id string) (entity.Sentence, error) {
				return testdata.sentence, testdata.errorMessage
			}

			//Send http request
			sut := FetchSingleSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Compare http response body to expected result
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testdata.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testdata.expected.responseBody))
		})
	}
}
