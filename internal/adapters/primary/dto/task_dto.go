package dto

import "time"

type CreateTaskRequest struct {
    Category    uint    `json:"category_id"`
    ParentID    *uint     `json:"parent_id"`
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description"`
    DueDate     time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    DueDate     time.Time `json:"due_date"`
}