package domain

type Category struct {
	ID             int
	Name           string
	ParentCategory *Category
	FilmCount      *int
}
