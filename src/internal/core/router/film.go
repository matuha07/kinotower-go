package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r *Router) filmRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", r.filmHandler.GetFilms)
	router.Get("/{film-id}/reviews", r.filmHandler.GetFilmReviews)
	router.Post("/", r.filmHandler.CreateFilm)
	router.Route("/{id}", func(rl chi.Router) {
		rl.Get("/", r.filmHandler.GetFilmByID)
		rl.Put("/", r.filmHandler.UpdateFilm)
		rl.Delete("/", r.filmHandler.DeleteFilm)
	})
	return router
}
