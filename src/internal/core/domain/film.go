package domain

type Film struct {
	ID            int
	Name          string
	Description   string
	ReleaseYear   int
	Rating        float64
	Age           int
	Genre         string
	LinkImg       string
	LinkKinopoisk string
	LinkVideo     string
	CreatedAt     string
	Country       Country
	Category      []Category
	RatingAvg     float64
	RewiewsInt    int
}
