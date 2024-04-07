package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/takumi616/english-example-sentence-api/config"
	"github.com/takumi616/english-example-sentence-api/handler"
	"github.com/takumi616/english-example-sentence-api/service"
	"github.com/takumi616/english-example-sentence-api/store"
)

func setUpRouting(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	router := chi.NewRouter()

	//An endpoint to check if http server is running correctly.
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	//Create validator.
	v := validator.New()

	//Get DB handle.
	//cleanup is used to close *sql.DB.
	dbHandle, cleanup, err := store.ConnectToDatabase(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	//Register http handler.
	repository := &store.Repository{DbHandle: dbHandle}
	c := &handler.CreateSentence{
		Service: &service.CreateSentence{
			Store: repository,
		},
		Validator: v,
	}

	fl := &handler.FetchSentenceList{
		Service: &service.FetchSentence{
			Store: repository,
		},
	}

	fs := &handler.FetchSingleSentence{
		Service: &service.FetchSentence{
			Store: repository,
		},
	}

	d := &handler.DeleteSentence{
		Service: &service.DeleteSentence{
			Store: repository,
		},
	}

	u := &handler.UpdateSentence{
		Service: &service.UpdateSentence{
			Store: repository,
		},
	}

	router.Route("/sentences", func(router chi.Router) {
		router.Post("/", c.ServeHTTP)
		router.Get("/", fl.ServeHTTP)
		router.Route("/{id}", func(router chi.Router) {
			router.Get("/", fs.ServeHTTP)
			router.Delete("/", d.ServeHTTP)
			router.Put("/", u.ServeHTTP)
		})
	})

	return router, cleanup, nil
}
