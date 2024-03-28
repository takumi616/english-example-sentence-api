package handler

import "context"

type SentenceCreater interface {
	CreateNewSentence(ctx context.Context, vocabularies []string, body string) (int, error)
}
