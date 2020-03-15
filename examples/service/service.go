package service

import (
	"net/http"

	"github.com/opentogo/togo"
)

// Service is a togo.Service that serves cats
type Service struct {
}

// NewService creates a new instance of this service
func NewService() Service {
	return Service{}
}

// Prefix defines the URL prefix this service listens to
func (s Service) Prefix() string {
	return "/svc/togo"
}

// Middleware defines the middleware stack that requests to this service will run through
func (s Service) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		togo.Log.Println("This middleware is enabled.")
		next.ServeHTTP(w, r)
	}
}

// Resources defines the set of resources (URL path, HTTP method, and handler) that this service responds to
func (s Service) Resources() []togo.Resource {
	return []togo.Resource{
		{
			Path:    "/cats/{id:[0-9]+}",
			Method:  http.MethodGet,
			Handler: GetCats("my-dependency"),
		},
	}
}
