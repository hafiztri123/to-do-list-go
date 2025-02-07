package usecase

import (
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/response"
)

type UserService struct {
	repo persistent.UserRepository
}

func NewUserService(repo persistent.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func(u *UserService) FindById(id uint) (*entity.User, error) {
	user, err := u.repo.FindById(id)
	if err != nil {
		return nil, response.NewAppError("404", "User not found")
	}
	return user, nil
}