package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/core/usecase"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(jwtService *usecase.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Error().Msg("No authorization header")
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		} 

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			log.Error().Msg("Invalid token format")
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(bearerToken[1])
		if err != nil {
			log.Error().Err(err).Msg("Failed to validate token")
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}