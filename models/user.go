package models

// User defines user object
type User struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
