package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/smoya/ghtop/pkg/contributor"
	"github.com/smoya/ghtop/pkg/handler"
	"github.com/smoya/ghtop/pkg/httpx/middleware"
	"github.com/smoya/ghtop/pkg/logx"
)

// Config ...
type Config struct {
	Port              int
	Env               string
	BasicAuthUser     string
	BasicAuthPassword string
}

// NewConfig creates new config.
func NewConfig(
	port int,
	env string,
	basicAuthUser string,
	basicAuthPassword string,
) Config {
	return Config{
		Port:              port,
		Env:               env,
		BasicAuthUser:     basicAuthUser,
		BasicAuthPassword: basicAuthPassword,
	}
}

// Server serves http requests.
type Server struct {
	Config
	Logger logx.Logger
	Router chi.Router
}

// NewServer creates new Server.
func NewServer(config Config, logger logx.Logger, contributorRepo contributor.Repository) *Server {
	r := chi.NewRouter()

	// Using gzip almost always.
	r.Use(chimiddleware.DefaultCompress)

	if config.Env == "dev" {
		r.Use(chimiddleware.Logger)
	}

	// Routes
	registerRoutes(r, config, logger, contributorRepo)

	return &Server{
		Config: config,
		Logger: logger,
		Router: r,
	}
}

func registerRoutes(r chi.Router, config Config, logger logx.Logger, contributorRepo contributor.Repository) {
	query := contributor.NewGetTopContributorsQuery(contributorRepo)

	if config.BasicAuthUser != "" {
		r = r.With(httpx.BasicAuthentication(config.BasicAuthUser, config.BasicAuthPassword))
	}

	r.Get("/top", handler.GetTop(query, logger))

}

// Run runs the server.
func (s *Server) Run(ctx context.Context) error {
	s.Logger.Info("Server started.", logx.NewField("port", s.Port))

	http.Handle("/", s.Router)

	return http.ListenAndServe(fmt.Sprintf(":%v", s.Port), s.Router)
}

// ServeHTTP serve just one request.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(w, req)
}
