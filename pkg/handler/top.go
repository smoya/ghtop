package handler

import (
	"net/http"

	"strconv"

	"github.com/smoya/ghtop/pkg/contributor"
	"github.com/smoya/ghtop/pkg/httpx"
	"github.com/smoya/ghtop/pkg/logx"
)

// GetTop handles GET /top requests.
func GetTop(query *contributor.GetTopContributorsQuery, logger logx.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.Debug("Error parsing the request form", logx.NewField("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loc := r.Form.Get("location")
		if loc == "" {
			logger.Debug("Missing location parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var limit int
		limitStr := r.Form.Get("limit")
		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				logger.Debug("Limit parameter should be a number")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		sort := r.Form.Get("sort")

		cont, err := query.Execute(r.Context(), loc, limit, sort)
		if err != nil {
			logger.Debug("Error fetching contributors", logx.NewField("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = httpx.WriteJSONOk(w, cont)
		if err != nil {
			logger.Debug("Error writing contributors output", logx.NewField("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
