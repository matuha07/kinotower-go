package film_repository

import (
	"database/sql"
	"fmt"

	core_database "github.com/matuha07/kinotower-go/src/internal/core/database"
	"github.com/matuha07/kinotower-go/src/internal/core/domain"
)

type FilmRepository interface {
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

type filmRepository struct {
	db core_database.Database
}

func NewFilmRepository(db core_database.Database) *filmRepository {
	return &filmRepository{db: db}
}

func (r *filmRepository) GetFilms(filter domain.Filter) (*domain.FilmList, error) {
	where := []string{"f.deleted_at IS NULL"}
	args := []any{}

	if filter.CategoryID > 0 {
		where = append(where, "EXISTS (SELECT 1 FROM categories_films cf WHERE cf.film_id = f.id AND cf.category_id = $"+itoa(len(args)+1)+")")
		args = append(args, filter.CategoryID)
	}
	if filter.CountryID > 0 {
		where = append(where, "f.country_id = $"+itoa(len(args)+1))
		args = append(args, filter.CountryID)
	}
	if filter.Search != "" {
		where = append(where, "f.name ILIKE $"+itoa(len(args)+1))
		args = append(args, "%"+filter.Search+"%")
	}

	total, err := r.countFilms(where, args)
	if err != nil {
		return nil, err
	}

	sortBy := "f.name"
	if filter.SortBy == "year" {
		sortBy = "f.year_of_issue"
	}
	if filter.SortBy == "rating" {
		sortBy = "rating_avg"
	}
	sortDir := "ASC"
	if filter.SortDir == "desc" {
		sortDir = "DESC"
	}

	args = append(args, filter.Limit(), filter.Offset())
	query := fmt.Sprintf(`
		SELECT f.id, f.name, f.duration, f.year_of_issue, f.age,
		       f.link_img, f.link_kinopoisk, f.link_video, f.created_at,
		       c.id AS country_id, c.name AS country_name,
		       COALESCE(ra.rating_avg, 0) AS rating_avg,
		       COALESCE(rv.review_count, 0) AS review_count
		FROM films f
		LEFT JOIN countries c ON c.id = f.country_id
		LEFT JOIN (
			SELECT film_id, AVG(ball)::float AS rating_avg
			FROM ratings
			GROUP BY film_id
		) ra ON ra.film_id = f.id
		LEFT JOIN (
			SELECT film_id, COUNT(*) AS review_count
			FROM reviews
			WHERE deleted_at IS NULL AND is_approved = TRUE
			GROUP BY film_id
		) rv ON rv.film_id = f.id
		WHERE %s
		ORDER BY %s %s, f.id ASC
		LIMIT $%d OFFSET $%d
	`, join(where, " AND "), sortBy, sortDir, len(args)-1, len(args))

	rows := []FilmRow{}
	if err := r.db.Select(&rows, query, args...); err != nil {
		return nil, err
	}

	films := make([]domain.Film, 0, len(rows))
	for _, row := range rows {
		film := rowToFilm(row)
		categories, err := r.getCategories(row.ID)
		if err != nil {
			return nil, err
		}
		film.Categories = categories
		films = append(films, film)
	}

	return &domain.FilmList{Page: filter.Page, Size: len(films), Total: total, Films: films}, nil
}

func (r *filmRepository) countFilms(where []string, args []any) (int, error) {
	query := "SELECT COUNT(*) FROM films f WHERE " + join(where, " AND ")
	var total int
	if err := r.db.Get(&total, query, args...); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *filmRepository) GetFilmByID(id int) (*domain.Film, error) {
	const query = `
		SELECT f.id, f.name, f.duration, f.year_of_issue, f.age,
		       f.link_img, f.link_kinopoisk, f.link_video, f.created_at,
		       c.id AS country_id, c.name AS country_name,
		       COALESCE(ra.rating_avg, 0) AS rating_avg,
		       COALESCE(rv.review_count, 0) AS review_count
		FROM films f
		LEFT JOIN countries c ON c.id = f.country_id
		LEFT JOIN (
			SELECT film_id, AVG(ball)::float AS rating_avg
			FROM ratings
			GROUP BY film_id
		) ra ON ra.film_id = f.id
		LEFT JOIN (
			SELECT film_id, COUNT(*) AS review_count
			FROM reviews
			WHERE deleted_at IS NULL AND is_approved = TRUE
			GROUP BY film_id
		) rv ON rv.film_id = f.id
		WHERE f.deleted_at IS NULL AND f.id = $1
	`

	var row FilmRow
	if err := r.db.Get(&row, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	filmValue := rowToFilm(row)
	film := &filmValue

	categories, err := r.getCategories(id)
	if err != nil {
		return nil, err
	}
	film.Categories = categories

	return film, nil
}

func (r *filmRepository) GetCategories() ([]domain.Category, error) {
	const query = `
		SELECT c.id, c.name,
		       pc.id AS parent_id,
		       pc.name AS parent_name,
		       COALESCE(fc.film_count, 0) AS film_count
		FROM categories c
		LEFT JOIN categories pc ON pc.id = c.parent_id
		LEFT JOIN (
			SELECT cf.category_id, COUNT(*) AS film_count
			FROM categories_films cf
			INNER JOIN films f ON f.id = cf.film_id AND f.deleted_at IS NULL
			GROUP BY cf.category_id
		) fc ON fc.category_id = c.id
		WHERE c.deleted_at IS NULL
		ORDER BY c.id ASC
	`

	rows := []categoryExtendedRow{}
	if err := r.db.Select(&rows, query); err != nil {
		return nil, err
	}

	out := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		filmCount := row.FilmCount
		item := domain.Category{ID: row.ID, Name: row.Name, FilmCount: &filmCount}

		out = append(out, item)
	}

	return out, nil
}

func (r *filmRepository) GetCountries() ([]domain.Country, error) {
	const query = `
		SELECT c.id, c.name, COALESCE(fc.film_count, 0) AS film_count
		FROM countries c
		LEFT JOIN (
			SELECT country_id, COUNT(*) AS film_count
			FROM films
			WHERE deleted_at IS NULL
			GROUP BY country_id
		) fc ON fc.country_id = c.id
		ORDER BY c.id ASC
	`

	rows := []countryExtendedRow{}
	if err := r.db.Select(&rows, query); err != nil {
		return nil, err
	}

	out := make([]domain.Country, 0, len(rows))
	for _, row := range rows {
		filmCount := row.FilmCount
		out = append(out, domain.Country{ID: row.ID, Name: row.Name, FilmCount: &filmCount})
	}
	return out, nil
}

func (r *filmRepository) GetGenders() ([]domain.Gender, error) {
	const query = `SELECT id, name FROM gender ORDER BY id ASC`
	rows := []domain.Gender{}
	if err := r.db.Select(&rows, query); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *filmRepository) GetFilmReviews(filmID int) ([]domain.FilmReview, error) {
	if _, err := r.GetFilmByID(filmID); err != nil {
		return nil, err
	}

	const query = `
		SELECT r.id, u.id AS user_id, u.fio AS user_fio, r.message, r.created_at
		FROM reviews r
		INNER JOIN users u ON u.id = r.user_id
		WHERE r.film_id = $1 AND r.deleted_at IS NULL AND r.is_approved = TRUE
		ORDER BY r.created_at DESC, r.id DESC
	`

	rows := []filmReviewRow{}
	if err := r.db.Select(&rows, query, filmID); err != nil {
		return nil, err
	}

	out := make([]domain.FilmReview, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.FilmReview{
			ID: row.ID,
			User: domain.UserShort{
				ID:  row.UserID,
				FIO: row.UserFIO,
			},
			Message:   row.Message,
			CreatedAt: row.CreatedAt,
		})
	}

	return out, nil
}

