package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FetchSingleSentence struct {
	//Interface to access service package's method
	Service SentenceFetcher
}

// Http handler to fetch a sentence by sentence id
func (fs *FetchSingleSentence) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//Get a param from http request URL
	id := chi.URLParam(r, "id")

	//Call service package's method using interface
	sentence, err := fs.Service.FetchSingleSentence(ctx, id)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	//Write http response to response writer
	RespondJSON(ctx, w, sentence, http.StatusOK)
}
