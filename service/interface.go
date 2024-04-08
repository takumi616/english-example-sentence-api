package service

import (
	"context"

	"github.com/takumi616/english-example-sentence-api/entity"
)

//go:generate go run github.com/matryer/moq -out moq_sentence_selecter.go . SentenceSelecter
//go:generate go run github.com/matryer/moq -out moq_sentence_deleter.go . SentenceDeleter
//go:generate go run github.com/matryer/moq -out moq_sentence_updater.go . SentenceUpdater

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

// Delete a sentence by sentence id
type SentenceDeleter interface {
	DeleteSentence(ctx context.Context, sentenceID int) (int, error)
}

// Update a sentence by sentence id
type SentenceUpdater interface {
	UpdateSentence(ctx context.Context, sentenceID int, body string) (entity.Sentence, error)
}