func (r *filmRepository) CreateFilm(input domain.FilmCreate) (*domain.Film, error) {
	const insert = `
		INSERT INTO films (name, country_id, duration, year_of_issue, age, link_img, link_kinopoisk, link_video)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	var id int
	var createdAt string
	if err := r.db.QueryRow(insert,
		input.Name,
		input.CountryID,
		input.Duration,
		input.YearOfIssue,
		input.Age,
		input.LinkImg,
		input.LinkKinopoisk,
		input.LinkVideo,
	).Scan(&id, &createdAt); err != nil {
		return nil, err
	}

	if len(input.CategoryIDs) > 0 {
		if err := r.setFilmCategories(id, input.CategoryIDs); err != nil {
			return nil, err
		}
	}

	return r.GetFilmByID(id)
}

func (r *filmRepository) UpdateFilm(id int, input domain.FilmUpdate) (*domain.Film, error) {
	set := []string{}
	args := []any{}

	if input.Name != nil {
		set = append(set, "name = $"+itoa(len(args)+1))
		args = append(args, *input.Name)
	}
	if input.CountryID != nil {
		set = append(set, "country_id = $"+itoa(len(args)+1))
		args = append(args, *input.CountryID)
	}
	if input.Duration != nil {
		set = append(set, "duration = $"+itoa(len(args)+1))
		args = append(args, *input.Duration)
	}
	if input.YearOfIssue != nil {
		set = append(set, "year_of_issue = $"+itoa(len(args)+1))
		args = append(args, *input.YearOfIssue)
	}
	if input.Age != nil {
		set = append(set, "age = $"+itoa(len(args)+1))
		args = append(args, *input.Age)
	}
	if input.LinkImg != nil {
		set = append(set, "link_img = $"+itoa(len(args)+1))
		args = append(args, *input.LinkImg)
	}
	if input.LinkKinopoisk != nil {
		set = append(set, "link_kinopoisk = $"+itoa(len(args)+1))
		args = append(args, *input.LinkKinopoisk)
	}
	if input.LinkVideo != nil {
		set = append(set, "link_video = $"+itoa(len(args)+1))
		args = append(args, *input.LinkVideo)
	}

	if len(set) > 0 {
		args = append(args, id)
		query := "UPDATE films SET " + join(set, ", ") + " WHERE id = $" + itoa(len(args))
		res, err := r.db.Exec(query, args...)
		if err != nil {
			return nil, err
		}
		if rows, err := res.RowsAffected(); err == nil && rows == 0 {
			return nil, ErrNotFound
		}
	}

	if input.CategoryIDs != nil {
		if err := r.setFilmCategories(id, *input.CategoryIDs); err != nil {
			return nil, err
		}
	}

	return r.GetFilmByID(id)
}

func (r *filmRepository) DeleteFilm(id int) error {
	const query = `UPDATE films SET deleted_at = NOW() WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return ErrNotFound
	}
	return nil
}
