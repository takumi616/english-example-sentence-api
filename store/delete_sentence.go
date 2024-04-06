package store

import (
	"context"
	"log"

	"github.com/takumi616/english-example-sentence-api/entity"
)

func (r *Repository) DeleteSentence(ctx context.Context, sentenceID int) (int, error) {
	//Begin a transaction
	tx, err := r.DbHandle.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start a transaction: %v", err)
		return 0, err
	}

	//Execute delete query and fetch a deleted record's id
	var deleted entity.Sentence
	query := "DELETE FROM sentence " + "WHERE id = $1 RETURNING id"
	err = tx.QueryRowContext(ctx, query, sentenceID).Scan(&deleted.SentenceID)
	if err != nil {
		//Execute roll back
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("Failed to rollback transaction: %v", rollbackErr)
		}
		log.Printf("Rolled back transaction: %v", err)
		return deleted.SentenceID, err
	}

	//Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	return deleted.SentenceID, nil

}
