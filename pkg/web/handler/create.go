package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/create" {
			handleNotFound(w, r)
			return
		}
		routes.CreatePage().Render(r.Context(), w)
	}
}
