package entity

type User struct {
	ID uint `json:"id" gorm:"primarykey"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}