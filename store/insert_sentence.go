package store

import (
	"context"
	"log"
	"time"

	"github.com/takumi616/generate-example/entity"
)

func (r *Repository) InsertNewSentence(ctx context.Context, sentence *entity.Sentence) (int, error) {
	//Begin a transaction
	tx, err := r.DbHandle.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start a transaction: %v", err)
		return 0, err
	}

	//Execute insert query and fetch a new inserted record's ID
	query := "INSERT INTO sentence (body, vocabularies, created, updated) VALUES($1, $2, $3, $4) RETURNING id"
	var inserted entity.Sentence
	err = tx.QueryRowContext(ctx, query, sentence.Body, sentence.Vocabularies, time.Now().String(), time.Now().String()).Scan(&inserted.SentenceID)
	if err != nil {
		//Roll back this transaction
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("Failed to rollback transaction: %v", rollbackErr)
		}
		log.Printf("Rolled back transaction: %v", err)
		return inserted.SentenceID, err
	}

	//Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	return inserted.SentenceID, nil
}
