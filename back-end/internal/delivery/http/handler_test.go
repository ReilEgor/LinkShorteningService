package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/ReilEgor/CleanArchitectureGolang/internal/delivery/http"
	"github.com/ReilEgor/CleanArchitectureGolang/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUsecase struct {
	mock.Mock
}

func (m *MockUsecase) AddTask(title string) (*domain.Task, error) {
	args := m.Called(title)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockUsecase) GetTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func TestGetTasksHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUc := new(MockUsecase)
	handler := api.NewTaskHandler(mockUc)

	r := gin.Default()
	r.GET("/tasks", handler.GetTasks)

	mockTasks := []domain.Task{
		{ID: 1, Title: "Test 1", Done: false},
		{ID: 2, Title: "Test 2", Done: true},
	}
	mockUc.On("GetTasks").Return(mockTasks, nil)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "Test 1")
	assert.Contains(t, w.Body.String(), "Test 2")

	mockUc.AssertExpectations(t)
}
