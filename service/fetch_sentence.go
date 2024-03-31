package service

import (
	"context"
	"log"
	"strconv"

	"github.com/takumi616/english-example-sentence-api/entity"
)

type FetchSentence struct {
	Store SentenceSelecter
}

func (f *FetchSentence) FetchSentenceList(ctx context.Context) ([]entity.Sentence, error) {
	sentences, err := f.Store.SelectSentenceList(ctx)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
		return nil, err
	}

	return sentences, nil
}

func (f *FetchSentence) FetchSingleSentence(ctx context.Context, id string) (entity.Sentence, error) {
	sentenceID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
	}

	sentence, err := f.Store.SelectSentenceById(ctx, sentenceID)
	if err != nil {
		log.Printf("Error occurred in service package: %v", err)
		return entity.Sentence{}, err
	}

	return sentence, nil
}
