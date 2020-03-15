package togo

import (
	"testing"

	"github.com/allisson/go-assert"
)

func TestResource(t *testing.T) {
	var (
		resource     = Resource{}
		expectedPath = "/svc/togo/cats/"
	)

	t.Run("sanitizing slashes", func(t *testing.T) {
		resource.Path = "/cats/"
		assert.Equal(t, expectedPath, resource.SanitizedPath("/svc/togo/"))
	})

	t.Run("sanitizing prefix", func(t *testing.T) {
		resource.Path = "/cats"
		assert.Equal(t, expectedPath, resource.SanitizedPath("/svc/togo"))
	})

	t.Run("sanitizing suffix", func(t *testing.T) {
		resource.Path = "cats/"
		assert.Equal(t, expectedPath, resource.SanitizedPath("svc/togo/"))
	})
}
