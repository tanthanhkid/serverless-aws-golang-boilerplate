package models

// User schema of the user table
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
}
