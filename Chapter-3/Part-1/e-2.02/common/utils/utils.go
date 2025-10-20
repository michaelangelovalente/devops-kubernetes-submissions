package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func Write(w http.ResponseWriter, status int, data string) error {
	w.Header().Set("Content-Type", "plain/text")
	_, err := w.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func ReadParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "n")
	if idParam == "" {
		return 0, errors.New("invalid parameter")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid parameter type")
	}

	return id, nil
}
