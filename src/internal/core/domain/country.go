package domain

type Country struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FilmCount *int   `json:"filmCount,omitempty"`
}
