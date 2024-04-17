package store

import (
	"context"
	"log"

	"github.com/takumi616/generate-example/entity"
)

func (r *Repository) UpdateSentence(ctx context.Context, sentenceID int, body string) (entity.Sentence, error) {
	//Begin a transaction
	tx, err := r.DbHandle.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start a transaction: %v", err)
		return entity.Sentence{}, err
	}

	//Execute update query and fetch an updated record
	var updated entity.Sentence
	query := "UPDATE sentence SET body = $2 WHERE id = $1 RETURNING *"
	err = tx.QueryRowContext(ctx, query, sentenceID, body).Scan(&updated.SentenceID, &updated.Body, &updated.Vocabularies, &updated.Created, &updated.Updated)
	if err != nil {
		//Execute roll back
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("Failed to rollback transaction: %v", rollbackErr)
		}
		log.Printf("Rolled back transaction: %v", err)
		return updated, err
	}

	//Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return updated, err
	}

	return updated, nil
}
