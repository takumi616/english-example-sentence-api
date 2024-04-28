package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/takumi616/go-postgres-docker-restapi/config"
)

type Repository struct {
	DbHandle *sql.DB
}

// To be able to close *sql.DB before finishing application process,
// this function needs to return a function which executes *sql.DB.Close()
func ConnectToDatabase(ctx context.Context, config *config.Config) (*Repository, func(), error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMODE)
	dbHandle, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Printf("Failed to open postgresql: %v", err)
		return nil, func() {}, err
	}

	closeDB := func() { _ = dbHandle.Close() }

	//Check if the connection to the database is still alive
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := dbHandle.PingContext(ctx); err != nil {
		return nil, closeDB, err
	}

	repository := &Repository{DbHandle: dbHandle}
	return repository, closeDB, nil
}
