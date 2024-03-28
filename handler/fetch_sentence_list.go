package handler

import (
	"net/http"
)

type FetchSentenceList struct {
	Service SentenceListFetcher
}

// Http handler to fetch all sentences
func (fl *FetchSentenceList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//Call service layer method using interface
	sentences, err := fl.Service.FetchSentenceList(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	RespondJSON(ctx, w, sentences, http.StatusOK)
}
