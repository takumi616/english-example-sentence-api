package store

import (
	"context"
	"log"
)

func (r *Repository) UpdateSentence(ctx context.Context, sentenceID int64, body string) (int64, error) {
	//Begin a transaction
	tx, err := r.DbHandle.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start a transaction: %v", err)
		return 0, err
	}

	//Execute update query and fetch an updated record
	query := "UPDATE sentence SET body = $2 WHERE id = $1"
	result, err := tx.ExecContext(ctx, query, sentenceID, body)
	if err != nil {
		//Execute roll back
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("Failed to rollback transaction: %v", rollbackErr)
		}
		log.Printf("Rolled back transaction: %v", err)
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected number: %v", err)
		return 0, err
	}
	if rows != 1 {
		log.Printf("Got an unexpected rows affected number: %d", rows)
		return 0, err
	}

	//Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return 0, err
	}

	return rows, nil
}
