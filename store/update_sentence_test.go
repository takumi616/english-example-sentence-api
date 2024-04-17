package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/takumi616/generate-example/entity"
)

func TestUpdateSentence(t *testing.T) {
	//Create a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	var updateParam struct {
		//Sentence ID which is used to specify the sentence
		sentenceID int
		//New sentence body
		body string
	}
	updateParam.sentenceID = 6
	updateParam.body = "After finishing the build process successfully, this application is packaged into a container and ready for deployment."

	//Existed sentence data
	testSentence := entity.Sentence{
		SentenceID:   6,
		Body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
		Vocabularies: pq.StringArray{"build", "deployment", "container"},
		Created:      "2024-04-06 20:16:35.47968413 +0000 UTC m=+25.323730179",
		Updated:      "2024-04-06 20:16:35.47969263 +0000 UTC m=+25.323738679",
	}

	//Prepare the sentence which is returned from mock DB
	rows := sqlmock.NewRows([]string{"id", "body", "vocabularies", "created", "updated"})
	rows.AddRow(testSentence.SentenceID, updateParam.body, testSentence.Vocabularies, testSentence.Created, testSentence.Updated)

	//Set query to mock
	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE sentence SET body \= \$2 WHERE id \= \$1 RETURNING \*`).WithArgs(updateParam.sentenceID, updateParam.body).WillReturnRows(rows)
	mock.ExpectCommit()

	//Create Repository
	repository := &Repository{DbHandle: db}

	//Call test target method
	updated, err := repository.UpdateSentence(context.Background(), updateParam.sentenceID, updateParam.body)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	//Compare updated result to expected
	if updated.Body == updateParam.body {
		t.Logf("Updated sentence successfully: %v", updated)
	} else {
		t.Errorf("Failed to update sentence. body is: %s", updated.Body)
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations were met: %v", err)
	}
}
