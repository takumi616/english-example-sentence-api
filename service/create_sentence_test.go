package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

func TestCreateNewSentence(t *testing.T) {
	//Expected returned value
	type expected struct {
		rowsAffectedNumber int64
		errorMessage       error
	}

	type testdata struct {
		body         string
		vocabularies pq.StringArray
		expected     expected
	}

	testCases := map[string]testdata{}

	//Case Ok
	testCases["Ok"] = testdata{
		body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
		vocabularies: pq.StringArray{"build", "deployment", "container"},
		expected: expected{
			rowsAffectedNumber: 1,
			errorMessage:       nil,
		},
	}

	//Case Status internal server error
	testCases["Internal server error"] = testdata{
		body:         "In the realm of technology, the network serves as the backbone of communication, facilitating seamless data transfer between devices. At the heart of this exchange lies the database, a repository of structured information meticulously organized for efficient retrieval and manipulation. As data flows through the network, it finds refuge in the memory, transient yet vital, where it is temporarily stored and processed to fuel the operations of countless applications, orchestrating the dance of bytes across the digital landscape.",
		vocabularies: pq.StringArray{"network", "database", "memory"},
		expected: expected{
			rowsAffectedNumber: 0,
			errorMessage:       errors.New("pq: value too long for type character varying(120)"),
		},
	}

	for name, testdata := range testCases {
		name := name
		testdata := testdata

		//Run in parallel
		t.Run(name, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//Create mock and set mock function
			ctx := context.Background()

			moq := &SentenceInserterMock{}
			moq.InsertNewSentenceFunc = func(ctx context.Context, sentence *entity.Sentence) (int64, error) {
				if name == "Ok" {
					return 1, nil
				} else if name == "Internal server error" {
					return 0, errors.New("pq: value too long for type character varying(120)")
				} else {
					return 0, errors.New("Not found matched test case")
				}
			}

			//Call test target method using mock interface
			c := &CreateSentence{Store: moq}
			rowsAffectedNumber, err := c.CreateNewSentence(ctx, testdata.vocabularies, testdata.body)
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
			} else {
				t.Error("Not found matched test case")
			}
		})
	}
}
