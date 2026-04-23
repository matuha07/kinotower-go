package user_repository

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	core_database "github.com/matuha07/kinotower-go/src/internal/core/database"
	"github.com/matuha07/kinotower-go/src/internal/core/domain"
)

var ErrNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already taken")
var ErrReviewNotFound = errors.New("review not found")
var ErrRatingNotFound = errors.New("rating not found")
var ErrFilmNotFound = errors.New("film not found")
var ErrScoreExists = errors.New("score exist")

type UserRepository interface {
	GetByID(id int) (*domain.User, error)
	GetByEmail(email string) (*userRow, error)
	Create(input domain.UserCreate, hashedPassword string) (*domain.User, error)
	Update(id int, input domain.UserUpdate) (*domain.User, error)
	Delete(id int) error
	CreateReview(userID int, input domain.ReviewCreate) (*domain.UserReview, error)
	GetReviews(userID int) ([]domain.UserReview, error)
	DeleteReview(userID, reviewID int) error
	CreateRating(userID int, input domain.RatingCreate) (*domain.UserRating, error)
	GetRatings(userID int) ([]domain.UserRating, error)
	DeleteRating(userID, ratingID int) error
}

type userRepository struct {
	db core_database.Database
}

func NewUserRepository(db core_database.Database) UserRepository {
	return &userRepository{db: db}
}

