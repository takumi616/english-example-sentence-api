package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/go-postgres-docker-restapi/entity"
	"github.com/takumi616/go-postgres-docker-restapi/testhelper"
)

func TestFetchSentenceList_handler(t *testing.T) {
	type expected struct {
		//Expected http status code
		statusCode int
		//Expected json response body
		responseBody string
	}

	type testData struct {
		//Test sentences data
		sentences []entity.Sentence
		//Expected result
		expected expected
	}

	//Prepare two test cases
	testCases := map[string]testData{}
	//OK test case
	testCases["Ok"] = testData{
		sentences: []entity.Sentence{
			{
				SentenceID:   1,
				Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
				Vocabularies: pq.StringArray{"build", "deployment", "container"},
				Created:      "2024-04-25 05:25:21.177099049 +0000 UTC m=+161.586254951",
				Updated:      "2024-04-25 05:25:21.177125382 +0000 UTC m=+161.586281242",
			},
			{
				SentenceID:   2,
				Body:         "The application communicates with the database server to retrieve and store data.",
				Vocabularies: pq.StringArray{"application", "store", "server"},
				Created:      "2024-04-25 05:28:38.567684293 +0000 UTC m=+358.976395166",
				Updated:      "2024-04-25 05:28:38.567821584 +0000 UTC m=+358.976532375",
			},
		},
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/fetchlist/ok_resp.json.golden",
		},
	}

	//Empty test case
	testCases["Empty"] = testData{
		sentences: []entity.Sentence{},
		expected: expected{
			statusCode:   http.StatusOK,
			responseBody: "../testhelper/golden/fetchlist/empty_resp.json.golden",
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
			r := httptest.NewRequest(http.MethodGet, "/sentences", nil)

			//SentenceFetcherMock　mocks SentenceFetcher interface
			//which is used to call service package method
			moq := &SentenceFetcherMock{}
			moq.FetchSentenceListFunc = func(ctx context.Context) ([]entity.Sentence, error) {
				return testdata.sentences, nil
			}

			//Send http request
			sut := FetchSentenceList{Service: moq}
			sut.ServeHTTP(w, r)

			//Compare http response body to expected result
			resp := w.Result()
			testhelper.CompareHTTPResponse(t, resp, testdata.expected.statusCode, testhelper.LoadJsonGoldenFile(t, testdata.expected.responseBody))
		})
	}
}
