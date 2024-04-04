package service

import (
	"context"

	"github.com/takumi616/english-example-sentence-api/entity"
)

//go:generate go run github.com/matryer/moq -out moq_sentence_selecter.go . SentenceSelecter

// Define interfaces to decouple each package from one another.
// These interfaces are used to access service package's methods.

// Insert a new sentence
type SentenceInserter interface {
	InsertNewSentence(ctx context.Context, sentence *entity.Sentence) (int, error)
}

// Select all sentences or single sentence by sentence id
type SentenceSelecter interface {
	SelectSentenceList(ctx context.Context) ([]entity.Sentence, error)
	SelectSentenceById(ctx context.Context, sentenceID int) (entity.Sentence, error)
}
