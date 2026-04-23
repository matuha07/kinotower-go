package user_service

import (
	"errors"

	"github.com/matuha07/kinotower-go/src/internal/core/domain"
	user_repository "github.com/matuha07/kinotower-go/src/internal/features/users/repository"
)

var ErrNotFound = errors.New("user not found")
var ErrReviewNotFound = errors.New("review not found")
var ErrRatingNotFound = errors.New("rating not found")
var ErrFilmNotFound = errors.New("film not found")
var ErrScoreExists = errors.New("score exist")
var ErrEmailTaken = errors.New("email already taken")

type UserService interface {
	GetByID(id int) (*domain.User, error)
	Update(id int, input domain.UserUpdate) (*domain.User, error)
	Delete(id int) error
	CreateReview(userID int, input domain.ReviewCreate) (*domain.UserReview, error)
	GetReviews(userID int) ([]domain.UserReview, error)
	DeleteReview(userID, reviewID int) error
	CreateRating(userID int, input domain.RatingCreate) (*domain.UserRating, error)
	GetRatings(userID int) ([]domain.UserRating, error)
	DeleteRating(userID, ratingID int) error
}

type userService struct {
	repo user_repository.UserRepository
}

func NewUserService(repo user_repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetByID(id int) (*domain.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) Update(id int, input domain.UserUpdate) (*domain.User, error) {
	u, err := s.repo.Update(id, input)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, user_repository.ErrEmailTaken) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *userService) CreateReview(userID int, input domain.ReviewCreate) (*domain.UserReview, error) {
	review, err := s.repo.CreateReview(userID, input)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, user_repository.ErrFilmNotFound) {
			return nil, ErrFilmNotFound
		}
		return nil, err
	}
	return review, nil
}

func (s *userService) GetReviews(userID int) ([]domain.UserReview, error) {
	reviews, err := s.repo.GetReviews(userID)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return reviews, nil
}

func (s *userService) DeleteReview(userID, reviewID int) error {
	err := s.repo.DeleteReview(userID, reviewID)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, user_repository.ErrReviewNotFound) {
			return ErrReviewNotFound
		}
		return err
	}
	return nil
}

func (s *userService) CreateRating(userID int, input domain.RatingCreate) (*domain.UserRating, error) {
	rating, err := s.repo.CreateRating(userID, input)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, user_repository.ErrFilmNotFound) {
			return nil, ErrFilmNotFound
		}
		if errors.Is(err, user_repository.ErrScoreExists) {
			return nil, ErrScoreExists
		}
		return nil, err
	}
	return rating, nil
}

func (s *userService) GetRatings(userID int) ([]domain.UserRating, error) {
	ratings, err := s.repo.GetRatings(userID)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return ratings, nil
}

func (s *userService) DeleteRating(userID, ratingID int) error {
	err := s.repo.DeleteRating(userID, ratingID)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, user_repository.ErrRatingNotFound) {
			return ErrRatingNotFound
		}
		return err
	}
	return nil
}
