package v1

import (
	"encoding/json"
	"errors"
	"net/http"
)

func writeJSON(w http.ResponseWriter, value any) {
	response, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func bindJSON(r *http.Request, dst any) error {
	if r.Header.Get("Content-Type") != "" {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			return errors.New("Content-Type header is not application/json")
		}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(&dst)
}
