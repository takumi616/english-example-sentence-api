package handler

import (
	"context"

	"github.com/takumi616/english-example-sentence-api/entity"
)

// Define interfaces to decouple each package from one another.
// These interfaces are used to access service package's methods.

// Create a sentence
type SentenceCreater interface {
	CreateNewSentence(ctx context.Context, vocabularies []string, body string) (int, error)
}

// Fetch all sentences or a specific sentence by sentence id
type SentenceFetcher interface {
	FetchSentenceList(ctx context.Context) ([]entity.Sentence, error)
	FetchSingleSentence(ctx context.Context, id string) (entity.Sentence, error)
}
