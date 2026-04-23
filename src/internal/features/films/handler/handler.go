package film_handler

import (
	film_service "github.com/matuha07/kinotower-go/src/internal/features/films/service"
)

type FilmHandler struct {
	filmService film_service.FilmService
}

func NewFilmHandler(filmService film_service.FilmService) *FilmHandler {
	return &FilmHandler{filmService: filmService}
}
