package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/english-example-sentence-api/entity"
)

type testSentenceList []entity.Sentence

func TestFetchSentenceList(t *testing.T) {
	//Prepare two test case (ok and empty response data pattern)
	testCases := map[string]testSentenceList{}

	//Case ok
	testCases["ok"] = testSentenceList{
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
	}

	//Case empty
	testCases["empty"] = testSentenceList{}

	for name, testdata := range testCases {
		testdata := testdata
		//Execute as parallel tests
		//Run runs function as a subtest of t called name n(first parameter of Run)
		//It runs function in a separate goroutine and blocks
		//until this function returns or calls t.Parallel to become a parallel test
		t.Run(name, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//SentenceSelecterMock　mocks SentenceSelecter interface
			//which is used to call service package method
			ctx := context.Background()
			moq := &SentenceSelecterMock{}
			moq.SelectSentenceListFunc = func(ctx context.Context) ([]entity.Sentence, error) {
				return testdata, nil
			}

			//Call test target method using mock interface
			f := &FetchSentence{Store: moq}
			got, err := f.FetchSentenceList(ctx)
			if len(got) != 0 && len(got) != 2 {
				t.Errorf("Failed to get expected result: %v", err)
			}
		})
	}
}

func TestFetchSingleSentence(t *testing.T) {
	//Prepare two test case (ok and empty response data pattern)
	testCases := map[string]testSentenceList{}

	//Case ok
	testCases["ok"] = testSentenceList{
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
	}

	//Case empty
	testCases["error"] = testSentenceList{
		{
			SentenceID:   1,
			Body:         "This is a test sentence.",
			Vocabularies: pq.StringArray{"This", "is", "test"},
			Created:      "2024-03-28T14:15:21.574757Z",
			Updated:      "2024-03-28T14:15:21.574758Z",
		},
	}

	for name, testdata := range testCases {
		name := name
		testdata := testdata
		//Execute as parallel tests
		//Run runs function as a subtest of t called name n(first parameter of Run)
		//It runs function in a separate goroutine and blocks
		//until this function returns or calls t.Parallel to become a parallel test
		t.Run(name, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//SentenceSelecterMock　mocks SentenceSelecter interface
			//which is used to call service package method
			ctx := context.Background()
			moq := &SentenceSelecterMock{}
			moq.SelectSentenceByIdFunc = func(ctx context.Context, sentenceID int) (entity.Sentence, error) {
				for _, sentence := range testdata {
					if sentence.SentenceID == sentenceID {
						return sentence, nil
					}
				}
				return entity.Sentence{}, errors.New("sql: no rows in result set")
			}

			//Call test target method using mock interface
			f := &FetchSentence{Store: moq}
			got, err := f.FetchSingleSentence(ctx, "2")
			if err != nil && err.Error() != "sql: no rows in result set" {
				t.Errorf("Failed to get expected result: %v", err)
			}

			if err == nil {
				t.Logf("Got result: %v", got)
			}
		})
	}
}
