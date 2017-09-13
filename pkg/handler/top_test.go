package handler

import (
	"context"
	"errors"
	"testing"

	"net/http/httptest"

	"net/http"

	"encoding/json"

	"net/url"

	"bytes"

	"strconv"

	"github.com/smoya/ghtop/pkg/contributor"
	"github.com/smoya/ghtop/pkg/logx"
	"github.com/stretchr/testify/assert"
)

func TestGetTop(t *testing.T) {
	locationParam := "barcelona"
	limitParam := 2
	sortParam := contributor.SortByRepositories

	contributorRepo := contributorRepoMock(locationParam, limitParam, sortParam)
	query := contributor.NewGetTopContributorsQuery(contributorRepo)
	w := httptest.NewRecorder()
	h := GetTop(query, logx.NewDiscardAll())

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Form = url.Values{}
	r.Form.Set("location", locationParam)
	r.Form.Set("limit", strconv.Itoa(limitParam))
	r.Form.Set("sort", sortParam)

	h.ServeHTTP(w, r)

	expectedOutput, err := json.Marshal(contributor.Fixtures()[:limitParam])
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, bytes.TrimRight(w.Body.Bytes(), "\n"))
}

func TestGetTop_InvalidLimit(t *testing.T) {
	locationParam := "barcelona"

	contributorRepo := contributorRepoMock(locationParam, contributor.LimitDefault, contributor.SortDefault)
	query := contributor.NewGetTopContributorsQuery(contributorRepo)

	w := httptest.NewRecorder()
	h := GetTop(query, logx.NewDiscardAll())

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Form = url.Values{}

	r.Form.Set("location", locationParam)
	r.Form.Set("limit", "Invalid Limit")

	h.ServeHTTP(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func contributorRepoMock(locationExpected string, limitExpected int, sortExpected string) *contributor.RepositoryMock {
	contributorRepo := contributor.RepositoryMock{}
	contributorRepo.FindByLocationFunc = func(ctx context.Context, location string, limit int, sortBy string) ([]contributor.Contributor, error) {
		if location != locationExpected {
			return nil, errors.New("invalid location param")
		}

		if limit != limitExpected {
			return nil, errors.New("invalid limit param")
		}

		if sortBy != sortExpected {
			return nil, errors.New("invalid sortBy param")
		}

		if len(contributor.Fixtures()) <= limitExpected {
			return contributor.Fixtures(), nil
		}

		return contributor.Fixtures()[:limit], nil
	}

	return &contributorRepo
}
