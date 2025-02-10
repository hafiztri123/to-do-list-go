package http

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/adapters/primary/dto"
	"github.com/hafiztri123/internal/core/response"
	"github.com/hafiztri123/internal/core/usecase"
)



type CategoryHandler struct {
	service *usecase.CategoryService
}

func NewCategoryHandler(service *usecase.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}


func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	BindJSON(c, &req)
	req.Name = strings.ToLower(req.Name)

	if err := h.service.CreateCategory(c.GetUint("user_id"), strings.ToLower(req.Name)); err != nil {
		c.JSON(errorCode(err), err)
		return
	}

	c.JSON(201, response.NewSuccessResponse("", 201, "Category created successfully"))
	
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := fetchCategoryIDFromParam(c)
	if err != nil {
		c.JSON(errorCode(err), err)
		return
	}
	
	
	err = h.service.DeleteCategory(uint(categoryID))

	if err != nil {
		c.JSON(errorCode(err), err)
		return
	}

	c.JSON(204, response.NewSuccessResponse("",204, "Category deleted successfully"))
}

func(h *CategoryHandler) GetAllCategory(c *gin.Context) {

	categories, err := h.service.GetAllCategory(c.GetUint("user_id"))
	if err != nil {
		c.JSON(errorCode(err), err)
		return
	}

	c.JSON(200, response.NewSuccessResponse(categories, 200, "Categories fetched successfully"))
}

func fetchCategoryIDFromParam(c *gin.Context) (uint64, error) {

	categoryID, err := strconv.ParseUint(c.Param("category_id"), 10, 32)
	if err != nil {
		c.JSON(400, response.NewAppError(400, "Invalid category ID"))
		return 0, err
	}

	return categoryID, nil
}