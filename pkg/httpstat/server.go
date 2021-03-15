package httpstat

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	client   *http.Client
	srv      *http.Server
	mux      *http.ServeMux
	logger   zerolog.Logger
	urls     []string
	interval time.Duration
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

func Timeout(timeout time.Duration) func(*Server) {
	return func(srv *Server) {
		srv.client.Timeout = timeout
	}
}

func Interval(interval time.Duration) func(*Server) {
	return func(srv *Server) {
		srv.interval = interval
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
		client:   &http.Client{Timeout: 5 * time.Second},
		srv:      &http.Server{Addr: addr, Handler: mux},
		mux:      mux,
		logger:   log.Logger,
		urls:     make([]string, 0, 0),
		interval: 15 * time.Second,
	}

	for _, option := range options {
		option(srv)
	}

	return srv
}

func (srv *Server) Start() error {
	srv.logger.Info().
		Str("addr", srv.srv.Addr).
		Msg("Started server")

	for _, url := range srv.urls {
		go srv.scrapeURL(url)
	}

	return srv.srv.ListenAndServe()
}

func (srv *Server) scrapeURL(url string) {
	for {
		start := time.Now()
		resp, err := srv.client.Get(url)
		if err != nil {
			urlUp.WithLabelValues(url).Set(0)

			srv.logger.Error().
				Err(err).
				Str("url", url).
				Msg("Failed to scrape URL")
		}
		defer resp.Body.Close()

		duration := time.Since(start).Milliseconds()
		urlResponseMS.WithLabelValues(url).Observe(float64(duration))

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			urlUp.WithLabelValues(url).Set(1)
		} else {
			urlUp.WithLabelValues(url).Set(0)
		}

		time.Sleep(srv.interval)
	}
}
