package response

import (
	"fmt"
	"time"
)

type AppResponse[T any] struct {
	Sucess 	bool		`json:"success"`
	Code 	string 		`json:"code"`
	Message string 		`json:"message"`
	Data 	*T 			`json:"data"`
	Date 	time.Time 	`json:"date"`
}

type AppError struct {
	Success bool
	Code string
	Message string
	Date time.Time
}


func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s, Date: %v", e.Code, e.Message, e.Date)
}

func NewAppError(code, message string) *AppError {
    return &AppError{
		Success: false,
        Code:    code,
        Message: message,
        Date:    time.Now(),
    }
}

func NewSuccessResponse[T any](data T, code, message string) AppResponse[T] {
	return AppResponse[T]{
		Sucess: true,
		Code: code,
		Message: message,
		Data: &data,
		Date: time.Now(),
	}
}

