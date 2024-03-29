package store

import (
	"context"
	"log"

	"github.com/takumi616/english-example-sentence-api/entity"
)

func (r *Repository) SelectSentenceList(ctx context.Context) ([]entity.Sentence, error) {
	//Execute select query
	query := "SELECT * FROM sentence"
	rows, err := r.DbHandle.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Failed to select records: %v", err)
		return nil, err
	}

	//Scan selected rows into go struct
	var selected []entity.Sentence
	for rows.Next() {
		sentence := entity.Sentence{}
		err = rows.Scan(&sentence.SentenceID, &sentence.Body, &sentence.Vocabularies, &sentence.Created, &sentence.Updated)
		if err != nil {
			log.Printf("Failed to scan a row into go struct: %v", err)
			return nil, err
		}
		selected = append(selected, sentence)
	}

	return selected, nil
}
