package domain

type Filter struct {
	Page       int
	Size       int
	SortBy     string
	SortDir    string
	GenreID    int
	CountryID  int
	CategoryID int
	Search     string
}

func (f *Filter) Limit() int {
	if f.Size <= 0 {
		return 10
	}
	return f.Size
}

func (f *Filter) Offset() int {
	if f.Page <= 0 {
		return 0
	}
	return (f.Page - 1) * f.Limit()
}
