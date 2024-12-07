package models

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Token struct {
	UserID  string
	Token   string
	Expires int64
}
