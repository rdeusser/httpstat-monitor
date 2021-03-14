package httpstat

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	srv    *http.Server
	mux    *http.ServeMux
	logger zerolog.Logger
	urls   []string
}

func MonitorURL(url string) func(*Server) {
	return func(srv *Server) {
		srv.urls = append(srv.urls, url)
	}
}

func Logger(logger zerolog.Logger) func(*Server) {
	return func(srv *Server) {
		srv.logger = logger
	}
}

func NewServer(addr string, options ...func(*Server)) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &Server{
		srv:    &http.Server{Addr: addr, Handler: mux},
		mux:    mux,
		logger: log.Logger,
		urls:   make([]string, 0, 0),
	}

	for _, option := range options {
		option(srv)
	}

	return srv
}

func (srv *Server) Start() error {
	srv.logger.Info().
		Str("addr", srv.srv.Addr).
		Msg("Starting server")

	return nil
}
