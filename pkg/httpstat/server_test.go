package httpstat

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("httpstat 200", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "200 OK")
		}))
		defer ts.Close()

		addr := getAddr(t, ts.URL)
		srv := NewServer(addr, MonitorURL(ts.URL))

		go func() {
			err := srv.Start()
			assert.NoError(t, err)
		}()

		time.Sleep(2 * time.Second)

		resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
		assert.NoError(t, err)

		data, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.True(t, strings.Contains(string(data), fmt.Sprintf(`sample_external_url_up{url="%s"}`, ts.URL)))
	})

	t.Run("httpstat 503", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "503 Service Unavailable")
		}))
		defer ts.Close()

		addr := getAddr(t, ts.URL)
		srv := NewServer(addr, MonitorURL(ts.URL))

		go func() {
			err := srv.Start()
			assert.NoError(t, err)
		}()

		time.Sleep(2 * time.Second)

		resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
		assert.NoError(t, err)

		data, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.True(t, strings.Contains(string(data), fmt.Sprintf(`sample_external_url_up{url="%s"}`, ts.URL)))
	})
}

func getAddr(t *testing.T, uri string) string {
	t.Helper()

	u, err := url.Parse(uri)
	assert.NoError(t, err)

	port, err := strconv.Atoi(u.Port())
	assert.NoError(t, err)

	port = port + 1 // increment port number

	return fmt.Sprintf("%s:%d", u.Hostname(), port)
}
