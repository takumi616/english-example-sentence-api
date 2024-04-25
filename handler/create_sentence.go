package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CreateSentence struct {
	//Interface to access service package's method
	Service SentenceCreater
	//Http request body validator
	Validator *validator.Validate
}

// Http handler to create a new Sentence
func (c *CreateSentence) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		Vocabularies []string `json:"vocabularies" validate:"required"`
		Body         string   `json:"body" validate:"required"`
	}

	//Convert json http request data into go struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJsonResponse(ctx, w, &ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	//Validate http request body
	err := validator.New().Struct(req)
	if err != nil {
		WriteJsonResponse(ctx, w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	//Call service package's method using interface.
	rowsAffected, err := c.Service.CreateNewSentence(ctx, req.Vocabularies, req.Body)
	if err != nil {
		WriteJsonResponse(ctx, w, &ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	//Write response body into response writer
	rsp := struct {
		RowsAffectedNumber int64 `json:"rows_affected_number"`
	}{RowsAffectedNumber: rowsAffected}
	WriteJsonResponse(ctx, w, rsp, http.StatusOK)
}
