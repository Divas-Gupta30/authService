package models

// User struct represents a user in memory
type User struct {
	Email    string
	Password string
}

// Token struct represents a token in memory
type Token struct {
	UserID  string
	Token   string
	Expires int64
}
