package models

// User defines user object
type User struct {
	ID    string `json:"_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
