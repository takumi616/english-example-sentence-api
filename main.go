package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/takumi616/english-example-sentence-api/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	dataSourceName := "host=" + config.DBHost + " port=" + config.DBPort + " user=" + config.DBUser + " password=" + config.DBPassword + " dbname=" + config.DBName + " sslmode=" + config.DBSSLMODE
	_, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to postgresql successfully.")
}
