package usecase_test

import (
	"testing"

	"github.com/ReilEgor/CleanArchitectureGolang/internal/domain"
	"github.com/ReilEgor/CleanArchitectureGolang/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) GetAll() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func TestAddTask(t *testing.T) {
	repo := new(MockRepository)
	uc := usecase.NewTaskUsecase(repo)

	title := "Test Task"

	repo.On("Create", mock.MatchedBy(func(task *domain.Task) bool {
		return task.Title == title && task.Done == false
	})).Return(nil)

	task, err := uc.AddTask(title)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, title, task.Title)

	repo.AssertExpectations(t)
}

func TestAddTask_EmptyTitle(t *testing.T) {
	repo := new(MockRepository)
	uc := usecase.NewTaskUsecase(repo)

	task, err := uc.AddTask("")

	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, domain.ErrEmptyTitle, err)

	repo.AssertNotCalled(t, "Create", mock.Anything)
}
