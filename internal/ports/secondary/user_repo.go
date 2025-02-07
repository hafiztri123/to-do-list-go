package secondary

import "github.com/hafiztri123/internal/core/entity"

type UserRepository interface {
	FindById(id uint) (*entity.User, error)
}

