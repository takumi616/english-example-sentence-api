package service

import (
	"context"
	"log"
	"strconv"
)

type DeleteSentence struct {
	//Interface to access store package's method
	Store SentenceDeleter
}

// Delete a sentence by sentence id
func (d *DeleteSentence) DeleteSentence(ctx context.Context, id string) (int64, error) {
	sentenceID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
		return 0, err
	}

	//Call store package's method, using interface
	deleted, err := d.Store.DeleteSentence(ctx, int64(sentenceID))
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
		return 0, err
	}

	return deleted, nil
}
