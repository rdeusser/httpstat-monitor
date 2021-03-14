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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
  <head><title>HTTP Stat Monitor</title></head>
  <body>
    <h1>HTTP Stat Monitor</h1>
    <p><a href='/metrics'>Metrics</a></p>
  </body>
</html>`))
	})

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

	return srv.srv.ListenAndServe()
}
