package service

import (
	"context"
	"log"
	"strconv"
)

type UpdateSentence struct {
	//Interface to access store package's method
	Store SentenceUpdater
}

// Update a sentence by sentence id
func (u *UpdateSentence) UpdateSentence(ctx context.Context, id string, body string) (int64, error) {
	sentenceID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Failed to get int type sentenceID: %v", err)
		return 0, err
	}

	//Call store package's method, using interface
	sentence, err := u.Store.UpdateSentence(ctx, int64(sentenceID), body)
	if err != nil {
		log.Printf("Failed to get updated sentence: %v", err)
		return 0, err
	}

	return sentence, nil
}
