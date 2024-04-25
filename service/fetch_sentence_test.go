package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/generate-example/entity"
)

func TestFetchSentenceList(t *testing.T) {
	//Prepare two test cases
	testCases := map[string][]entity.Sentence{}

	//Case
	testCases["ok"] = []entity.Sentence{
		{
			SentenceID:   5,
			Body:         "The application communicates with the database server to retrieve and store data.",
			Vocabularies: pq.StringArray{"application", "store", "server"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		{
			SentenceID:   6,
			Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			Vocabularies: pq.StringArray{"build", "deployment", "container"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
	}

	//Empty
	testCases["empty"] = []entity.Sentence{}

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

			//SentenceSelecterMock　mocks SentenceSelecter interface
			//which is used to call service package method
			ctx := context.Background()
			moq := &SentenceSelecterMock{}
			moq.SelectSentenceListFunc = func(ctx context.Context) ([]entity.Sentence, error) {
				return testData, nil
			}

			//Call test target method using mock interface
			f := &FetchSentence{Store: moq}
			fetchedSentenceList, err := f.FetchSentenceList(ctx)
			if len(fetchedSentenceList) != 0 && len(fetchedSentenceList) != 2 {
				t.Errorf("Failed to get expected result: %v", err)
			} else {
				t.Logf("Fetched sentence list: %v", fetchedSentenceList)
			}
		})
	}
}

func TestFetchSingleSentence(t *testing.T) {
	//Prepare two test case (ok and empty response data pattern)
	testCases := map[string]entity.Sentence{}

	//OK
	testCases["ok"] = entity.Sentence{
		SentenceID:   6,
		Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
		Vocabularies: pq.StringArray{"build", "deployment", "container"},
		Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
		Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
	}

	//Data does not exist
	testCases["error"] = entity.Sentence{
		SentenceID:   5,
		Body:         "The application communicates with the database server to retrieve and store data.",
		Vocabularies: pq.StringArray{"application", "store", "server"},
		Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
		Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
	}

	for name, _ := range testCases {
		//TODO
		//Temp edit
		//testSentence := testSentence
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
			moq.SelectSentenceByIdFunc = func(ctx context.Context, sentenceID int64) (entity.Sentence, error) {
				// if testSentence.SentenceID == sentenceID {
				// 	return testSentence, nil
				// }
				return entity.Sentence{}, errors.New("sql: no rows in result set")
			}

			//Call test target method using mock interface
			f := &FetchSentence{Store: moq}
			fetchedSentence, err := f.FetchSingleSentence(ctx, "6")
			if err != nil && err.Error() != "sql: no rows in result set" {
				t.Errorf("Failed to get expected result: %v", err)
			} else {
				t.Logf("Fetched sentence: %v", fetchedSentence)
			}
		})
	}
}
