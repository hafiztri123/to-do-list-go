package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/adapters/primary/dto"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/usecase"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)


type AuthHandler struct {
	service *usecase.AuthService
	jwtService *usecase.JWTService
}

func NewAuthHandler(service *usecase.AuthService, jwtService *usecase.JWTService) *AuthHandler {
	return &AuthHandler{
		service: service,
		jwtService: jwtService,
	}
}

func (a *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Error().
            Err(err).
            Str("endpoint", "register").
            Msg("Failed to bind request")
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := hashingPassword(req.Password)
    if err != nil {
        log.Error().
            Err(err).
            Str("endpoint", "register").
            Msg("Failed to hash password")
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    user := &entity.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }

    if err := a.service.Register(user); err != nil {
        log.Error().
            Err(err).
            Str("email", req.Email).
            Msg("Failed to register user")
        c.JSON(400, gin.H{"error": "Email already exist"})
        return
    }

    log.Info().
        Str("email", req.Email).
        Msg("User registered successfully")
    c.JSON(201, gin.H{"message": "register success"})
}

func (a *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Error().
            Err(err).
            Str("endpoint", "login").
            Msg("Failed to bind request")
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    user, err := a.service.FindByEmail(req.Email)
    if err != nil {
        log.Error().
            Err(err).
            Str("email", req.Email).
            Msg("User not found")
        c.JSON(404, gin.H{"error": "user not found"})
        return
    }

    if err := isPasswordMatch(user.Password, req.Password); err != nil {
        log.Error().
            Str("email", req.Email).
            Msg("Invalid password attempt")
        c.JSON(401, gin.H{"error": "invalid password"})
        return
    }

	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("Failed to generate token")
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

    log.Info().
        Str("email", req.Email).
        Msg("User logged in successfully")
    c.JSON(200, gin.H{"message": token})
}

func hashingPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func isPasswordMatch(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

