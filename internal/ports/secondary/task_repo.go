package secondary

import "github.com/hafiztri123/internal/core/entity"

type TaskRepository interface {
    Create(task *entity.Task) error
    FindByID(id uint) (*entity.Task, error)
    FindByUserID(userID uint) ([]entity.Task, error)
    FindSubTasks(taskID uint) ([]entity.Task, error)
    Update(task *entity.Task) error
    Delete(id uint) error
}