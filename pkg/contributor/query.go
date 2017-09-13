package contributor

import (
	"context"
	"errors"
	"fmt"
)

// Errors on executing the GetTopContributorsQuery.Execute.
var (
	ErrMissingLocation     = errors.New("missing location")
	ErrLimitGreaterThanMax = fmt.Errorf("limit can not be greater than %v", LimitMax)
	ErrInvalidSort         = fmt.Errorf("sort should be one of: %s|%s|%s", SortByRepositories, SortByFollowers, SortByDateJoined)
)

// GetTopContributorsQuery is the use case that fetches the top contributors by location.
type GetTopContributorsQuery struct {
	repo Repository
}

// Execute executes the use case.
func (q *GetTopContributorsQuery) Execute(ctx context.Context, location string, limit int, sort string) ([]Contributor, error) {
	if location == "" {
		return nil, ErrMissingLocation
	}

	if limit == 0 {
		limit = LimitDefault
	}

	if limit > LimitMax {
		return nil, ErrLimitGreaterThanMax
	}

	if sort == "" {
		sort = SortDefault
	} else if sort != SortByRepositories && sort != SortByFollowers && sort != SortByDateJoined {
		return nil, ErrInvalidSort
	}

	cont, err := q.repo.FindByLocation(ctx, location, limit, sort)
	if err != nil {
		return nil, err
	}

	return cont, nil
}

// NewGetTopContributorsQuery creates a GetTopContributorsQuery.
func NewGetTopContributorsQuery(repo Repository) *GetTopContributorsQuery {
	return &GetTopContributorsQuery{
		repo: repo,
	}
}
