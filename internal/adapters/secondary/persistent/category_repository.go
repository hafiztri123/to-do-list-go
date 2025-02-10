package persistent

import (
	"time"

	"github.com/hafiztri123/internal/core/entity"
	"gorm.io/gorm"
)


type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(userid uint, name string) error {
	
	category := &entity.Category{
		UserID: userid,
		Name: name,
		CreatedAt: time.Now(),
	}

	result := r.db.Create(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CategoryRepository) Delete(id uint) error {


	result := r.db.Delete(&entity.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CategoryRepository) FindByUserID(userID uint) (*[]entity.Category, error) {
	var category []entity.Category
	result := r.db.Where("user_id = ?", userID).Find(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}



func(r *CategoryRepository) IsCategoryExistByID(id uint) bool {
	var count int64
	r.db.Model(&entity.Category{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func(r *CategoryRepository) IsCategoryExistByNameAndUserID(name string, userID uint) bool {
	var count int64
	r.db.Model(&entity.Category{}).Where("name = ? AND user_id = ?", name, userID).Count(&count)
	return count > 0
}