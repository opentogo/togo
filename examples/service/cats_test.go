package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allisson/go-assert"
	"github.com/opentogo/router"
)

func TestService_GetCats(t *testing.T) {
	var (
		router = &router.Router{}
		r      = httptest.NewRequest(http.MethodGet, "/svc/togo/cats/123", nil)
		w      = httptest.NewRecorder()
	)

	router.Handler(http.MethodGet, "/svc/togo/cats/:id", GetCats("none"))
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `my "123" cat with "none" dependency`, w.Body.String())
}
