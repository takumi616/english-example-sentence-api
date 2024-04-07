package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UpdateSentence struct {
	//Interface to access service package's method
	Service SentenceUpdater
	//Http request body validator
	Validator *validator.Validate
}

// Http handler to create a new Sentence
func (u *UpdateSentence) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//Get a param from http request URL
	id := chi.URLParam(r, "id")

	var req struct {
		Body string `json:"body" validate:"required"`
	}

	//Convert json http request data into go struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	//Validate http request body
	err := validator.New().Struct(req)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	//Call service package's method using interface
	sentence, err := u.Service.UpdateSentence(ctx, id, req.Body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	//Write a response to http response writer
	RespondJSON(ctx, w, sentence, http.StatusOK)
}
