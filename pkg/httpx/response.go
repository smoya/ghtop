package httpx

import (
	"net/http"
)

// WriteJSONOk writes a 200 http response encoding given data.
func WriteJSONOk(w http.ResponseWriter, data interface{}) error {
	err := JSONEncoder(w, data)
	if err != nil {
		return err
	}

	return nil
}
