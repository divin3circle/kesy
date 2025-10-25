package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

type ErrRecordNotFound struct{
	Message string
}

func (e ErrRecordNotFound) Error() string {
	return e.Message
}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	
	return err
}

func ReadParam(r *http.Request, key string) (string, error) {
	idParam := chi.URLParam(r, key)
	if idParam == "" {
		return "", ErrRecordNotFound{Message: "id is required"}
	}
	return idParam, nil
}