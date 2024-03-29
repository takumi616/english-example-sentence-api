package service

import (
	"context"
	"log"

	"github.com/takumi616/english-example-sentence-api/entity"
)

type FetchSentenceList struct {
	Store SentenceSelecter
}

func (fl *FetchSentenceList) FetchSentenceList(ctx context.Context) ([]entity.Sentence, error) {
	sentences, err := fl.Store.SelectSentenceList(ctx)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
		return nil, err
	}

	return sentences, nil
}
