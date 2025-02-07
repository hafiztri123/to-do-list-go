package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/adapters/primary/dto"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/response"
	"github.com/hafiztri123/internal/core/usecase"
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
    BindJSON(c, &req)

    hashedPassword, err := hashingPassword(req.Password)
    if err != nil {
        c.JSON(500, err)
        return
    }

    user := &entity.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }

    if err := a.service.Register(user); err != nil {
        c.JSON(400, err)
        return
    }

    c.JSON(201, response.NewSuccessResponse(user, "201", "User created successfully"))
}

func (a *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    BindJSON(c, &req)

    user, err := a.service.FindByEmail(req.Email)
    if err != nil {
        c.JSON(404, err)
        return
    }

    if err := isPasswordMatch(user.Password, req.Password); err != nil {
        c.JSON(401, err)
        return
    }

	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {

		c.JSON(500, err)
		return
	}


    c.JSON(200, response.NewSuccessResponse(token, "200", "Login successful"))
}

func hashingPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
        return "", &response.AppError{
            Code: "500",
            Message: "Failed to hash password",
        }
	}
	return string(hashedPassword), nil
}

func isPasswordMatch(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return response.NewAppError("401", "Invalid credentials")
    }
    return nil
}

func BindJSON(c *gin.Context, v interface{}) error {
    if err := c.ShouldBindJSON(v); err != nil {
        return &response.AppError{
            Code:    "400",
            Message: err.Error(),
        }
    }
    return nil
}



