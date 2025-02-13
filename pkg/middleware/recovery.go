package middleware

import (
	"log" // Standard library logger
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/core/response"
)

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf(
                    "Server error recovered: %v | Method: %s | Path: %s",
                    err,
                    c.Request.Method,
                    c.Request.URL.Path,
                )

                appError := response.NewAppError(
                    http.StatusInternalServerError,
                    "An unexpected error occurred",
                )
                c.AbortWithStatusJSON(appError.Code, appError)
                return
            }
        }()
		
        c.Next()
    }
}