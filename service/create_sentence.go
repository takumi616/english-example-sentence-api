package service

import (
	"context"
	"log"
	"time"

	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

type CreateSentence struct {
	//Interface to access store package's method
	Store SentenceInserter
}

// Create a sentence
func (c *CreateSentence) CreateNewSentence(ctx context.Context, vocabularies []string, body string) (int64, error) {
	sentence := &entity.Sentence{
		Body:         body,
		Vocabularies: vocabularies,
		Created:      time.Now().String(),
		Updated:      time.Now().String(),
	}

	//Call store package's method, using interface
	sentenceID, err := c.Store.InsertNewSentence(ctx, sentence)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
	}

	return sentenceID, nil
}
