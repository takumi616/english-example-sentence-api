package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

func TestUpdateSentence(t *testing.T) {
	//Create a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	//Existed sentence data
	testSentence := entity.Sentence{
		SentenceID:   6,
		Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
		Vocabularies: pq.StringArray{"build", "deployment", "container"},
		Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
		Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
	}

	//Set query to mock
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE sentence SET body \= \$2 WHERE id \= \$1`).WithArgs(testSentence.SentenceID, testSentence.Body).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//Create Repository
	repository := &Repository{DbHandle: db}

	//Call test target method
	rows, err := repository.UpdateSentence(context.Background(), testSentence.SentenceID, testSentence.Body)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	//Compare updated result to expected
	if rows == 1 {
		t.Logf("Updated sentence successfully: %d", rows)
	} else {
		t.Errorf("Failed to update sentence. Rows affected number is: %d", rows)
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations were met: %v", err)
	}
}
