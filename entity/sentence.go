package entity

import (
	"github.com/lib/pq"
)

type Sentence struct {
	SentenceID   int            `json:"sentence_id"`
	Body         string         `json:"body"`
	Vocabularies pq.StringArray `json:"vocabularies"`
	Created      string         `json:"created"`
	Updated      string         `json:"updated"`
}
