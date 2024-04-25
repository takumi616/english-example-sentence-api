package store

import (
	"context"
	"log"

	"github.com/takumi616/generate-example/entity"
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

func (r *Repository) SelectSentenceById(ctx context.Context, sentenceID int64) (entity.Sentence, error) {
	//Execute select query and scan a selected row into go struct
	query := "SELECT * FROM sentence WHERE id = $1"
	var selected entity.Sentence
	err := r.DbHandle.QueryRowContext(ctx, query, sentenceID).Scan(&selected.SentenceID, &selected.Body, &selected.Vocabularies, &selected.Created, &selected.Updated)
	if err != nil {
		log.Printf("Failed to select and scan a row into go struct: %v", err)
		return selected, err
	}

	return selected, nil
}
