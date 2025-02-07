package persistent

import (

	"github.com/hafiztri123/internal/core/entity"
	"gorm.io/gorm"
)


type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository (db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}



func (a *AuthRepository) Create(user *entity.User) error  {

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

func (a *AuthRepository) IsEmailExist(email string) bool {
	var count int64
	a.db.Model(&entity.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}