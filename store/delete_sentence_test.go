package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDeleteSentence(t *testing.T) {
	//Create a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	//Prepare returned sentenceID
	sentenceID := 6
	rows := mock.NewRows([]string{"id"})
	rows.AddRow(sentenceID)

	//Set query to mock DB
	mock.ExpectBegin()
	mock.ExpectQuery(`DELETE FROM sentence WHERE id \= \$1 RETURNING id`).WithArgs(sentenceID).WillReturnRows(rows)
	mock.ExpectCommit()

	//Call test target method
	repository := &Repository{DbHandle: db}
	deletedID, err := repository.DeleteSentence(context.Background(), sentenceID)
	if err != nil {
		t.Errorf("Failed to delete a sentence: %v", err)
	} else {
		t.Logf("Successfully deleted a sentence and its sentenceID is: %d", deletedID)
	}
}
