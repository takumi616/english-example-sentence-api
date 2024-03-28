package entity

import (
	"time"

	"github.com/lib/pq"
)

type Sentence struct {
	SentenceID   int            `json:"sentence_id"`
	Body         string         `json:"body"`
	Vocabularies pq.StringArray `json:"vocabularies"`
	Created      time.Time      `json:"created"`
	Updated      time.Time      `json:"updated"`
}
