package contributor

import (
	"context"

	"net/http"

	"fmt"

	"time"

	"sync"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

const (
	// SortByRepositories sorts by the amount of repositories the contributor has.
	SortByRepositories = "repositories"

	// SortByFollowers sorts by the amount of followers the contributor has.
	SortByFollowers = "followers"

	// SortByDateJoined sorts by the date the contributor joined.
	SortByDateJoined = "joined"

	// SortDefault is the default sorting.
	SortDefault = SortByRepositories

	// LimitMax Max contributor results.
	LimitMax = 150

	// LimitDefault is the default amount of contributor results.
	LimitDefault = LimitMax
)

// Repository represents
type Repository interface {
	FindByLocation(ctx context.Context, location string, limit int, sortBy string) ([]Contributor, error)
}

// Note about Github Search API rate limit:
// For requests using Basic Authentication, OAuth, or client ID and secret, you can make up to 30 requests per minute.
// For unauthenticated requests, the rate limit allows you to make up to 10 requests per minute.
type githubRepository struct {
	client *github.Client
}

// FindByLocation fetches contributors by location from github.
func (r *githubRepository) FindByLocation(ctx context.Context, location string, limit int, sort string) ([]Contributor, error) {
	result, response, err := r.client.Search.Users(ctx, fmt.Sprintf("location:%s", location), &github.SearchOptions{
		Sort: sort,
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	})

	if err != nil {
		if response != nil && response.StatusCode != http.StatusOK {
			return nil, errors.Wrap(err, "Error during github search.")
		}

		return nil, err
	}

	var cont []Contributor
	for _, u := range result.Users {
		cont = append(cont, FromGithub(u))
	}

	return cont, nil
}

// NewGihubRepository creates a githubRepository.
func NewGihubRepository(client *github.Client) Repository {
	return &githubRepository{
		client: client,
	}
}

type cachedRepository struct {
	wrappedRepo Repository
	cache       *sync.Map
	ttl         time.Duration
}

type cachedResult struct {
	data     []Contributor
	storedAt time.Time
}

// FindByLocation fetches contributors
func (r *cachedRepository) FindByLocation(ctx context.Context, location string, limit int, sortBy string) ([]Contributor, error) {
	key := fmt.Sprintf("%s#%v#%s", location, limit, sortBy)
	result, ok := r.cache.Load(key)
	if ok {
		cachedResult := result.(cachedResult)
		if time.Now().Sub(cachedResult.storedAt) < r.ttl {
			return cachedResult.data, nil
		}
	}

	freshData, err := r.wrappedRepo.FindByLocation(ctx, location, limit, sortBy)
	if err != nil {
		return nil, err
	}

	r.cache.Store(key, cachedResult{
		data:     freshData,
		storedAt: time.Now(),
	})

	return freshData, nil
}

// WithCache wraps a repository with a cached one.
func WithCache(repo Repository, ttl time.Duration) Repository {
	return &cachedRepository{
		wrappedRepo: repo,
		cache:       new(sync.Map),
		ttl:         ttl,
	}
}
