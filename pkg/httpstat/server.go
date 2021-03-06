package httpstat

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Server is the server component of httpstat-monitor that performs the actual
// monitoring of URLs.
type Server struct {
	client   *http.Client
	srv      *http.Server
	mux      *http.ServeMux
	logger   zerolog.Logger
	urls     []string
	interval time.Duration
}

// MonitorURL adds a URL to the list of URLs to monitor.
func MonitorURL(url string) func(*Server) {
	return func(srv *Server) {
		srv.urls = append(srv.urls, url)
	}
}

// Logger sets the logger of the server.
func Logger(logger zerolog.Logger) func(*Server) {
	return func(srv *Server) {
		srv.logger = logger
	}
}

// Timeout sets the timeout of the HTTP client used to monitor URLs.
func Timeout(timeout time.Duration) func(*Server) {
	return func(srv *Server) {
		srv.client.Timeout = timeout
	}
}

// Interval sets the interval between calls to each URL.
func Interval(interval time.Duration) func(*Server) {
	return func(srv *Server) {
		srv.interval = interval
	}
}

// NewServer creates a new Server instance with which to monitor URLs.
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
		urls:     make([]string, 0),
		interval: 15 * time.Second,
	}

	for _, option := range options {
		option(srv)
	}

	return srv
}

// Starts starts the server.
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
	srv.logger.Debug().
		Str("url", url).
		Msg("Scraping URL")

	for {
		timer := prometheus.NewTimer(urlResponseMS.WithLabelValues(url))

		resp, err := srv.client.Get(url)
		if err != nil {
			urlUp.WithLabelValues(url).Set(0)

			srv.logger.Error().
				Err(err).
				Str("url", url).
				Msg("Failed to scrape URL")
		}

		timer.ObserveDuration()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			urlUp.WithLabelValues(url).Set(1)
		} else {
			urlUp.WithLabelValues(url).Set(0)
		}

		resp.Body.Close()

		time.Sleep(srv.interval)
	}
}
