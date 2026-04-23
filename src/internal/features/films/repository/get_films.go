package film_repository

import "github.com/matuha07/kinotower-go/src/internal/features/films/domain"

func (r *filmRepository) GetFilms() ([]domain.Film, error) {
	return []domain.Film{}, nil
}
