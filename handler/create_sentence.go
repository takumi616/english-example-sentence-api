package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CreateSentence struct {
	Service   SentenceCreater
	Validator *validator.Validate
}

// Http handler to add a new task
func (c *CreateSentence) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		Vocabularies []string `json:"vocabularies" validate:"required"`
		Body         string   `json:"body" validate:"required"`
	}

	//Convert json http request data into go struct type
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//Validate request body
	err := validator.New().Struct(req)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	//Call service layer method using interface
	sentenceID, err := c.Service.CreateNewSentence(ctx, req.Vocabularies, req.Body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//Create json response body writing newly created record's id
	rsp := struct {
		SentenceID int `json:"sentence_id"`
	}{SentenceID: sentenceID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
