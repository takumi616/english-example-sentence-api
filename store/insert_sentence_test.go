package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

func TestInsertNewSentence(t *testing.T) {
	//Expected data
	type expected struct {
		inserted int64
	}

	type testdata struct {
		sentence       entity.Sentence
		expected       expected
		lastInsertedId int64
		rowsAffected   int64
	}

	//Prepare test cases
	testcases := map[string]testdata{}

	//Ok test case
	testcases["OK"] = testdata{
		sentence: entity.Sentence{
			Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			Vocabularies: pq.StringArray{"build", "deployment", "container"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expected:       expected{inserted: 6},
		lastInsertedId: 6,
		rowsAffected:   1,
	}

	//Fail test case
	testcases["Fail"] = testdata{
		sentence: entity.Sentence{
			Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			Vocabularies: pq.StringArray{"build", "deployment", "container"},
			Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
			Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
		},
		expected:       expected{inserted: 0},
		lastInsertedId: 0,
		rowsAffected:   0,
	}

	for testcase, testdata := range testcases {
		//Run runs function as a subtest of t called name n(first parameter of Run)
		t.Run(testcase, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()
			//Create a new mock DB
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			//Set expected query to mock
			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO sentence \(body, vocabularies, created, updated\) VALUES\(\$1, \$2, \$3, \$4\)`).
				WithArgs(testdata.sentence.Body, testdata.sentence.Vocabularies, testdata.sentence.Created, testdata.sentence.Updated).
				WillReturnResult(sqlmock.NewResult(testdata.lastInsertedId, testdata.rowsAffected))
			mock.ExpectCommit()

			//Call test target function using mock db
			repository := &Repository{DbHandle: db}
			actualInserted, err := repository.InsertNewSentence(context.Background(), &testdata.sentence)
			if err != nil {
				t.Errorf("Failed to insert new sentence: %v", err)
			}

			//Compare actual to expected
			if actualInserted == testdata.expected.inserted {
				t.Logf("Successfully inserted new sentence and its ID is: %d", actualInserted)
			} else {
				t.Errorf("Unexpected sentence ID. Expected: %d  Actual: %d", testdata.expected.inserted, actualInserted)
			}
		})
	}
}
