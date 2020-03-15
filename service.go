package togo

import "net/http"

// Service is the basic interface that defines what to expect from any service.
type Service interface {
	Prefix() string
	Middleware(http.HandlerFunc) http.HandlerFunc
	Resources() []Resource
}
