package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/generate-example/entity"
)

func TestDeleteSentence(t *testing.T) {
	//Prepare two test case (ok and error)
	testCases := map[string][]entity.Sentence{}

	//Case ok
	testCases["ok"] = []entity.Sentence{
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

	//Case error
	testCases["error"] = []entity.Sentence{
		{
			SentenceID:   1,
			Body:         "This is a test sentence.",
			Vocabularies: pq.StringArray{"This", "is", "test"},
			Created:      "2024-03-28T14:15:21.574757Z",
			Updated:      "2024-03-28T14:15:21.574758Z",
		},
	}

	for name, testSentenceList := range testCases {
		name := name
		testSentenceList := testSentenceList
		//Execute as parallel tests
		//Run runs function as a subtest of t called name n(first parameter of Run)
		//It runs function in a separate goroutine and blocks
		//until this function returns or calls t.Parallel to become a parallel test
		t.Run(name, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//SentenceDeleterMock　mocks SentenceDeleter interface
			//which is used to call service package method
			ctx := context.Background()
			moq := &SentenceDeleterMock{}
			moq.DeleteSentenceFunc = func(ctx context.Context, sentenceID int) (int, error) {
				for _, sentence := range testSentenceList {
					if sentence.SentenceID == sentenceID {
						return sentence.SentenceID, nil
					}
				}
				return 0, errors.New("sql: no rows in result set")
			}

			//Call test target method using mock interface
			d := &DeleteSentence{Store: moq}
			got, err := d.DeleteSentence(ctx, "2")
			if err != nil && err.Error() != "sql: no rows in result set" {
				t.Errorf("Failed to get expected result: %v", err)
			}

			if err == nil {
				t.Logf("Got result: %v", got)
			}
		})
	}
}
