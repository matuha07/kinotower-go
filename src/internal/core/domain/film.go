package domain

import "time"

type Film struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Duration      int        `json:"duration"`
	YearOfIssue   int        `json:"year_of_issue"`
	Age           int        `json:"age"`
	LinkImg       *string    `json:"link_img"`
	LinkKinopoisk *string    `json:"link_kinopoisk"`
	LinkVideo     *string    `json:"link_video"`
	CreatedAt     time.Time  `json:"created_at"`
	Country       *Country   `json:"country,omitempty"`
	Categories    []Category `json:"categories,omitempty"`
	RatingAvg     float64    `json:"ratingAvg"`
	ReviewCount   int        `json:"reviewCount"`
}

type FilmList struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int    `json:"total"`
	Films []Film `json:"films"`
}

type FilmCreate struct {
	Name          string
	CountryID     *int
	Duration      int
	YearOfIssue   int
	Age           int
	LinkImg       *string
	LinkKinopoisk *string
	LinkVideo     *string
	CategoryIDs   []int
}

type FilmUpdate struct {
	Name          *string
	CountryID     *int
	Duration      *int
	YearOfIssue   *int
	Age           *int
	LinkImg       *string
	LinkKinopoisk *string
	LinkVideo     *string
	CategoryIDs   *[]int
}
