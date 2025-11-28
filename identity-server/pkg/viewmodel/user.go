package viewmodel

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
