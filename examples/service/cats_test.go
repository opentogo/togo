package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allisson/go-assert"
)

func TestService_GetCats(t *testing.T) {
	var (
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		w = httptest.NewRecorder()
	)

	http.HandlerFunc(GetCats("none")).ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.String(), `my cats with "none" dependency`)
}
