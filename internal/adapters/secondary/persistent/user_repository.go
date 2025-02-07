package persistent

import (
	"github.com/hafiztri123/internal/core/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindById(id uint) (*entity.User, error) {
	var user entity.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func(r *UserRepository) IsUserExistByID(id uint) bool {
	var count int64
	r.db.Model(&entity.User{}).Where("id = ?", id).Count(&count)
	return count > 0
}