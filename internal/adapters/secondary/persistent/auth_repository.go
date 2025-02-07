package persistent

import (
	"errors"
	"fmt"

	"github.com/hafiztri123/internal/core/entity"
	"gorm.io/gorm"
)


type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository (db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}



// auth_repository.go
func (a *AuthRepository) Register(user *entity.User) error  {
    existingUser, err := a.FindByEmail(user.Email)
    
    // If we found a user (no error), then it's a duplicate
    if err == nil && existingUser != nil {
        return fmt.Errorf("user with email %s already exists", user.Email)
    }
    
    // If error is anything other than "record not found", it's a real error
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("error checking existing user: %w", err)
    }

    // At this point, we know the email doesn't exist, so create the user
    result := a.db.Create(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (a *AuthRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := a.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
	
}