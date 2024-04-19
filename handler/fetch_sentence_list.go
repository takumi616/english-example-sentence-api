package handler

import (
	"net/http"
)

type FetchSentenceList struct {
	//Interface to access service package's method
	Service SentenceFetcher
}

// Http handler to fetch all sentences
func (fl *FetchSentenceList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//Call service package's method using interface
	sentences, err := fl.Service.FetchSentenceList(ctx)
	if err != nil {
		WriteJsonResponse(ctx, w, &ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	//Write http response to response writer
	WriteJsonResponse(ctx, w, sentences, http.StatusOK)
}
