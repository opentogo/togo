package togo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/opentogo/logger"
	"github.com/opentogo/router"
)

// Log is the global logger for the server.
var Log = logger.NewLogger(os.Stdout, "", 0)

// Togo is a microservice application server. To use, you create an instance, passing a name and a configuration.
// You then register services to your instance, and then call `Run()`.
// Example:
//  t := togo.Init("my-togo", togo.LoadJSONFile("config.json"))
//  t.Register(myservice.New())
//
//  if err := t.Run(); err != nil {
//      fmt.Fprintf(os.Stderr, "%v", err)
//      os.Exit(1)
//  }
type Togo struct {
	appName     string
	logFilename string
	Router      *router.Router
	server      *http.Server
}

// Init will set up the name, logging, and the HTTP server.
func Init(appName string, config Config) *Togo {
	Log.SetPrefix(fmt.Sprintf("[%s] ", appName))

	return &Togo{
		appName:     appName,
		logFilename: config.LogFilename,
		Router:      &router.Router{},
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", config.HTTPAddr, config.HTTPPort),
			ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
		},
	}
}

// Register adds the specified togo.Service to the Togo server.
func (t *Togo) Register(service Service) {
	resources := service.Resources()
	resources = append(resources, Resource{
		Path:   "/hc",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if _, err := io.WriteString(w, fmt.Sprintf("ok-%s", t.appName)); err != nil {
				Log.Printf("Unable to write healthcheck response: %s", err)
			}
		},
	})
	for _, resource := range resources {
		t.Router.Handler(resource.Method, resource.SanitizedPath(service.Prefix()), service.Middleware(resource.Handler))
	}
	t.Router.NotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write(nil); err != nil {
			Log.Printf("Unable to write response for `not found` handler.")
		}
	})
	if t.logFilename != "" {
		logFile(t.logFilename)
	}
	t.server.Handler = Log.Handler(t.Router)
}

// Run will start the HTTP server for Togo.
func (t Togo) Run() error {
	Log.Printf("Running at %s\n", t.server.Addr)
	return t.server.ListenAndServe()
}

func logFile(filename string) {
	var (
		err error
		out io.Writer
	)
	if out, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		Log.Printf("Unable to opening file %q: %s", filename, err.Error())
		return
	}
	Log.SetOutput(out)
}
