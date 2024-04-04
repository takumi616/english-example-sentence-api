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

func TestFetchSingleSentence(t *testing.T) {
	type wantResult struct {
		statusCode   int
		responseFile string
	}

	type testCase struct {
		sentence entity.Sentence
		want     wantResult
	}

	//Prepare two test case (ok and empty response data pattern)
	testCases := map[string]testCase{}
	//Case ok
	testCases["ok"] = testCase{
		sentence: entity.Sentence{
			SentenceID:   1,
			Body:         "This is a test sentence.",
			Vocabularies: pq.StringArray{"This", "is", "test"},
			Created:      "2024-03-28T14:15:21.574757Z",
			Updated:      "2024-03-28T14:15:21.574758Z",
		},
		want: wantResult{
			statusCode:   http.StatusOK,
			responseFile: "../testhelper/golden/fetchsingle/ok_resp.json.golden",
		},
	}
	//Case error
	testCases["error"] = testCase{
		sentence: entity.Sentence{},
		want: wantResult{
			statusCode:   http.StatusInternalServerError,
			responseFile: "../testhelper/golden/fetchsingle/err_resp.json.golden",
		},
	}

	for n, testCase := range testCases {
		testCase := testCase
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
			r := httptest.NewRequest(http.MethodGet, "/sentences/2", nil)

			//SentenceFetcherMock　mocks SentenceFetcher interface
			//which is used to call service package method
			moq := &SentenceFetcherMock{}
			moq.FetchSingleSentenceFunc = func(ctx context.Context, id string) (entity.Sentence, error) {
				if testCase.want.statusCode == http.StatusOK {
					return testCase.sentence, nil
				}
				return testCase.sentence, errors.New("sql: no rows in result set")
			}

			//Send http request
			sut := FetchSingleSentence{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CheckOutHTTPResponse(t,
				resp, testCase.want.statusCode, testhelper.LoadJsonGoldenFile(t, testCase.want.responseFile),
			)
		})
	}
}
