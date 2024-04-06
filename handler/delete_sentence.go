package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DeleteSentence struct {
	//Interface to access service package's method
	Service SentenceDeleter
}

// Http handler to delete a Sentence
func (d *DeleteSentence) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//Get a param from http request URL
	id := chi.URLParam(r, "id")

	//Call service package's method using interface
	sentenceID, err := d.Service.DeleteSentence(ctx, id)
	if err != nil {
		RespondJSON(ctx, w, ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
	}

	//Create json response body writing deleted record's id.
	rsp := struct {
		SentenceID int `json:"sentence_id"`
	}{SentenceID: sentenceID}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
