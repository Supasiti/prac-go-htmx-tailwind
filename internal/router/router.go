package router

import (
	"net/http"

	"github.com/supasiti/prac-go-htmx-tailwind/components/page"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		page.Page().Render(r.Context(), w)
	})

	return mux
}
