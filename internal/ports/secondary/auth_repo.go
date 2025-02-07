package secondary

import "github.com/hafiztri123/internal/core/entity"


type AuthRepository interface {
	Register (*entity.User) error
	FindByEmail(string) (*entity.User, error)
}