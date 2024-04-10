package service

import (
	"context"
	"log"

	"github.com/takumi616/generate-example/entity"
)

type CreateSentence struct {
	//Interface to access store package's method
	Store SentenceInserter
}

// Create a sentence
func (c *CreateSentence) CreateNewSentence(ctx context.Context, vocabularies []string, body string) (int, error) {
	sentence := &entity.Sentence{
		Body:         body,
		Vocabularies: vocabularies,
	}

	//Call store package's method, using interface
	sentenceID, err := c.Store.InsertNewSentence(ctx, sentence)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
	}

	return sentenceID, nil
}
