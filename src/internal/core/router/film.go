package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r *Router) filmRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", r.filmHandler.GetFilms)
	return router
}
