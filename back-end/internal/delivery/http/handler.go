package http

import (
	"errors"
	"net/http"

	"github.com/ReilEgor/CleanArchitectureGolang/internal/domain"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	uc domain.TaskUsecase
}

func NewTaskHandler(uc domain.TaskUsecase) *TaskHandler {
	return &TaskHandler{uc: uc}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.uc.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	task, err := h.uc.AddTask(input.Title)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, task)
}
