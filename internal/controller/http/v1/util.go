package v1

import (
	"encoding/json"
	"errors"
	"net/http"
)

func writeJSON(w http.ResponseWriter, value any, code int) {
	response, err := json.Marshal(value)
	if err != nil {
		// TODO: log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func bindJSON(r *http.Request, dst any) error {
	if r.Header.Get("Content-Type") != "" {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			return errors.New("Content-Type header is not application/json")
		}
	} else {
		return nil
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// TODO: better error messages for the end user
	return decoder.Decode(&dst)
}
