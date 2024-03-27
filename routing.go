package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/takumi616/english-example-sentence-api/config"
)

func setUpRouting(ctx context.Context, cfg *config.Config) http.Handler {
	router := chi.NewRouter()

	//An endpoint to check if http server is running correctly.
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	return router
}
