package contributor

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetTopContributorsQuery(t *testing.T) {
	location := "barcelona"
	limit := 2
	sort := SortByRepositories

	contributorRepo := contributorRepoMock(location, limit, sort)
	q := NewGetTopContributorsQuery(contributorRepo)

	con, err := q.Execute(context.Background(), location, limit, sort)
	assert.NoError(t, err)

	assert.Equal(t, Fixtures()[:limit], con)
}

func TestGetTopContributorsQuery_MissingLocation(t *testing.T) {
	q := NewGetTopContributorsQuery(&RepositoryMock{})

	con, err := q.Execute(context.Background(), "", 0, "")
	assert.Equal(t, ErrMissingLocation, err)
	assert.Empty(t, con)
}

func TestGetTopContributorsQuery_InvalidLimit(t *testing.T) {
	q := NewGetTopContributorsQuery(&RepositoryMock{})

	con, err := q.Execute(context.Background(), "barcelona", LimitMax+1, "")
	assert.Equal(t, ErrLimitGreaterThanMax, err)
	assert.Empty(t, con)
}

func TestGetTopContributorsQuery_DefaultValues(t *testing.T) {
	location := "barcelona"

	contributorRepo := contributorRepoMock(location, LimitDefault, SortDefault)
	q := NewGetTopContributorsQuery(contributorRepo)

	con, err := q.Execute(context.Background(), location, 0, "")

	assert.NoError(t, err)
	assert.Equal(t, Fixtures(), con)
}

func contributorRepoMock(locationExpected string, limitExpected int, sortExpected string) *RepositoryMock {
	contributorRepo := RepositoryMock{}
	contributorRepo.FindByLocationFunc = func(ctx context.Context, location string, limit int, sortBy string) ([]Contributor, error) {
		if location != locationExpected {
			return nil, errors.New("Invalid location param")
		}

		if limit != limitExpected {
			return nil, errors.New("Invalid limit param")
		}

		if sortBy != sortExpected {
			return nil, errors.New("Invalid sortBy param")
		}

		if len(Fixtures()) <= limitExpected {
			return Fixtures(), nil
		}

		return Fixtures()[:limit], nil
	}

	return &contributorRepo
}
