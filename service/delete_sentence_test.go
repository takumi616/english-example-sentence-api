package service

import (
	"context"
	"errors"
	"testing"
)

func TestDeleteSentence(t *testing.T) {
	//Expected returned value
	type expected struct {
		rowsAffectedNumber int64
		errorMessage       error
	}

	type testdata struct {
		expected expected
	}

	//Prepare test cases
	testCases := map[string]testdata{}

	//OK
	testCases["Ok"] = testdata{
		expected: expected{
			rowsAffectedNumber: 1,
			errorMessage:       nil,
		},
	}

	//Internal server error(No rows)
	testCases["Internal server error"] = testdata{
		expected: expected{
			rowsAffectedNumber: 0,
			errorMessage:       errors.New("sql: no rows in result set"),
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

			//SentenceDeleterMock　mocks SentenceDeleter interface
			//which is used to call service package method
			ctx := context.Background()
			moq := &SentenceDeleterMock{}
			moq.DeleteSentenceFunc = func(ctx context.Context, sentenceID int64) (int64, error) {
				if name == "Ok" {
					return 1, nil
				} else if name == "Internal server error" {
					return 0, errors.New("sql: no rows in result set")
				} else {
					return 0, errors.New("Not found matched test case.")
				}
			}

			//Call test target method using mock interface
			d := &DeleteSentence{Store: moq}
			rowsAffectedNumber, err := d.DeleteSentence(ctx, "1")

			if name == "Ok" {
				if rowsAffectedNumber == testdata.expected.rowsAffectedNumber {
					t.Log("Got an expected rows affected number")
				} else {
					t.Errorf("Failed to get an expected rows affected number: %v", err)
				}
			} else if name == "Internal server error" {
				if err.Error() == testdata.expected.errorMessage.Error() {
					t.Log("Got an expected error message")
				} else {
					t.Errorf("Failed to get an expected error message: %v", err)
				}
			}
		})
	}
}
