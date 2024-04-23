package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/takumi616/generate-example/entity"
)

func TestInsertNewSentence(t *testing.T) {
	//Create a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	testSentence := &entity.Sentence{
		Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
		Vocabularies: pq.StringArray{"build", "deployment", "container"},
		Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
		Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
	}

	expected := 6
	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(6)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO sentence \(body, vocabularies, created, updated\) VALUES\(\$1, \$2, \$3, \$4\) RETURNING id`).WithArgs(testSentence.Body, testSentence.Vocabularies, testSentence.Created, testSentence.Updated).WillReturnRows(rows)
	mock.ExpectCommit()

	repository := &Repository{DbHandle: db}
	insertedID, err := repository.InsertNewSentence(context.Background(), testSentence)
	if err != nil {
		t.Errorf("Failed to insert new sentence: %v", err)
	}

	if insertedID == expected {
		t.Logf("Successfully inserted new sentence and its ID is: %d", insertedID)
	} else {
		t.Errorf("Unexpected sentence ID. Expected: %d  Actual: %d", expected, insertedID)
	}
}
