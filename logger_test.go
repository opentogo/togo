package togo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/allisson/go-assert"
)

func TestLogger(t *testing.T) {
	var (
		mux = http.NewServeMux()
		r   = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		w   = httptest.NewRecorder()
	)

	r.RemoteAddr = "127.0.0.1"
	now = func() time.Time {
		t, _ := time.Parse("2006-01-02T15:04:05", "1980-01-01T12:12:12")
		return t
	}

	Logger.SetOutput(&buf)
	t.Run("logging without header data", func(t *testing.T) {
		buf.Reset()

		loggingHandler(mux)(w, r)
		assert.Equal(t, "127.0.0.1 - - [01/Jan/1980:12:12:12 +0000] \"GET http://example.com HTTP/1.1\" 301 54 \"-\" \"-\" 0.0000\n", buf.String())
	})

	t.Run("logging with referer and user-agent data", func(t *testing.T) {
		buf.Reset()

		r.Header.Set("Referer", "http://example.com")
		r.Header.Set("User-Agent", "go/testing")

		loggingHandler(mux)(w, r)
		assert.Equal(t, "127.0.0.1 - - [01/Jan/1980:12:12:12 +0000] \"GET http://example.com HTTP/1.1\" 301 0 \"http://example.com\" \"go/testing\" 0.0000\n", buf.String())
	})

	t.Run("logging with encoded request", func(t *testing.T) {
		buf.Reset()

		r.URL, _ = url.Parse("http://example.com/testing?message=foo%20bar&go=lang%3F")

		loggingHandler(mux)(w, r)
		assert.Equal(t, "127.0.0.1 - - [01/Jan/1980:12:12:12 +0000] \"GET http://example.com HTTP/1.1\" 404 19 \"http://example.com\" \"go/testing\" 0.0000\n", buf.String())
	})

	t.Run("logging with connect method", func(t *testing.T) {
		buf.Reset()

		r.Method = http.MethodConnect
		r.Host = "www.example.com:443"
		r.Proto = "HTTP/2.0"
		r.ProtoMajor = 2
		r.ProtoMinor = 0
		r.URL = &url.URL{Host: "www.example.com:443"}

		loggingHandler(mux)(w, r)
		assert.Equal(t, "127.0.0.1 - - [01/Jan/1980:12:12:12 +0000] \"CONNECT www.example.com:443 HTTP/2.0\" 404 19 \"http://example.com\" \"go/testing\" 0.0000\n", buf.String())
	})
}
