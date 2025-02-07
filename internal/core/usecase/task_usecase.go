package usecase

import (
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/rs/zerolog/log"
)

type TaskService struct {
    repo *persistent.TaskRepository
}

func NewTaskService(repo *persistent.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *entity.Task) error {
    log.Info().
        Uint("user_id", task.UserID).
        Str("title", task.Title).
        Msg("Creating new task")

    if err := s.repo.Create(task); err != nil {
        log.Error().Err(err).Msg("Failed to create task")
        return err
    }
    return nil
}

func (s *TaskService) GetUserTasks(userID uint) ([]entity.Task, error) {
    log.Info().
        Uint("user_id", userID).
        Msg("Fetching user tasks")

    tasks, err := s.repo.FindByUserID(userID)
    if err != nil {
        log.Error().Err(err).Msg("Failed to fetch user tasks")
        return nil, err
    }
    return tasks, nil
}

func (s *TaskService) GetSubTasks(taskID uint) ([]entity.Task, error) {
    log.Info().
        Uint("task_id", taskID).
        Msg("Fetching sub-tasks")

    tasks, err := s.repo.FindSubTasks(taskID)
    if err != nil {
        log.Error().Err(err).Msg("Failed to fetch sub-tasks")
        return nil, err
    }
    return tasks, nil
}

func (s *TaskService) UpdateTask(task *entity.Task) error {
    log.Info().
        Uint("task_id", task.ID).
        Msg("Updating task")

    if err := s.repo.Update(task); err != nil {
        log.Error().Err(err).Msg("Failed to update task")
        return err
    }
    return nil
}

func (s *TaskService) DeleteTask(taskID uint) error {
    log.Info().
        Uint("task_id", taskID).
        Msg("Deleting task")

    if err := s.repo.Delete(taskID); err != nil {
        log.Error().Err(err).Msg("Failed to delete task")
        return err
    }
    return nil
}