package httpx

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/smoya/ghtop/pkg/httpx"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuthentication(t *testing.T) {
	user := "foo"
	pass := "bar"

	// Here we test the first response of the server saying "Eh! you must authenticate"
	NonAuthenticatedHandler := &httpx.HandlerMock{}
	NonAuthenticatedHandler.ServeHTTPFunc = func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "This should not happen since you are not authenticated")
	}
	h := BasicAuthentication(user, pass)(NonAuthenticatedHandler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(recorder, request)

	assert.Equal(t, "Basic", recorder.Header().Get("WWW-Authenticate"))

	// Here we test the response after being authenticated.
	var visited bool
	AuthenticatedHandler := &httpx.HandlerMock{}
	AuthenticatedHandler.ServeHTTPFunc = func(w http.ResponseWriter, r *http.Request) {
		visited = true
	}

	h = BasicAuthentication(user, pass)(AuthenticatedHandler)

	request = httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Add("Authorization", "Basic Zm9vOmJhcg==")
	h.ServeHTTP(recorder, request)

	assert.True(t, visited)
}
