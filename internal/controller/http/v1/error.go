package v1

import (
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, errorResponse{Error: err.Error()}, http.StatusBadRequest)
}
