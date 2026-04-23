package film_repository

import (
	"strconv"
	"strings"
	"time"

	"github.com/matuha07/kinotower-go/src/internal/core/domain"
)

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err == nil {
		return parsed
	}
	parsed, err = time.Parse(time.RFC3339, value)
	if err == nil {
		return parsed
	}
	parsed, _ = time.Parse("2006-01-02 15:04:05-07", value)
	return parsed
}

func safeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func itoa(value int) string {
	return strconv.Itoa(value)
}

func join(values []string, sep string) string {
	return strings.Join(values, sep)
}

func rowToFilm(row FilmRow) domain.Film {
	film := domain.Film{
		ID:            row.ID,
		Name:          row.Name,
		Duration:      row.Duration,
		YearOfIssue:   row.YearOfIssue,
		Age:           row.Age,
		LinkImg:       row.LinkImg,
		LinkKinopoisk: row.LinkKinopoisk,
		LinkVideo:     row.LinkVideo,
		CreatedAt:     parseTime(row.CreatedAt),
		RatingAvg:     row.RatingAvg,
		ReviewCount:   row.ReviewCount,
	}
	if row.CountryID != nil {
		film.Country = &domain.Country{ID: *row.CountryID, Name: safeString(row.CountryName)}
	}
	return film
}
