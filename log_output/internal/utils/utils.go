package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Envelope map[string]any

func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(js); err != nil {
		log.Printf("Failed to write response: %v", err)
		return err
	}

	return nil
}
