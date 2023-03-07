package domain

type User struct {
	Base     Base   `json:",inline"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
