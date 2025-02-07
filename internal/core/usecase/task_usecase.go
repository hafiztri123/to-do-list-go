package usecase

import (
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/response"
)

type TaskService struct {
    taskRepo *persistent.TaskRepository
    userRepo *persistent.UserRepository
    
}

func NewTaskService(taskRepo *persistent.TaskRepository, userRepo *persistent.UserRepository) *TaskService {
    return &TaskService{
        taskRepo: taskRepo,
        userRepo: userRepo,
    }
}

func (s *TaskService) CreateTask(task *entity.Task) error {

    if err := s.taskRepo.Create(task); err != nil {
        return response.NewAppError("500", err.Error())
    }
    return nil
}

func (s *TaskService) GetUserTasks(userID uint) ([]entity.Task, error) {
    if !s.userRepo.IsUserExistByID(userID) {
        return nil, response.NewAppError("404", "User not found")
    }

    tasks, err := s.taskRepo.FindByUserID(userID)
    if err != nil {
        return nil, response.NewAppError("500", err.Error())

    }
    return tasks, nil
}

func (s *TaskService) GetSubTasks(taskID uint) ([]entity.Task, error) {
    if !s.taskRepo.IsTaskExistByID(taskID) {
        return nil, response.NewAppError("404", "Task not found")
    }

    tasks, err := s.taskRepo.FindSubTasks(taskID)
    if err != nil {
        return nil, response.NewAppError("500", err.Error())
    }

    return tasks, nil
}

func (s *TaskService) UpdateTask(task *entity.Task) error {

    if err := s.taskRepo.Update(task); err != nil {
        return response.NewAppError("500", err.Error())
    }
    return nil
}

func (s *TaskService) DeleteTask(taskID uint) error {

    if err := s.taskRepo.Delete(taskID); err != nil {
        return response.NewAppError("500", err.Error())
    }
    return nil
}