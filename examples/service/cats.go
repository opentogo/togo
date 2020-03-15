package service

import (
	"fmt"
	"net/http"

	"github.com/opentogo/router"
)

// GetCats returns the cats with the specified dependency
func GetCats(anyDependency interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			params   = router.Params(r)
			response = fmt.Sprintf("my %q cat with %q dependency", params["id"], anyDependency)
		)
		if _, err := w.Write([]byte(response)); err != nil {
			http.Error(w, "error trying to return response", http.StatusInternalServerError)
		}
	}
}
