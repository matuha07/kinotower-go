package auth_service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/matuha07/kinotower-go/src/internal/core/domain"
	user_repository "github.com/matuha07/kinotower-go/src/internal/features/users/repository"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrEmailTaken = errors.New("email already taken")

type AuthService interface {
	SignUp(input domain.SignUpInput) (*domain.TokenPair, error)
	SignIn(input domain.SignInInput) (*domain.TokenPair, error)
}

type authService struct {
	userRepo  user_repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo user_repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) SignUp(input domain.SignUpInput) (*domain.TokenPair, error) {
	hashed := hashPassword(input.Password)
	user, err := s.userRepo.Create(domain.UserCreate{
		FIO:       input.FIO,
		Email:     input.Email,
		Password:  hashed,
		Birthday:  input.Birthday,
		CountryID: input.CountryID,
		GenderID:  input.GenderID,
	}, hashed)
	if err != nil {
		if errors.Is(err, user_repository.ErrEmailTaken) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.TokenPair{
		Status: "success",
		Token:  token,
		ID:     user.ID,
		FIO:    user.FIO,
	}, nil
}

func (s *authService) SignIn(input domain.SignInInput) (*domain.TokenPair, error) {
	row, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if row.Password != hashPassword(input.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(row.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.TokenPair{Status: "success", Token: token, ID: row.ID, FIO: row.FIO}, nil
}

func (s *authService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// hashPassword uses SHA-256 for simplicity (replace with bcrypt for production).
func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}
