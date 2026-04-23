package film_repository

import "time"

type Film struct {
	ID            int     `db:"id"`
	Name          string  `db:"name"`
	Duration      int     `db:"duration"`
	YearOfIssue   int     `db:"year_of_issue"`
	Age           int     `db:"age"`
	LinkImg       string  `db:"link_img"`
	LinkKinopoisk string  `db:"link_kinopoisk"`
	LinkVideo     string  `db:"link_video"`
	CreatedAt     string  `db:"created_at"`
	CountryID     int     `db:"country_id"`
	RatingAvg     float64 `db:"rating_avg"`
	ReviewCount   int     `db:"review_count"`
}

type FilmRow struct {
	ID            int     `db:"id"`
	Name          string  `db:"name"`
	Duration      int     `db:"duration"`
	YearOfIssue   int     `db:"year_of_issue"`
	Age           int     `db:"age"`
	LinkImg       *string `db:"link_img"`
	LinkKinopoisk *string `db:"link_kinopoisk"`
	LinkVideo     *string `db:"link_video"`
	CreatedAt     string  `db:"created_at"`
	CountryID     *int    `db:"country_id"`
	CountryName   *string `db:"country_name"`
	RatingAvg     float64 `db:"rating_avg"`
	ReviewCount   int     `db:"review_count"`
}

type CategoryRow struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type categoryExtendedRow struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	ParentID   *int    `db:"parent_id"`
	ParentName *string `db:"parent_name"`
	FilmCount  int     `db:"film_count"`
}

type countryExtendedRow struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	FilmCount int    `db:"film_count"`
}

type filmReviewRow struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	UserFIO   string    `db:"user_fio"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}
