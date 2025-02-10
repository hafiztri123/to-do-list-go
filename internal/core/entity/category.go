package entity

import "time"

type Category struct {
	ID 			uint 		`json:"id" gorm:"primarykey"`
	UserID 		uint 		`json:"user_id"`
	Name 		string 		`json:"name"`
	Tasks 		[]Task 		`json:"tasks" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}