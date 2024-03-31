package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FetchSingleSentence struct {
	Service SentenceFetcher
}

func (fs *FetchSingleSentence) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	sentence, err := fs.Service.FetchSingleSentence(ctx, id)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	RespondJSON(ctx, w, sentence, http.StatusOK)
}
