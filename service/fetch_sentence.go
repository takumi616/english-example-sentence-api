package service

import (
	"context"
	"log"
	"strconv"

	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

type FetchSentence struct {
	//Interface to access store package's method
	Store SentenceSelecter
}

// Fetch all sentences
func (f *FetchSentence) FetchSentenceList(ctx context.Context) ([]entity.Sentence, error) {
	//Call store package's method, using interface
	sentences, err := f.Store.SelectSentenceList(ctx)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
		return nil, err
	}

	return sentences, nil
}

// Fetch single sentence by sentence id
func (f *FetchSentence) FetchSingleSentence(ctx context.Context, id string) (entity.Sentence, error) {
	// Convert string into int
	sentenceID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Failed to convert string type into int: %v", err)
	}

	// Call store package's method, using interface
	fetched, err := f.Store.SelectSentenceById(ctx, int64(sentenceID))
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
		return entity.Sentence{}, err
	}

	return fetched, nil
}
