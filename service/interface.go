package service

import (
	"context"

	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

//go:generate go run github.com/matryer/moq -out moq_sentence_selecter.go . SentenceSelecter
//go:generate go run github.com/matryer/moq -out moq_sentence_deleter.go . SentenceDeleter
//go:generate go run github.com/matryer/moq -out moq_sentence_updater.go . SentenceUpdater
//go:generate go run github.com/matryer/moq -out moq_sentence_inserter.go . SentenceInserter

// Define interfaces to decouple each package from one another.
// These interfaces are used to access service package's methods.

// Insert a new sentence
type SentenceInserter interface {
	InsertNewSentence(ctx context.Context, sentence *entity.Sentence) (int64, error)
}

// Select all sentences or single sentence by sentence id
type SentenceSelecter interface {
	SelectSentenceList(ctx context.Context) ([]entity.Sentence, error)
	SelectSentenceById(ctx context.Context, sentenceID int64) (entity.Sentence, error)
}

// Delete a sentence by sentence id
type SentenceDeleter interface {
	DeleteSentence(ctx context.Context, sentenceID int64) (int64, error)
}

// Update a sentence by sentence id
type SentenceUpdater interface {
	UpdateSentence(ctx context.Context, sentenceID int64, body string) (int64, error)
}
