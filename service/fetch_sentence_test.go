package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lib/pq"
	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

func TestFetchSentenceList(t *testing.T) {
	type testdata struct {
		expectedRowsNumber int
		sentenceList       []entity.Sentence
	}

	//Prepare two test cases
	testCases := map[string]testdata{}

	//Case
	testCases["Ok"] = testdata{
		expectedRowsNumber: 2,
		sentenceList: []entity.Sentence{
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
	}

	//Empty
	testCases["Empty"] = testdata{
		expectedRowsNumber: 0,
		sentenceList:       []entity.Sentence{},
	}

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
				return testdata.sentenceList, nil
			}

			//Call test target method using mock interface
			f := &FetchSentence{Store: moq}
			actual, err := f.FetchSentenceList(ctx)
			if err != nil {
				t.Errorf("Got an unexpected error: %v", err)
			} else {
				if len(actual) == testdata.expectedRowsNumber {
					t.Log("Successfully got sentences.")
				} else {
					t.Error("Got an unexpected number of sentences.")
				}
			}
		})
	}
}

func TestFetchSingleSentence(t *testing.T) {
	type expected struct {
		sentence     entity.Sentence
		errorMessage error
	}

	type testdata struct {
		expected expected
	}

	//Prepare test cases
	testCases := map[string]testdata{}

	//OK
	testCases["Ok"] = testdata{
		expected: expected{
			sentence: entity.Sentence{
				SentenceID:   1,
				Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
				Vocabularies: pq.StringArray{"build", "deployment", "container"},
				Created:      "2024-04-25 05:25:21.177099049 +0000 UTC m=+161.586254951",
				Updated:      "2024-04-25 05:25:21.177125382 +0000 UTC m=+161.586281242",
			},
			errorMessage: nil,
		},
	}

	//Internal server error(No rows)
	testCases["Internal server error"] = testdata{
		expected: expected{
			sentence:     entity.Sentence{},
			errorMessage: errors.New("sql: no rows in result set"),
		},
	}

	for name, testdata := range testCases {
		testdata := testdata
		name := name
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
				sentence := entity.Sentence{}
				if name == "Ok" {
					sentence.SentenceID = 1
					sentence.Body = "After completing the build process, the application is packaged into a container and ready for deployment."
					sentence.Vocabularies = pq.StringArray{"build", "deployment", "container"}
					sentence.Created = "2024-04-25 05:25:21.177099049 +0000 UTC m=+161.586254951"
					sentence.Updated = "2024-04-25 05:25:21.177125382 +0000 UTC m=+161.586281242"
					return sentence, nil
				} else if name == "Internal server error" {
					return sentence, errors.New("sql: no rows in result set")
				} else {
					return sentence, errors.New("Not found matched test case")
				}
			}

			//Call test target method using mock interface
			f := &FetchSentence{Store: moq}
			actual, err := f.FetchSingleSentence(ctx, "1")
			t.Logf("actual body: %s", actual.Body)

			if name == "Ok" {
				if diff := cmp.Diff(testdata.expected.sentence, actual); diff != "" {
					t.Errorf("Some differences were found: (-expected +actual)\n%s", diff)
				} else {
					t.Log("Got an expected Sentence entity.")
				}
			} else if name == "Internal server error" {
				if testdata.expected.errorMessage.Error() == err.Error() {
					t.Log("Got an expected error message.")
				} else {
					t.Errorf("Failed to get an expected error message: %v", err)
				}
			} else {
				t.Error("Not found matched test case.")
			}
		})
	}
}
