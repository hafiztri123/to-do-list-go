package usecase

import (
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
)


type AuthService struct {
	authRepo *persistent.AuthRepository
	
}

func NewAuthService(authrepo *persistent.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authrepo,
	}
}

func (a *AuthService) Register(user *entity.User) error {


    if err := a.authRepo.Register(user); err != nil {
        return err
    }

    return nil
}

func (a *AuthService) FindByEmail(email string) (*entity.User, error) {
	return a.authRepo.FindByEmail(email)
}