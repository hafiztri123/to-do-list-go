package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/adapters/primary/dto"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/response"
	"github.com/hafiztri123/internal/core/usecase"
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
    BindJSON(c, &req)
    
  
    task := &entity.Task{
        UserID:      userID,
        ParentID:    req.ParentID,
        Title:       req.Title,
        Description: req.Description,
        Status:      "pending",
        DueDate:     req.DueDate,
    }

    if err := h.service.CreateTask(task); err != nil {
        c.JSON(400, err)
        return
    }

    c.JSON(201, response.NewSuccessResponse(task, "201", "Task created successfully"))
}

func (h *TaskHandler) GetUserTasks(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    tasks, err := h.service.GetUserTasks(userID)
    if err != nil {
        c.JSON(404, err)
        return
    }

    c.JSON(200, response.NewSuccessResponse(tasks, "200", "Tasks fetched successfully"))
}

func (h *TaskHandler) GetSubTasks(c *gin.Context) {
    taskID, err := fetchIDFromParam(c)
    if err != nil {
        c.JSON(400, err)
        return
    }
	

    tasks, err := h.service.GetSubTasks(uint(taskID))
    if err != nil {
        c.JSON(404, err)
        return
    }

    c.JSON(200, response.NewSuccessResponse(tasks, "200", "Sub tasks fetched successfully"))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
    taskID, err := fetchIDFromParam(c)
    if err != nil {
        c.JSON(400, err)
        return
    }
	
    var req dto.UpdateTaskRequest
    BindJSON(c, &req)

    task := &entity.Task{
        ID:          uint(taskID),
        Title:       req.Title,
        Description: req.Description,
        Status:      req.Status,
        DueDate:     req.DueDate,
    }

    if err := h.service.UpdateTask(task); err != nil {
        c.JSON(400, err)
        return
    }

    c.JSON(200, response.NewSuccessResponse("", "200", "Task updated successfully"))
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
    taskID, err := fetchIDFromParam(c)
    if err != nil {
        c.JSON(400, err)
        return
    }
	
    
    if err := h.service.DeleteTask(uint(taskID)); err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, response.NewSuccessResponse("", "200", "Task deleted successfully"))
}

func fetchIDFromParam(c *gin.Context) (uint64, error) {
    taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)
	if err != nil {
		return 0, response.NewAppError("400", "Invalid task ID")
	}

    return taskID, nil

}
