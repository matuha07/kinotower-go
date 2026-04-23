package domain

type SignUpInput struct {
	FIO       string
	Email     string
	Password  string
	Birthday  string
	CountryID *int
	GenderID  *int
}

type SignInInput struct {
	Email    string
	Password string
}

type TokenPair struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	ID     int    `json:"id"`
	FIO    string `json:"fio"`
}

type InvalidAuthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
