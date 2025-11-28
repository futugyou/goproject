package viewmodel

type UserView struct {
	ID       string `json:"id"`
	Name     string `json:"name"` 
	Email    string `json:"email"`
}

type SearchUserRequest struct {
	Name string `json:"name"`
}