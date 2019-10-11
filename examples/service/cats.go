package service

import (
	"fmt"
	"net/http"
)

// GetCats returns the cats with the specified dependency
func GetCats(anyDependency interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf("my cats with %q dependency", anyDependency.(string))
		w.Write([]byte(response))
	}
}
