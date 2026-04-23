package film_repository

import (
	"fmt"

	"github.com/matuha07/kinotower-go/src/internal/core/domain"
)

func (r *filmRepository) getCategories(filmID int) ([]domain.Category, error) {
	const query = `
		SELECT c.id, c.name
		FROM categories c
		INNER JOIN categories_films cf ON cf.category_id = c.id
		WHERE cf.film_id = $1
		  AND c.deleted_at IS NULL
		ORDER BY c.name ASC
	`

	rows := []CategoryRow{}
	if err := r.db.Select(&rows, query, filmID); err != nil {
		return nil, fmt.Errorf("failed to get categories for film %d: %w", filmID, err)
	}

	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, domain.Category{ID: row.ID, Name: row.Name})
	}

	return categories, nil
}

func (r *filmRepository) setFilmCategories(filmID int, categoryIDs []int) error {
	const clear = `DELETE FROM categories_films WHERE film_id = $1`
	if _, err := r.db.Exec(clear, filmID); err != nil {
		return err
	}

	if len(categoryIDs) == 0 {
		return nil
	}

	const insert = `INSERT INTO categories_films (category_id, film_id) VALUES ($1, $2)`
	for _, id := range categoryIDs {
		if _, err := r.db.Exec(insert, id, filmID); err != nil {
			return err
		}
	}

	return nil
}
