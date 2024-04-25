package handler

import (
	"context"

	"github.com/takumi616/generate-example/entity"
)

//go:generate go run github.com/matryer/moq -out moq_sentence_fetcher.go . SentenceFetcher
//go:generate go run github.com/matryer/moq -out moq_sentence_deleter.go . SentenceDeleter
//go:generate go run github.com/matryer/moq -out moq_sentence_updater.go . SentenceUpdater
//go:generate go run github.com/matryer/moq -out moq_sentence_creater.go . SentenceCreater

// Define interfaces to decouple each package from one another.
// These interfaces are used to access service package's methods.

// Create a sentence
type SentenceCreater interface {
	CreateNewSentence(ctx context.Context, vocabularies []string, body string) (int64, error)
}

// Fetch all sentences or a specific sentence by sentence id
type SentenceFetcher interface {
	FetchSentenceList(ctx context.Context) ([]entity.Sentence, error)
	FetchSingleSentence(ctx context.Context, id string) (entity.Sentence, error)
}

// Delete a sentence by sentence id
type SentenceDeleter interface {
	DeleteSentence(ctx context.Context, id string) (int64, error)
}

// Update a sentence by sentence id
type SentenceUpdater interface {
	UpdateSentence(ctx context.Context, id string, body string) (int64, error)
}
