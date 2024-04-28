package store

import (
	"context"
	"log"

	"github.com/takumi616/go-postgres-docker-restapi/entity"
)

// Insert a new sentence into db
func (r *Repository) InsertNewSentence(ctx context.Context, sentence *entity.Sentence) (int64, error) {
	//Begin a transaction
	tx, err := r.DbHandle.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start a transaction: %v", err)
		return 0, err
	}

	//Execute insert query and fetch a new inserted record's ID
	query := "INSERT INTO sentence (body, vocabularies, created, updated) VALUES($1, $2, $3, $4)"
	result, err := tx.ExecContext(ctx, query, sentence.Body, sentence.Vocabularies, sentence.Created, sentence.Updated)
	if err != nil {
		//Roll back transaction
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Failed to rollback transaction: %v", rollbackErr)
			return 0, err
		}
		log.Printf("Rolled back transaction: %v", err)
		return 0, err
	}

	//Get rows affected number
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get a rows affected number: %v", err)
		return 0, err
	}
	if rows != 1 {
		log.Printf("Got an unexpected rows affected number: %d", rows)
		return 0, err
	}

	//Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
	}

	return rows, nil
}
