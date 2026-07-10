package models

import "time"

// User represents an administrator in the database.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password hash in JSON
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest defines parameters for the login API
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse defines the returned payload upon successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
