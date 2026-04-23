package domain

import "time"

type User struct {
	ID          int       `json:"id"`
	FIO         string    `json:"fio"`
	Email       string    `json:"email"`
	Birthday    string    `json:"birthday"`
	CountryID   *int      `json:"country_id,omitempty"`
	GenderID    *int      `json:"gender_id,omitempty"`
	Gender      *Gender   `json:"gender,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ReviewCount int       `json:"reviewCount"`
	RatingCount int       `json:"ratingCount"`
}

type Gender struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserCreate struct {
	FIO       string
	Email     string
	Password  string
	Birthday  string
	CountryID *int
	GenderID  *int
}

type UserUpdate struct {
	FIO       *string
	Email     *string
	Birthday  *string
	CountryID *int
	GenderID  *int
}

type UserReview struct {
	ID         int       `json:"id"`
	Film       FilmShort `json:"film"`
	Message    string    `json:"message"`
	IsApproved int       `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
}

type FilmReview struct {
	ID        int       `json:"id"`
	User      UserShort `json:"user"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRating struct {
	ID        int       `json:"id"`
	Film      FilmShort `json:"film"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type FilmShort struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserShort struct {
	ID  int    `json:"id"`
	FIO string `json:"fio"`
}

type ReviewCreate struct {
	FilmID  int
	Message string
}

type RatingCreate struct {
	FilmID int
	Ball   int
}
