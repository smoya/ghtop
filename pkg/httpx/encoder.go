package httpx

import (
	"encoding/json"
	"net/http"
)

// M represents a map to encode json.
type M map[string]interface{}

// Encoder encodes and writes an http response. It can modify headers.
type Encoder func(w http.ResponseWriter, data interface{}) error

// JSONEncoder encodes data to json and writes the response + the proper content-type header.
func JSONEncoder(w http.ResponseWriter, data interface{}) error {
	// This should happen before call the Encode method. If not, this header will not be written.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}
