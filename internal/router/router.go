package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router() *chi.Mux {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	return r
}
