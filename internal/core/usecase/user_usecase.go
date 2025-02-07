package usecase

import (
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/ports/secondary"
	
)

type UserService struct {
	repo secondary.UserRepository
}

func NewUserService(repo secondary.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func(u *UserService) FindById(id uint) (*entity.User, error) {
	return u.repo.FindById(id)
}