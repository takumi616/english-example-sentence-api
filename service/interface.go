package service

import (
	"context"

	"github.com/takumi616/english-example-sentence-api/entity"
)

type SentenceInserter interface {
	InsertNewSentence(ctx context.Context, sentence *entity.Sentence) (int, error)
}
