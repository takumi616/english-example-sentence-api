package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/takumi616/generate-example/entity"
)

// Register test data to mock DB
func prepareTestRecord(testFuncName string) *sqlmock.Rows {
	//Create rows
	rows := sqlmock.NewRows([]string{"id", "body", "vocabularies", "created", "updated"})

	//Prepare test data
	testSentenceList := []entity.Sentence{
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

	//Change the number of records, depending on test target function
	//Test target function:  TestSelectSentenceList
	if testFuncName == "TestSelectSentenceList" {
		rows.AddRow(testSentenceList[0].SentenceID, testSentenceList[0].Body, testSentenceList[0].Vocabularies, testSentenceList[0].Created, testSentenceList[0].Updated)
		rows.AddRow(testSentenceList[1].SentenceID, testSentenceList[1].Body, testSentenceList[1].Vocabularies, testSentenceList[1].Created, testSentenceList[1].Updated)
		//Test target function:  TestSelectSentenceById
	} else {
		rows.AddRow(testSentenceList[0].SentenceID, testSentenceList[0].Body, testSentenceList[0].Vocabularies, testSentenceList[0].Created, testSentenceList[0].Updated)
	}

	return rows
}

func TestSelectSentenceList(t *testing.T) {
	//Create a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	//Prepare mock DB records
	rows := prepareTestRecord("TestSelectSentenceList")

	//Set query to mock
	mock.ExpectQuery(`SELECT \* FROM sentence`).WillReturnRows(rows)

	//Create Repository
	repository := &Repository{DbHandle: db}

	//Call test target method
	selected, err := repository.SelectSentenceList(context.Background())
	if err != nil {
		t.Errorf("Failed to fetch all sentences: %v", err)
	}

	//Compare selected records number to expected
	if len(selected) == 2 {
		t.Log("Fetched all sentences successfully.")
	} else {
		t.Logf("Unexpected number of records: %v", err)
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations were met: %v", err)
	}

	t.Logf("Fetched sentences: %v", selected)
}

func TestSelectSentenceById(t *testing.T) {
	//Create a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	//Register test record to mock DB
	row := prepareTestRecord("TestSelectSentenceByid")

	//Set query to mock
	mock.ExpectQuery(`SELECT \* FROM sentence WHERE id = \$1`).WillReturnRows(row)

	//Create Repository
	repository := &Repository{DbHandle: db}

	//Call test target method
	sentenceID := 5
	selected, err := repository.SelectSentenceById(context.Background(), sentenceID)
	if err != nil {
		t.Errorf("Failed to fetch a sentence: %v", err)
	}

	//Compare selected record to expected
	if selected.SentenceID == sentenceID {
		t.Logf("Successfully fetched a sentence by sentence ID: %v", selected)
	} else {
		t.Errorf("Fetched unexpected sentence: %v", selected)
	}
}
