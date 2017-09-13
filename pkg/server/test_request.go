package server

import (
	"net/http"
	"net/http/httptest"
)

// Request launches request on given server and return response.
func Request(r *Server, req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	return resp
}
