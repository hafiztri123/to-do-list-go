package entity

import "time"

type Task struct {
    ID          uint      `json:"id" gorm:"primarykey"`
    UserID      uint      `json:"user_id"`
    ParentID    *uint     `json:"parent_id"` // Nullable, for nested tasks
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status" gorm:"type:varchar(20)"` // e.g., "pending", "completed"
    DueDate     time.Time `json:"due_date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    SubTasks    []Task    `json:"sub_tasks" gorm:"foreignKey:ParentID"` // For nested tasks
    User        User      `json:"-" gorm:"foreignKey:UserID"`
}