package model

type User struct {
	Base     `json:",inline"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Email    string `json:"email" gorm:"column:email"`
}

func (r User) TableName() string {
	return "users"
}

type Users []User
