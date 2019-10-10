package togo

import (
	"fmt"
	"net/http"
	"strings"
)

type Resource struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func (r Resource) SanitizedPath(prefix string) string {
	return fmt.Sprintf(
		"/%s/%s",
		strings.TrimPrefix(strings.TrimSuffix(prefix, "/"), "/"),
		strings.TrimPrefix(strings.TrimSuffix(r.Path, "/"), "/"),
	)
}
