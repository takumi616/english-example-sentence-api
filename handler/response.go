package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	//Content of error message
	Message string `json:"message"`
}

// Write http response to http response writer
func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	//Set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//Get json encoding of body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed to encode response correctly: %v", err)
		//Set error status code
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			Message: err.Error(),
		}
		//Write error response into response writer.
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("Failed to write error response correctly: %v", err)
		}
		return
	}

	//Write status code and response body into response writer.
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("Failed to write response correctly: %v", err)
	}
}
