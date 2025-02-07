package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/adapters/primary/dto"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/usecase"
	"github.com/rs/zerolog/log"
)

type TaskHandler struct {
    service *usecase.TaskService
}

func NewTaskHandler(service *usecase.TaskService) *TaskHandler {
    return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
    userID := c.GetUint("user_id") // From JWT middleware
    var req dto.CreateTaskRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Error().Err(err).Msg("Invalid request body")
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    task := &entity.Task{
        UserID:      userID,
        ParentID:    req.ParentID,
        Title:       req.Title,
        Description: req.Description,
        Status:      "pending",
        DueDate:     req.DueDate,
    }

    if err := h.service.CreateTask(task); err != nil {
        c.JSON(500, gin.H{"error": "Failed to create task"})
        return
    }

    c.JSON(201, task)
}

func (h *TaskHandler) GetUserTasks(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    tasks, err := h.service.GetUserTasks(userID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch tasks"})
        return
    }

    c.JSON(200, tasks)
}

func (h *TaskHandler) GetSubTasks(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)  
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID"})
		return
	}  

    tasks, err := h.service.GetSubTasks(uint(taskID))
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch sub-tasks"})
        return
    }

    c.JSON(200, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
    taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID"})
		return
	}
    var req dto.UpdateTaskRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    task := &entity.Task{
        ID:          uint(taskID),
        Title:       req.Title,
        Description: req.Description,
        Status:      req.Status,
        DueDate:     req.DueDate,
    }

    if err := h.service.UpdateTask(task); err != nil {
        c.JSON(500, gin.H{"error": "Failed to update task"})
        return
    }

    c.JSON(200, gin.H{"message": "Task updated successfully"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
    taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID"})
		return
	}
	
    
    if err := h.service.DeleteTask(uint(taskID)); err != nil {
        c.JSON(500, gin.H{"error": "Failed to delete task"})
        return
    }

    c.JSON(200, gin.H{"message": "Task deleted successfully"})
}