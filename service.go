package togo

import "net/http"

type Service interface {
	Prefix() string
	Middleware(http.HandlerFunc) http.HandlerFunc
	Resources() []Resource
}
