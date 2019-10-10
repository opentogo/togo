package togo

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const logPattern = "%s - - [%s] \"%s %s %s\" %d %d %q %q %.4f\n"

var (
	// Logger is the global logger for the server.
	Logger = log.New(os.Stdout, "", 0)

	now = func() time.Time {
		return time.Now().UTC()
	}
)

type logWriter struct {
	w         http.ResponseWriter
	r         *http.Request
	timestamp time.Time
	status    int
	size      int
}

func (l *logWriter) Header() http.Header {
	return l.w.Header()
}

func (l *logWriter) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *logWriter) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *logWriter) Status() int {
	return l.status
}

func (l *logWriter) Size() int {
	return l.size
}

func (l *logWriter) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func loggingHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := now()

		writer := &logWriter{
			w:         w,
			r:         r,
			status:    http.StatusOK,
			timestamp: s,
		}

		handler.ServeHTTP(writer, r)
		if r.MultipartForm != nil {
			r.MultipartForm.RemoveAll()
		}

		apacheLogFormat(writer)
	}
}

func apacheLogFormat(writer *logWriter) {
	var (
		finishTime = now()
		uri        = writer.r.RequestURI
		referer    = writer.r.Referer()
		userAgent  = writer.r.UserAgent()
	)
	host, _, err := net.SplitHostPort(writer.r.RemoteAddr)
	if err != nil {
		host = writer.r.RemoteAddr
	}
	if writer.r.ProtoMajor == 2 && writer.r.Method == "CONNECT" {
		uri = writer.r.Host
	}
	if uri == "" {
		uri = writer.r.URL.RequestURI()
	}
	if referer == "" {
		referer = "-"
	}
	if userAgent == "" {
		userAgent = "-"
	}
	Logger.Printf(
		logPattern,
		host,
		finishTime.Format("02/Jan/2006:15:04:05 -0700"),
		writer.r.Method,
		uri,
		writer.r.Proto,
		writer.Status(),
		writer.Size(),
		referer,
		userAgent,
		finishTime.Sub(writer.timestamp).Seconds(),
	)
}
