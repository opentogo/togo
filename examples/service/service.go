package service

import (
	"net/http"

	"github.com/opentogo/togo"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Prefix() string {
	return "/svc/togo"
}

func (s Service) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		togo.Logger.Println("This middleware is enabled.")
		next.ServeHTTP(w, r)
	}
}

func (s Service) Resources() []togo.Resource {
	return []togo.Resource{
		{
			Path:    "/cats",
			Method:  http.MethodGet,
			Handler: GetCats("my-dependency"),
		},
	}
}
