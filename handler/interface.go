package handler

import (
	"context"

	"github.com/takumi616/english-example-sentence-api/entity"
)

type SentenceCreater interface {
	CreateNewSentence(ctx context.Context, vocabularies []string, body string) (int, error)
}

type SentenceFetcher interface {
	FetchSentenceList(ctx context.Context) ([]entity.Sentence, error)
	FetchSingleSentence(ctx context.Context, id string) (entity.Sentence, error)
}
