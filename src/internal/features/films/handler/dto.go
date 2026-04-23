package film_handler

import "github.com/matuha07/kinotower-go/src/internal/core/domain"

type createFilmRequest struct {
	Name          string  `json:"name"`
	CountryID     *int    `json:"country_id"`
	Duration      int     `json:"duration"`
	YearOfIssue   int     `json:"year_of_issue"`
	Age           int     `json:"age"`
	LinkImg       *string `json:"link_img"`
	LinkKinopoisk *string `json:"link_kinopoisk"`
	LinkVideo     *string `json:"link_video"`
	CategoryIDs   []int   `json:"category_ids"`
}

type updateFilmRequest struct {
	Name          *string `json:"name"`
	CountryID     *int    `json:"country_id"`
	Duration      *int    `json:"duration"`
	YearOfIssue   *int    `json:"year_of_issue"`
	Age           *int    `json:"age"`
	LinkImg       *string `json:"link_img"`
	LinkKinopoisk *string `json:"link_kinopoisk"`
	LinkVideo     *string `json:"link_video"`
	CategoryIDs   *[]int  `json:"category_ids"`
}

func (r createFilmRequest) toDomain() domain.FilmCreate {
	return domain.FilmCreate{
		Name:          r.Name,
		CountryID:     r.CountryID,
		Duration:      r.Duration,
		YearOfIssue:   r.YearOfIssue,
		Age:           r.Age,
		LinkImg:       r.LinkImg,
		LinkKinopoisk: r.LinkKinopoisk,
		LinkVideo:     r.LinkVideo,
		CategoryIDs:   r.CategoryIDs,
	}
}

func (r updateFilmRequest) toDomain() domain.FilmUpdate {
	return domain.FilmUpdate{
		Name:          r.Name,
		CountryID:     r.CountryID,
		Duration:      r.Duration,
		YearOfIssue:   r.YearOfIssue,
		Age:           r.Age,
		LinkImg:       r.LinkImg,
		LinkKinopoisk: r.LinkKinopoisk,
		LinkVideo:     r.LinkVideo,
		CategoryIDs:   r.CategoryIDs,
	}
}
