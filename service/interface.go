package service

import (
	"context"

	"github.com/takumi616/english-example-sentence-api/entity"
)

type SentenceInserter interface {
	InsertNewSentence(ctx context.Context, sentence *entity.Sentence) (int, error)
}

type SentenceSelecter interface {
	SelectSentenceList(ctx context.Context) ([]entity.Sentence, error)
}
