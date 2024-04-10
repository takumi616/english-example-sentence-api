package service

import (
	"context"
	"log"
	"strconv"

	"github.com/takumi616/generate-example/entity"
)

type UpdateSentence struct {
	//Interface to access store package's method
	Store SentenceUpdater
}

// Update a sentence by sentence id
func (u *UpdateSentence) UpdateSentence(ctx context.Context, id string, body string) (entity.Sentence, error) {
	sentenceID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Failed to get int type sentenceID: %v", err)
		return entity.Sentence{}, err
	}

	//Call store package's method, using interface
	sentence, err := u.Store.UpdateSentence(ctx, sentenceID, body)
	if err != nil {
		log.Printf("Failed to get updated sentence: %v", err)
		return entity.Sentence{}, err
	}

	return sentence, nil
}
