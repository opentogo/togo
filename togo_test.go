package togo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/allisson/go-assert"
)

type fakeService struct{}

func (s fakeService) Prefix() string {
	return "/svc/togo"
}

func (s fakeService) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Log.Println("This middleware was called")
		next.ServeHTTP(w, r)
	}
}

func (s fakeService) Resources() []Resource {
	return []Resource{
		{
			Path:   "/test",
			Method: http.MethodGet,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("testing")) // nolint: errcheck
			},
		},
	}
}

var (
	buf     bytes.Buffer
	service = Init("togo-testing", Config{
		HTTPAddr:     "0.0.0.0",
		HTTPPort:     3000,
		IdleTimeout:  30,
		ReadTimeout:  5,
		WriteTimeout: 10,
		LogFilename:  "/var/log/svc/togo",
	})
)

func TestTogo(t *testing.T) {
	Log.SetOutput(&buf)
	service.Register(fakeService{})

	t.Run("checking builder func value", func(t *testing.T) {
		assert.Equal(t, service.appName, "togo-testing")
		assert.Equal(t, service.logFilename, "/var/log/svc/togo")
		assert.Equal(t, service.server.Addr, "0.0.0.0:3000")
		assert.Equal(t, service.server.IdleTimeout, time.Duration(30)*time.Second)
		assert.Equal(t, service.server.ReadTimeout, time.Duration(5)*time.Second)
		assert.Equal(t, service.server.WriteTimeout, time.Duration(10)*time.Second)
	})

	t.Run("checking invalid log filename", func(t *testing.T) {
		assert.Equal(t, "[togo-testing] Unable to opening file \"/var/log/svc/togo\": open /var/log/svc/togo: no such file or directory\n", buf.String())
	})

	t.Run("checking unregistered endpoints", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/svc/togo/invalid/", nil)
		)

		service.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("checking registered endpoints", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/svc/togo/test/", nil)
		)

		service.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "testing", w.Body.String())
	})

	t.Run("requesting an endpoint via a not allowed method", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodPost, "/svc/togo/test/", nil)
		)

		service.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("checking middleware call", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/svc/togo/test/", nil)
		)

		buf.Reset()
		service.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, true, strings.Contains(buf.String(), "This middleware was called"))
	})

	t.Run("checking healthcheck endpoint", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/svc/togo/hc/", nil)
		)

		service.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "ok-togo-testing", w.Body.String())
	})
}
