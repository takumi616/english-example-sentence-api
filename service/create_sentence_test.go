package service

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

func TestCreateNewSentence(t *testing.T) {
	//Test request body
	type requestBody struct {
		body         string
		vocabularies pq.StringArray
	}

	type testData struct {
		//Expected returned value
		expectedSentenceID int64
		//Test request body
		requestBody requestBody
	}

	//Case Ok
	testCases := map[string]testData{}
	testCases["ok"] = testData{
		expectedSentenceID: 6,
		requestBody: requestBody{
			body:         "After completing the build process, the application is packaged into a container and ready for deployment.",
			vocabularies: pq.StringArray{"build", "deployment", "container"},
		},
	}

	for name, testData := range testCases {
		name := name
		testData := testData

		//Run in parallel
		t.Run(name, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests
			t.Parallel()

			//Create mock and set mock function
			ctx := context.Background()
			moq := &SentenceInserterMock{}
			moq.InsertNewSentenceFunc = func(ctx context.Context, sentence *entity.Sentence) (int64, error) {
				return testData.expectedSentenceID, nil
			}

			//Call test target method using mock interface
			c := &CreateSentence{Store: moq}
			createdSentenceID, err := c.CreateNewSentence(ctx, testData.requestBody.vocabularies, testData.requestBody.body)
			if err != nil {
				t.Errorf("Failed to create new sentence: %v", err)
			}

			//Compare sentenceID
			if createdSentenceID != testData.expectedSentenceID {
				t.Errorf("Unexpected sentenceID. Expected: %d  Result: %d", testData.expectedSentenceID, createdSentenceID)
			} else {
				t.Logf("Inserted sentence ID is: %d", createdSentenceID)
			}
		})
	}
}
