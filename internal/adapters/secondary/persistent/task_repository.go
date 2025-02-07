package persistent

import (
	"github.com/hafiztri123/internal/core/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TaskRepository struct{
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *entity.Task) error {
    log.Info().
        Uint("user_id", task.UserID).
        Str("title", task.Title).
        Msg("Creating task in database")

    result := r.db.Create(task)
    if result.Error != nil {
        log.Error().
            Err(result.Error).
            Msg("Database error while creating task")
        return result.Error
    }
    return nil
}

func(r *TaskRepository) FindByID(id uint) (*entity.Task, error) {
	var task entity.Task

    result := r.db.First(&task, id)
    if result.Error != nil {
        log.Error().
            Err(result.Error).
            Uint("task_id", id).
            Msg("Failed to find task")
        return nil, result.Error
    }

    return &task, nil
		
}

func (r *TaskRepository) FindByUserID(userID uint) ([]entity.Task, error) {
    var tasks []entity.Task
    
    // Get tasks with their subtasks 
    result := r.db.Debug(). // Add Debug() to see the SQL query
        Where("user_id = ? AND parent_id IS NULL", userID).
        Preload("SubTasks"). // This will load one level of subtasks
        Find(&tasks)
    
    if result.Error != nil {
        log.Error().
            Err(result.Error).
            Uint("user_id", userID).
            Msg("Failed to fetch user tasks")
        return nil, result.Error
    }
    
    return tasks, nil
}
func (r *TaskRepository) FindSubTasks(taskID uint) ([]entity.Task, error) {
    var tasks []entity.Task
    
    result := r.db.Debug().Where("parent_id = ?", taskID).
        Preload("SubTasks"). // Also load next level of subtasks
        Find(&tasks)
    
    if result.Error != nil {
        log.Error().
            Err(result.Error).
            Uint("task_id", taskID).
            Msg("Failed to fetch subtasks")
        return nil, result.Error
    }
    
    return tasks, nil
}

func (r *TaskRepository) Update(task *entity.Task) error {
    // First check if task exists
    existingTask, err := r.FindByID(task.ID)
    if err != nil {
        return err
    }

    // Only update allowed fields
    updates := map[string]interface{}{
        "title":       task.Title,
        "description": task.Description,
        "status":      task.Status,
        "due_date":    task.DueDate,
    }

    result := r.db.Model(existingTask).Updates(updates)
    if result.Error != nil {
        log.Error().
            Err(result.Error).
            Uint("task_id", task.ID).
            Msg("Failed to update task")
        return result.Error
    }

    return nil
}


func (r *TaskRepository) Delete(id uint) error {
    // Start a transaction because we need to delete subtasks too
    return r.db.Transaction(func(tx *gorm.DB) error {
        // First recursively delete all subtasks
        if err := r.deleteSubTasks(tx, id); err != nil {
            log.Error().
                Err(err).
                Uint("task_id", id).
                Msg("Failed to delete subtasks")
            return err
        }

        // Then delete the task itself
        if result := tx.Delete(&entity.Task{}, id); result.Error != nil {
            log.Error().
                Err(result.Error).
                Uint("task_id", id).
                Msg("Failed to delete task")
            return result.Error
        }

        return nil
    })
}



func (r *TaskRepository) deleteSubTasks(tx *gorm.DB, taskID uint) error {
    var subtasks []entity.Task
    if err := tx.Where("parent_id = ?", taskID).Find(&subtasks).Error; err != nil {
        return err
    }

    // Recursively delete each subtask's subtasks
    for _, subtask := range subtasks {
        if err := r.deleteSubTasks(tx, subtask.ID); err != nil {
            return err
        }
    }

    // Delete all subtasks of this task
    if err := tx.Where("parent_id = ?", taskID).Delete(&entity.Task{}).Error; err != nil {
        return err
    }

    return nil
}


