package viewmodel

import "time"

type UserView struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type SearchUserRequest struct {
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type LoginResponse struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	LoginAt time.Time `json:"login_at"`
}
