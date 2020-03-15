package service

import (
	"net/http"

	"github.com/opentogo/middlewares"
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
	middlewares.Use(
		middlewares.NewIPSpoofing(),
		middlewares.NewPathTraversal(),
		middlewares.NewRemoteReferer([]string{
			http.MethodDelete,
			http.MethodGet,
			http.MethodPatch,
			http.MethodPost,
			http.MethodPut,
		}),
		middlewares.NewXSS("block", true),
	)

	return middlewares.Handle(next)
}

// Resources defines the set of resources (URL path, HTTP method, and handler) that this service responds to
func (s Service) Resources() []togo.Resource {
	return []togo.Resource{
		{
			Path:    "/cats/:id",
			Method:  http.MethodGet,
			Handler: GetCats("my-dependency"),
		},
	}
}
