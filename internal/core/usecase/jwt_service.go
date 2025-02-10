package usecase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/internal/core/response"
	"github.com/hafiztri123/pkg/config"
)

type JWTService struct {
	config *config.JWTConfig
	
}

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTService(config *config.JWTConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}

func (j *JWTService) GenerateToken(userID uint) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(j.config.TokenLifetime))),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.SecretKey))

	if err != nil {
		return "", response.NewAppError(500, "Failed to generate token")
	}

	return tokenString, nil
}

func(j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, response.NewAppError(401, "Invalid signing method")
		}

		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, response.NewAppError(401, "Invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, response.NewAppError(401, "Invalid token")
	}
	
	return claims, nil
}


