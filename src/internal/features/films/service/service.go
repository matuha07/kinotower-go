package film_service

import (
	"errors"

	"github.com/matuha07/kinotower-go/src/internal/core/domain"
	film_repository "github.com/matuha07/kinotower-go/src/internal/features/films/repository"
)

type FilmService interface {
	GetFilms(filter domain.Filter) (*domain.FilmList, error)
	GetFilmByID(id int) (*domain.Film, error)
	GetCategories() ([]domain.Category, error)
	GetCountries() ([]domain.Country, error)
	GetGenders() ([]domain.Gender, error)
	GetFilmReviews(filmID int) ([]domain.FilmReview, error)
	CreateFilm(input domain.FilmCreate) (*domain.Film, error)
	UpdateFilm(id int, input domain.FilmUpdate) (*domain.Film, error)
	DeleteFilm(id int) error
}

type filmService struct {
	filmRepository film_repository.FilmRepository
}

var ErrNotFound = errors.New("film not found")

func NewFilmService(filmRepo film_repository.FilmRepository) *filmService {
	return &filmService{
		filmRepository: filmRepo,
	}
}

func (s *filmService) GetFilms(filter domain.Filter) (*domain.FilmList, error) {
	return s.filmRepository.GetFilms(filter)
}

func (s *filmService) GetFilmByID(id int) (*domain.Film, error) {
	film, err := s.filmRepository.GetFilmByID(id)
	if err != nil {
		if errors.Is(err, film_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return film, nil
}

func (s *filmService) GetCategories() ([]domain.Category, error) {
	return s.filmRepository.GetCategories()
}

func (s *filmService) GetCountries() ([]domain.Country, error) {
	return s.filmRepository.GetCountries()
}

func (s *filmService) GetGenders() ([]domain.Gender, error) {
	return s.filmRepository.GetGenders()
}

func (s *filmService) GetFilmReviews(filmID int) ([]domain.FilmReview, error) {
	reviews, err := s.filmRepository.GetFilmReviews(filmID)
	if err != nil {
		if errors.Is(err, film_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return reviews, nil
}

func (s *filmService) CreateFilm(input domain.FilmCreate) (*domain.Film, error) {
	return s.filmRepository.CreateFilm(input)
}

func (s *filmService) UpdateFilm(id int, input domain.FilmUpdate) (*domain.Film, error) {
	film, err := s.filmRepository.UpdateFilm(id, input)
	if err != nil {
		if errors.Is(err, film_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return film, nil
}

func (s *filmService) DeleteFilm(id int) error {
	if err := s.filmRepository.DeleteFilm(id); err != nil {
		if errors.Is(err, film_repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
