package usecase

import (
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/response"
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
    if a.authRepo.IsEmailExist(user.Email){
        return response.NewAppError("400", "Email already exists")
    }


    if err := a.authRepo.Create(user); err != nil {
        return err
    }

    return nil
}

func (a *AuthService) FindByEmail(email string) (*entity.User, error) { 
    user, err := a.authRepo.FindByEmail(email)
    if err != nil {
        return nil, response.NewAppError("404", err.Error())
    }
    return user, nil
}