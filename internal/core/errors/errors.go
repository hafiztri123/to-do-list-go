package errors

import (
	"fmt"
	"time"
)

type AppError struct {
	Code string
	Message string
	Path string
	Date time.Time
}


func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s, Path: %s, Date: %v", e.Code, e.Message, e.Path, e.Date)
}

func NewAppError(code, message, path string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Path:    path,
        Date:    time.Now(),
    }
}