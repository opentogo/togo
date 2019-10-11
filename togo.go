package togo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

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
	server      *http.Server
}

// Init will set up the name, logging, and the HTTP server.
func Init(appName string, config Config) *Togo {
	Logger.SetPrefix(fmt.Sprintf("[%s] ", appName))

	return &Togo{
		appName:     appName,
		logFilename: config.LogFilename,
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.HTTPAddr, config.HTTPPort),
			ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
		},
	}
}

// Register adds the specified togo.Service to the Togo server.
func (t *Togo) Register(service Service) {
	mux := http.NewServeMux()

	for _, resource := range service.Resources() {
		resource.Handler = service.Middleware(resource.Handler)

		mux.HandleFunc(resource.SanitizedPath(service.Prefix()), func(w http.ResponseWriter, r *http.Request) {
			if resource.Method == r.Method {
				resource.Handler.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(nil)
		})
	}

	if t.logFilename != "" {
		logFile(t.logFilename)
	}
	t.server.Handler = loggingHandler(mux)
}

// Run will start the HTTP server for Togo.
func (t Togo) Run() error {
	Logger.Printf("Running at %s\n", t.server.Addr)
	return t.server.ListenAndServe()
}

func logFile(filename string) {
	var (
		err error
		out io.Writer
	)
	if out, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		Logger.Printf("Unable to opening file %q: %s", filename, err.Error())
		return
	}
	Logger.SetOutput(out)
}
