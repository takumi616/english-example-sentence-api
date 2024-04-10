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

func TestFetchSentenceList(t *testing.T) {
	type wantResult struct {
		statusCode   int
		responseFile string
	}

	type testCase struct {
		sentences []entity.Sentence
		want      wantResult
	}

	//Prepare two test case (ok and empty response data pattern)
	testCases := map[string]testCase{}
	//Case ok
	testCases["ok"] = testCase{
		sentences: []entity.Sentence{
			{
				SentenceID:   1,
				Body:         "This is a test sentence.",
				Vocabularies: pq.StringArray{"This", "is", "test"},
				Created:      "2024-03-28T14:15:21.574757Z",
				Updated:      "2024-03-28T14:15:21.574758Z",
			},
			{
				SentenceID:   2,
				Body:         "This is also test sentence.",
				Vocabularies: pq.StringArray{"This", "is", "also"},
				Created:      "2024-03-28T14:15:25.954024Z",
				Updated:      "2024-03-28T14:15:25.954024Z",
			},
		},
		want: wantResult{
			statusCode:   http.StatusOK,
			responseFile: "../testhelper/golden/fetchlist/ok_resp.json.golden",
		},
	}
	//Case empty
	testCases["empty"] = testCase{
		sentences: []entity.Sentence{},
		want: wantResult{
			statusCode:   http.StatusOK,
			responseFile: "../testhelper/golden/fetchlist/empty_resp.json.golden",
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
			r := httptest.NewRequest(http.MethodGet, "/sentences", nil)

			//SentenceFetcherMock　mocks SentenceFetcher interface
			//which is used to call service package method
			moq := &SentenceFetcherMock{}
			moq.FetchSentenceListFunc = func(ctx context.Context) ([]entity.Sentence, error) {
				if testCase.sentences != nil {
					return testCase.sentences, nil
				}
				return nil, errors.New("Found errors in mock.")
			}

			//Send http request
			sut := FetchSentenceList{Service: moq}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testhelper.CheckOutHTTPResponse(t,
				resp, testCase.want.statusCode, testhelper.LoadJsonGoldenFile(t, testCase.want.responseFile),
			)
		})
	}
}
