package togo

import (
	"fmt"
	"net/http"
	"strings"
)

// Resource holds info required for configuring endpoints and resources for the server.
type Resource struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

// SanitizedPath trims forward slash characters from the provided string prefix
// and Path of the resource.
// Returns the sanitized path string.
func (r Resource) SanitizedPath(prefix string) string {
	return fmt.Sprintf(
		"/%s/%s/",
		strings.TrimPrefix(strings.TrimSuffix(prefix, "/"), "/"),
		strings.TrimPrefix(strings.TrimSuffix(r.Path, "/"), "/"),
	)
}
