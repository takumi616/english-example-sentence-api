package service

import (
	"context"
	"log"

	"github.com/takumi616/english-example-sentence-api/entity"
)

type CreateSentence struct {
	Store SentenceInserter
}

func (c *CreateSentence) CreateNewSentence(ctx context.Context, vocabularies []string, body string) (int, error) {
	sentence := &entity.Sentence{
		Body:         body,
		Vocabularies: vocabularies,
	}

	sentenceID, err := c.Store.InsertNewSentence(ctx, sentence)
	if err != nil {
		log.Printf("Error occured in service package: %v", err)
	}

	return sentenceID, nil
}
