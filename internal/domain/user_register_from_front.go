package domain

type UserRegisterFromFront struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
