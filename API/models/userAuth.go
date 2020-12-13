package models

// UserAuth provides a password.
type UserAuth struct {
	User
	Password string `json:"password"` // password
}
