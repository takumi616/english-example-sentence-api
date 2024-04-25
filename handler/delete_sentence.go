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
		WriteJsonResponse(ctx, w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	//Struct of http response body
	rsp := struct {
		SentenceID int64 `json:"sentence_id"`
	}{SentenceID: sentenceID}

	//Write a response to http response writer
	WriteJsonResponse(ctx, w, rsp, http.StatusOK)
}
