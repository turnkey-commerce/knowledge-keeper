package models

// UserRegistration provides a password with the User model for registration.
type UserRegistration struct {
	User
	Password string `json:"password"` // password
}

// Token returns a token from login.
type Token struct {
	Token string `json:"token"` // password
}

// UserLogin provides the email and password for login.
type UserLogin struct {
	Email    string `json:"email"`    // password
	Password string `json:"password"` // password
}

// UserUpdate provides the user elements that can be updated.
type UserUpdate struct {
	Email     string `json:"email"`      // password
	Password  string `json:"password"`   // password
	FirstName string `json:"first_name"` // first_name
	LastName  string `json:"last_name"`  // last_name
	IsAdmin   bool   `json:"is_admin"`   // is_admin
}
