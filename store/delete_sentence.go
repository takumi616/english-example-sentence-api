package store

import (
	"context"
	"log"

	"github.com/takumi616/generate-example/entity"
)

func (r *Repository) DeleteSentence(ctx context.Context, sentenceID int64) (int64, error) {
	//Begin a transaction
	tx, err := r.DbHandle.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start a transaction: %v", err)
		return 0, err
	}

	//Execute delete query and fetch a deleted record's id
	var deleted entity.Sentence
	query := "DELETE FROM sentence WHERE id = $1"
	result, err := tx.ExecContext(ctx, query, sentenceID)
	if err != nil {
		//Execute roll back
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("Failed to rollback transaction: %v", rollbackErr)
		}
		log.Printf("Rolled back transaction: %v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected number: %v", err)
		return 0, err
	}
	if rowsAffected != 1 {
		log.Printf("Got an unexpected rows affected number: %d", rowsAffected)
		return 0, err
	}

	//Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	return deleted.SentenceID, nil
}
