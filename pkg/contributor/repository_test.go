package contributor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"context"

	"time"

	"github.com/google/go-github/github"
	"github.com/smoya/ghtop/pkg/httpx"
	"github.com/stretchr/testify/assert"
)

func TestGithubRepositoryFindByLocation(t *testing.T) {
	server := httptest.NewServer(&githubTestHandler{})
	serverURL, err := url.Parse(server.URL + "/")
	assert.NoError(t, err)

	c := http.DefaultClient
	ghClient := github.NewClient(c)
	ghClient.BaseURL = serverURL

	r := NewGihubRepository(ghClient)
	con, err := r.FindByLocation(context.Background(), "", 0, "")

	assert.NoError(t, err)
	assert.Equal(t, Fixtures(), con)
}

func TestCachedRepositoryFindByLocation(t *testing.T) {
	ttl := 5 * time.Minute // Time enough.

	var visitedTimes int
	wrappedRepo := RepositoryMock{}
	wrappedRepo.FindByLocationFunc = func(ctx context.Context, location string, limit int, sortBy string) ([]Contributor, error) {
		visitedTimes++
		return Fixtures(), nil
	}

	r := WithCache(&wrappedRepo, ttl)

	for i := 0; i < 50; i++ {
		_, err := r.FindByLocation(context.Background(), "barcelona", 0, "")
		assert.NoError(t, err)
		assert.Equal(t, 1, visitedTimes)
	}
}

func TestCachedRepositoryFindByLocation_ExpiredData(t *testing.T) {
	ttl := 1 * time.Millisecond

	var visitedTimes int
	wrappedRepo := RepositoryMock{}
	wrappedRepo.FindByLocationFunc = func(ctx context.Context, location string, limit int, sortBy string) ([]Contributor, error) {
		visitedTimes++
		return Fixtures(), nil
	}

	r := WithCache(&wrappedRepo, ttl)

	_, err := r.FindByLocation(context.Background(), "barcelona", 0, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, visitedTimes)

	time.Sleep(ttl)

	_, err = r.FindByLocation(context.Background(), "barcelona", 0, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, visitedTimes)
}

type githubTestHandler struct{}

func (h *githubTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := githubUsersFromContributors(Fixtures())
	err := httpx.WriteJSONOk(w, result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func githubUsersFromContributors(con []Contributor) github.UsersSearchResult {
	users := make([]github.User, len(con))
	for key, c := range con {
		contrib := c
		users[key] = github.User{
			Login:     &contrib.Username,
			URL:       &contrib.ProfileURL,
			AvatarURL: &contrib.AvatarURL,
		}
	}

	total := len(users)
	incomplete := false

	return github.UsersSearchResult{
		Total:             &total,
		IncompleteResults: &incomplete,
		Users:             users,
	}
}
