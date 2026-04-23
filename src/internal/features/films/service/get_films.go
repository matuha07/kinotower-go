package film_service

import "github.com/matuha07/kinotower-go/src/internal/features/films/domain"

func (s *filmService) GetFilms() ([]domain.Film, error) {
	return []domain.Film{}, nil
}
