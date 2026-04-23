package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r *Router) userRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/{id}", r.userHandler.GetByID)
	router.Put("/", r.userHandler.Update)
	router.Delete("/", r.userHandler.Delete)

	router.Post("/{user-id}/reviews", r.userHandler.CreateReview)
	router.Get("/{user-id}/reviews", r.userHandler.GetReviews)
	router.Delete("/{user-id}/reviews/{id}", r.userHandler.DeleteReview)

	router.Post("/{user-id}/ratings", r.userHandler.CreateRating)
	router.Get("/{user-id}/ratings", r.userHandler.GetRatings)
	router.Delete("/{user-id}/ratings/{id}", r.userHandler.DeleteRating)

	return router
}