type userRow struct {
	ID          int       `db:"id"`
	FIO         string    `db:"fio"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Birthday    string    `db:"birthday"`
	CountryID   *int      `db:"country_id"`
	GenderID    *int      `db:"gender_id"`
	GenderName  *string   `db:"gender_name"`
	ReviewCount int       `db:"review_count"`
	RatingCount int       `db:"rating_count"`
	CreatedAt   time.Time `db:"created_at"`
}

type userReviewRow struct {
	ID         int       `db:"id"`
	FilmID     int       `db:"film_id"`
	FilmName   string    `db:"film_name"`
	Message    string    `db:"message"`
	IsApproved bool      `db:"is_approved"`
	CreatedAt  time.Time `db:"created_at"`
}

type userRatingRow struct {
	ID        int       `db:"id"`
	FilmID    int       `db:"film_id"`
	FilmName  string    `db:"film_name"`
	Score     int       `db:"score"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	const q = `
		SELECT u.id, u.fio, u.email, u.password, u.birthday::text, u.country_id, u.gender_id,
		       g.name AS gender_name,
		       COALESCE(rv.review_count, 0) AS review_count,
		       COALESCE(rt.rating_count, 0) AS rating_count,
		       u.created_at
		FROM users u
		LEFT JOIN gender g ON g.id = u.gender_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) AS review_count
			FROM reviews
			WHERE deleted_at IS NULL
			GROUP BY user_id
		) rv ON rv.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(*) AS rating_count
			FROM ratings
			GROUP BY user_id
		) rt ON rt.user_id = u.id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`
	var row userRow
	if err := r.db.Get(&row, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return rowToUser(row), nil
}

func (r *userRepository) GetByEmail(email string) (*userRow, error) {
	const q = `
		SELECT id, fio, email, password, birthday::text, country_id, gender_id, NULL::text AS gender_name,
		       0 AS review_count, 0 AS rating_count, created_at
		FROM users WHERE email = $1 AND deleted_at IS NULL
	`
	var row userRow
	if err := r.db.Get(&row, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &row, nil
}

func (r *userRepository) Create(input domain.UserCreate, hashedPassword string) (*domain.User, error) {
	const q = `
		INSERT INTO users (fio, email, password, birthday, country_id, gender_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, fio, email, password, birthday::text, country_id, gender_id, NULL::text AS gender_name,
		          0 AS review_count, 0 AS rating_count, created_at
	`
	var row userRow
	err := r.db.QueryRowx(q,
		input.FIO, input.Email, hashedPassword, input.Birthday, input.CountryID, input.GenderID,
	).StructScan(&row)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return rowToUser(row), nil
}

func (r *userRepository) Update(id int, input domain.UserUpdate) (*domain.User, error) {
	set := []string{}
	args := []any{}

	if input.FIO != nil {
		set = append(set, "fio = $"+itoa(len(args)+1))
		args = append(args, *input.FIO)
	}
	if input.Birthday != nil {
		set = append(set, "birthday = $"+itoa(len(args)+1))
		args = append(args, *input.Birthday)
	}
	if input.Email != nil {
		set = append(set, "email = $"+itoa(len(args)+1))
		args = append(args, *input.Email)
	}
	if input.CountryID != nil {
		set = append(set, "country_id = $"+itoa(len(args)+1))
		args = append(args, *input.CountryID)
	}
	if input.GenderID != nil {
		set = append(set, "gender_id = $"+itoa(len(args)+1))
		args = append(args, *input.GenderID)
	}

	if len(set) == 0 {
		return r.GetByID(id)
	}

	args = append(args, id)
	query := "UPDATE users SET " + strings.Join(set, ", ") + " WHERE id = $" + itoa(len(args)) + " AND deleted_at IS NULL"
	res, err := r.db.Exec(query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, ErrNotFound
	}
	return r.GetByID(id)
}

func (r *userRepository) Delete(id int) error {
	const q = `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.Exec(q, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func rowToUser(row userRow) *domain.User {
	var gender *domain.Gender
	if row.GenderID != nil {
		gender = &domain.Gender{ID: *row.GenderID, Name: safeString(row.GenderName)}
	}

	return &domain.User{
		ID:          row.ID,
		FIO:         row.FIO,
		Email:       row.Email,
		Birthday:    row.Birthday,
		CountryID:   row.CountryID,
		GenderID:    row.GenderID,
		Gender:      gender,
		CreatedAt:   row.CreatedAt,
		ReviewCount: row.ReviewCount,
		RatingCount: row.RatingCount,
	}
}

func (r *userRepository) CreateReview(userID int, input domain.ReviewCreate) (*domain.UserReview, error) {
	if err := r.ensureUser(userID); err != nil {
		return nil, err
	}
	filmName, err := r.getFilmName(input.FilmID)
	if err != nil {
		return nil, err
	}

	const query = `
		INSERT INTO reviews (film_id, user_id, message, is_approved)
		VALUES ($1, $2, $3, FALSE)
		RETURNING id, created_at
	`

	var id int
	var createdAt time.Time
	if err := r.db.QueryRow(query, input.FilmID, userID, input.Message).Scan(&id, &createdAt); err != nil {
		return nil, err
	}

	return &domain.UserReview{
		ID: id,
		Film: domain.FilmShort{
			ID:   input.FilmID,
			Name: filmName,
		},
		Message:    input.Message,
		IsApproved: 0,
		CreatedAt:  createdAt,
	}, nil
}

func (r *userRepository) GetReviews(userID int) ([]domain.UserReview, error) {
	if err := r.ensureUser(userID); err != nil {
		return nil, err
	}

	const query = `
		SELECT r.id, f.id AS film_id, f.name AS film_name, r.message, r.is_approved, r.created_at
		FROM reviews r
		INNER JOIN films f ON f.id = r.film_id
		WHERE r.user_id = $1 AND r.deleted_at IS NULL
		ORDER BY r.created_at DESC, r.id DESC
	`

	rows := []userReviewRow{}
	if err := r.db.Select(&rows, query, userID); err != nil {
		return nil, err
	}

	out := make([]domain.UserReview, 0, len(rows))
	for _, row := range rows {
		approved := 0
		if row.IsApproved {
			approved = 1
		}
		out = append(out, domain.UserReview{
			ID: row.ID,
			Film: domain.FilmShort{
				ID:   row.FilmID,
				Name: row.FilmName,
			},
			Message:    row.Message,
			IsApproved: approved,
			CreatedAt:  row.CreatedAt,
		})
	}
	return out, nil
}

func (r *userRepository) DeleteReview(userID, reviewID int) error {
	if err := r.ensureUser(userID); err != nil {
		return err
	}
	const query = `UPDATE reviews SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	res, err := r.db.Exec(query, reviewID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrReviewNotFound
	}
	return nil
}

func (r *userRepository) CreateRating(userID int, input domain.RatingCreate) (*domain.UserRating, error) {
	if err := r.ensureUser(userID); err != nil {
		return nil, err
	}
	filmName, err := r.getFilmName(input.FilmID)
	if err != nil {
		return nil, err
	}

	const existing = `SELECT COUNT(*) FROM ratings WHERE film_id = $1 AND user_id = $2`
	var count int
	if err := r.db.Get(&count, existing, input.FilmID, userID); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrScoreExists
	}

	const query = `
		INSERT INTO ratings (film_id, user_id, ball)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	var id int
	var createdAt time.Time
	if err := r.db.QueryRow(query, input.FilmID, userID, input.Ball).Scan(&id, &createdAt); err != nil {
		return nil, err
	}

	return &domain.UserRating{
		ID: id,
		Film: domain.FilmShort{
			ID:   input.FilmID,
			Name: filmName,
		},
		Score:     input.Ball,
		CreatedAt: createdAt,
	}, nil
}

func (r *userRepository) GetRatings(userID int) ([]domain.UserRating, error) {
	if err := r.ensureUser(userID); err != nil {
		return nil, err
	}
	const query = `
		SELECT rt.id, f.id AS film_id, f.name AS film_name, rt.ball AS score, rt.created_at
		FROM ratings rt
		INNER JOIN films f ON f.id = rt.film_id
		WHERE rt.user_id = $1
		ORDER BY rt.created_at DESC, rt.id DESC
	`
	rows := []userRatingRow{}
	if err := r.db.Select(&rows, query, userID); err != nil {
		return nil, err
	}
	out := make([]domain.UserRating, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.UserRating{
			ID: row.ID,
			Film: domain.FilmShort{
				ID:   row.FilmID,
				Name: row.FilmName,
			},
			Score:     row.Score,
			CreatedAt: row.CreatedAt,
		})
	}
	return out, nil
}

func (r *userRepository) DeleteRating(userID, ratingID int) error {
	if err := r.ensureUser(userID); err != nil {
		return err
	}
	const query = `DELETE FROM ratings WHERE id = $1 AND user_id = $2`
	res, err := r.db.Exec(query, ratingID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrRatingNotFound
	}
	return nil
}

func (r *userRepository) ensureUser(userID int) error {
	const q = `SELECT COUNT(*) FROM users WHERE id = $1 AND deleted_at IS NULL`
	var count int
	if err := r.db.Get(&count, q, userID); err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepository) getFilmName(filmID int) (string, error) {
	const q = `SELECT name FROM films WHERE id = $1 AND deleted_at IS NULL`
	var name string
	if err := r.db.Get(&name, q, filmID); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrFilmNotFound
		}
		return "", err
	}
	return name, nil
}

func safeString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func itoa(v int) string { return strconv.Itoa(v) }
